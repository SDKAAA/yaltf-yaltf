// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package models

import (
	"encoding/json"
	"time"
)

// LicenseInfo holds (package, license) pairs.
type LicenseInfo map[string]string

type SingleResult struct {
	TargetName  string
	LicenseInfo LicenseInfo
}

// LicenseInfos holds target, LicenseInfo pairs.
type LicenseInfos map[string]LicenseInfo

type ErrorLevel int

const (
	Warning ErrorLevel = iota
	Critical
)

func (el ErrorLevel) String() string {
	switch el {
	case Warning:
		return "Warning"
	case Critical:
		return "Critical"
	default:
		return "Invalid Error Level"
	}
}

func (el ErrorLevel) MarshalJSON() ([]byte, error) {
	// It is assumed Suit implements fmt.Stringer.
	return json.Marshal(el.String())
}

type Error struct {
	Time    time.Time  `json:"time"`
	Level   ErrorLevel `json:"severity"`
	Message string     `json:"message"`
}

// Errors holds target, Error pairs
type Errors map[string]Error

type SingleError struct {
	TargetName string
	Error      Error
}

type ScanResult struct {
	Version   string       `json:"version"`
	ScannedAt string       `json:"scannedAt"`
	Results   LicenseInfos `json:"results"`
	Errors    Errors       `json:"errors,omitempty"`
}
