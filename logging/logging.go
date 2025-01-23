// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package logging

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func init() {
	InfoLog = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lmsgprefix)
	ErrorLog = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
