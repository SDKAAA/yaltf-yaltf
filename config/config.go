// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package config

const Version = "1.0"

type Target struct {
	Host string `toml:"host,omitempty"`
	Port string `toml:"port,omitempty"`
	User string `toml:"user,omitempty"`
}

type CommonOpts struct {
	SSHKeyPath string `toml:"sshkeypath"`
	SSHTimeout int    `toml:"timeout_seconds"`
}

// Targets represents the computers that will be scanned.
type Config struct {
	Common  CommonOpts
	Targets map[string]Target
}

var Conf Config
