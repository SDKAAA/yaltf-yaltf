// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package models

// LicenseInfo holds package, license pairs.
type LicenseInfo map[string]string

type SingleResult struct {
	TargetName  string
	LicenseInfo LicenseInfo
}

// LicenseInfos holds target, LicenseInfo pairs.
type LicenseInfos map[string]LicenseInfo

type ScanResult struct {
	Version   string       `json:"version"`
	ScannedAt string       `json:"scannedAt"`
	Results   LicenseInfos `json:"results"`
}
