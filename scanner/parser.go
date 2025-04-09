// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package scanner

import (
	"strings"

	"github.com/yaltf/yaltf/models"
)

func parseRPM(stdout string, licenseInfo models.LicenseInfo) {
	lines := strings.Split(stdout, "\n")

	for _, line := range lines {
		// Examples for line:
		// samba-common GPL-3.0-or-later AND LGPL-3.0-or-later
		// intel-audio-firmware LicenseRef-Callaway-Redistributable-no-modification-permitted
		trimmed := strings.TrimSpace(line)
		name, license, found := strings.Cut(trimmed, " ")

		if !found {
			continue
		}

		licenseInfo[name] = license
	}
}

func parseDPKG(stdout string, licenseInfo models.LicenseInfo) {
	lines := strings.Split(stdout, "\n")

	for _, line := range lines {
		// Examples for line:
		// busybox unknown
		// bzip2 BSD-variant GPL-2
		// ca-certificates GPL-2+ MPL-2.0
		trimmed := strings.TrimSpace(line)
		name, license, found := strings.Cut(trimmed, " ")

		if !found {
			continue
		}

		licenseInfo[name] = license
	}
}

func parseWIN(stdout string, licenseInfo models.LicenseInfo) {
	lines := strings.Split(stdout, "\n")

	for _, line := range lines {
		// Examples for line:
		// samba-common GPL-3.0-or-later AND LGPL-3.0-or-later
		// intel-audio-firmware LicenseRef-Callaway-Redistributable-no-modification-permitted
		trimmed := strings.TrimSpace(line)
		name, license, found := strings.Cut(trimmed, "|&|")

		if !found {
			continue
		}

		licenseInfo[name] = license
	}
}
