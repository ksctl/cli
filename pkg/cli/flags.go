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

package cli

import "github.com/spf13/cobra"

func MarkFlagsRequired(command *cobra.Command, flagNames ...string) error {
	for _, flagName := range flagNames {
		if err := command.MarkFlagRequired(flagName); err != nil {
			return err
		}
	}
	return nil
}

func AddVerboseFlag(command *cobra.Command, verbose *int) {
	command.PersistentFlags().IntVarP(verbose, "verbose", "v", 0, "Enable verbose output")
}

func AddDryRunFlag(command *cobra.Command, dryRun *bool) {
	command.PersistentFlags().BoolVar(dryRun, "dry-run", false, "Perform a dry run")
}
