// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package scanner

import (
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/yaltf/yaltf/config"
	"github.com/yaltf/yaltf/models"
	"github.com/yaltf/yaltf/reporter"

	"golang.org/x/crypto/ssh"
)

type Scanner struct {
	Targets map[string]config.Target
}

// Scan execute scan
func (s Scanner) Scan() error {
	var wg sync.WaitGroup

	scanResult := models.ScanResult{
		Version:   config.Version,
		ScannedAt: time.Now().Format("2006-01-02T15:04:05"),
		Results:   make(models.LicenseInfos),
	}

	resultsCh := make(chan models.SingleResult, len(s.Targets))

	wg.Add(1)
	go s.collectResults(&scanResult, resultsCh, &wg)

	for name, target := range s.Targets {
		wg.Add(1)
		go scanTarget(name, target, resultsCh, &wg)
	}

	wg.Wait()

	reporter.WriteResults(&scanResult)
	return nil
}

func scanTarget(name string, target config.Target, resultsCh chan models.SingleResult, wg *sync.WaitGroup) {
	defer wg.Done()

	slog.Info("Scanning target.", "name", name)

	keyPath := os.ExpandEnv(config.Conf.Common.SSHKeyPath)
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		slog.Error("Failed to read private key.", "error", err)
		resultsCh <- models.SingleResult{TargetName: name}
		return
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		slog.Error("Failed to parse private key.", "error", err)
		resultsCh <- models.SingleResult{TargetName: name}
		return
	}

	sshClientConfig := &ssh.ClientConfig{
		User: target.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := net.JoinHostPort(target.Host, target.Port)
	client, err := ssh.Dial("tcp", address, sshClientConfig)
	if err != nil {
		slog.Error("Connection failed.", "host", target.Host, "error", err)
		resultsCh <- models.SingleResult{TargetName: name}
		return
	}

	slog.Debug("Connected to host.", "host", target.Host)
	defer client.Close()

	licenseInfo := make(models.LicenseInfo)

	switch targetOS := getTargetOS(client); targetOS {
	case "fedora", "opensuse-leap", "centos":
		output, err := runCommand(client, `rpm -qa --queryformat "%{NAME} %{LICENSE}\n"`)

		if err != nil {
			slog.Error("Failed to query packages.")
			return
		}

		parseRPM(string(output), licenseInfo)
	}

	resultsCh <- models.SingleResult{TargetName: name, LicenseInfo: licenseInfo}

	slog.Info("Finished scanning.", "name", name)
}

func (s Scanner) collectResults(scanResult *models.ScanResult, resultsCh chan models.SingleResult, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0

	for singleResult := range resultsCh {
		scanResult.Results[singleResult.TargetName] = singleResult.LicenseInfo
		i++
		if i == len(s.Targets) {
			break
		}
	}
}

func runCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		slog.Error("failed to create session.", "error", err)
		return nil, err
	}
	defer session.Close()

	output, err := session.Output(command)

	if err != nil {
		slog.Error("failed to start command", "command", command, "error", err)
		return nil, err
	}

	return output, nil
}

func getTargetOS(client *ssh.Client) string {
	output, err := runCommand(client, "cat /etc/os-release")

	if err != nil {
		slog.Error("Failed to read /etc/os-release.", "error", err.Error())
		return "Unknown"
	}

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		key, value, _ := strings.Cut(trimmed, "=")
		if strings.ToUpper(key) == "ID" {
			return strings.Trim(value, `'"`)
		}
	}

	return "Unknown"
}

func (s Scanner) LocalScan() {
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		slog.Error("Failed to read /etc/os-release.", "error", err.Error())
		return
	}

	lines := strings.Split(string(b), "\n")
	osType := "Unknown"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		key, value, _ := strings.Cut(trimmed, "=")
		if strings.ToUpper(key) == "ID" {
			osType = strings.Trim(value, `'"`)
			break
		}
	}

	licenseInfo := make(models.LicenseInfo)

	switch osType {
	case "fedora", "opensuse-leap", "centos":
		cmd := exec.Command("rpm", "-qa", "--queryformat", "%{NAME} %{LICENSE}\n")
		output, err := cmd.Output()
		if err != nil {
			slog.Error("Failed to query packages.")
		}

		parseRPM(string(output), licenseInfo)
	}

	scanResult := models.ScanResult{
		Version:   config.Version,
		ScannedAt: time.Now().Format("2006-01-02T15:04:05"),
		Results:   models.LicenseInfos{"localhost": licenseInfo},
	}

	reporter.WriteResults(&scanResult)
}
