// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/yaltf/yaltf/logging"
	"github.com/yaltf/yaltf/models"
)

func WriteResults(scanResult *models.ScanResult) {
	outputBytes, err := json.MarshalIndent(scanResult, "", "  ")

	if err != nil {
		slog.Error("Failed to convert results to json.")
		return
	}

	fileName := fmt.Sprintf("%s.json", scanResult.ScannedAt)
	// fileName := fmt.Sprintf("%s-Licenses.json", scanResult.ScannedAt)
	// fileName := fmt.Sprintf("%s-Versions.json", scanResult.ScannedAt)

	err = os.Mkdir("results", os.ModePerm)

	if errors.Is(err, fs.ErrPermission) {
		slog.Error("Failed to create 'results' directory.", "error", err.Error())
		return
	}

	fullPath := filepath.Join("results", fileName)

	logging.InfoLog.Printf("Results written to: %s", fullPath)

	err = os.WriteFile(fullPath, outputBytes, 0666)

	if err != nil {
		slog.Error("Failed to write results to file.", "error", err.Error())
		return
	}
}
