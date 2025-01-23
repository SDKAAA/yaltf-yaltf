// SPDX-License-Identifier: GPL-3.0-only.
// SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.
//
// This software is licensed under GPL v3.0 license.
// See LICENSE file at the root of the project.

package subcmds

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/yaltf/yaltf/config"
	"github.com/yaltf/yaltf/logging"
	"github.com/yaltf/yaltf/scanner"

	"github.com/BurntSushi/toml"
	"github.com/google/subcommands"
)

type ScanCmd struct {
	localhost bool
}

func (*ScanCmd) Name() string { return "scan" }

func (*ScanCmd) Synopsis() string { return "Scan licenses" }

func (*ScanCmd) Usage() string {
	return `scan [-localhost]:
	Scans licenses of packages on targets (defined in config.toml) or localhost.
`
}

func (p *ScanCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.localhost, "localhost", false, "Scan localhost only.")
}

func (p *ScanCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logging.InfoLog.Printf("yaltf v%s started.", config.Version)
	logging.InfoLog.Print("Reading configuration.")

	_, err := toml.DecodeFile("config.toml", &config.Conf)

	if err != nil {
		slog.Error("Could not decode config.toml file.", "error", err)
		os.Exit(1)
	}

	logging.InfoLog.Print("Start scanning.")

	targets := config.Conf.Targets

	s := scanner.Scanner{
		Targets: targets,
	}

	if p.localhost {
		s.LocalScan()
	} else {
		s.Scan()
	}

	logging.InfoLog.Print("Scan finished.")

	return subcommands.ExitSuccess
}
