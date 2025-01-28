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

import "github.com/spf13/cobra"

func (k *KsctlCommand) Version() *cobra.Command {

	cmd := &cobra.Command{
		Use: "version",
		Example: `
ksctl version --help
		`,
		Short: "ksctl version",
		Long:  "To get version for ksctl components",
		Run: func(cmd *cobra.Command, args []string) {
			k.l.Box(k.Ctx, "Version", "cli: dev\nksctl: dev")
		},
	}

	return cmd
}
