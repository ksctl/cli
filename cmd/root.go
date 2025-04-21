// Copyright 2025 Ksctl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/telemetry"

	cLogger "github.com/ksctl/cli/v2/pkg/logger"
	"github.com/spf13/cobra"
)

// for the newController we should be able to pass some option fields for control more things
// for example whther it is a dry-run for testing

func (k *KsctlCommand) NewRootCmd() *cobra.Command {

	v := false

	cmd := &cobra.Command{
		Use:   "ksctl",
		Short: "CLI tool for managing multiple K8s clusters",
		Long:  "CLI tool which can manage multiple K8s clusters from local clusters to cloud provider specific clusters.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			telemetry.IntegrityCheck()

			if k.debugMode {
				k.CliLog.Box(k.Ctx, "CLI Mode", "CLI is running in debug mode")
				k.menuDriven = cli.NewDebugMenuDriven()
			} else {
				k.menuDriven = cli.NewMenuDriven()
			}

			_ = k.menuDriven.GetProgressAnimation() // Just boot it up...

			if v {
				k.CliLog.Box(k.Ctx, "CLI Mode", "Verbose mode is enabled")
				k.verbose = -1
			}

			k.l = cLogger.NewLogger(k.verbose, os.Stdout)

			k.telemetry = telemetry.NewTelemetry(k.KsctlConfig.Telemetry)

			cmdName := cmd.Name()
			if cmdName != "self-update" && cmdName != "version" {
				hasUpdates, err := k.CheckForUpdates()
				if err == nil && hasUpdates {
					k.NotifyAvailableUpdates()
				}
			}
		},
	}

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cli.AddDebugMode(cmd, &k.debugMode)
	cli.AddVerboseFlag(cmd, &v)

	return cmd
}

func (k *KsctlCommand) Cluster() *cobra.Command {

	cmd := &cobra.Command{
		Use: "cluster",
		Example: `
ksctl cluster --help
		`,
		Short: "Use to work with clusters",
		Long:  "It is used to work with cluster",
	}

	return cmd
}
