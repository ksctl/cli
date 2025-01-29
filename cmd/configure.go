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

	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/cli/pkg/config"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Configure() *cobra.Command {

	cmd := &cobra.Command{
		Use: "configure",

		Short: "Configure ksctl cli",
		Long:  "It will help you to configure the ksctl cli",
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := cli.DropDown(
				"What should be your default storageDriver?",
				map[string]string{
					"MongoDb": string(consts.StoreExtMongo),
					"Local":   string(consts.StoreLocal),
				},
				"Local",
			); err != nil {
				k.l.Error("Failed to get the storageDriver", "Reason", err)
				os.Exit(1)
			} else {
				k.l.Note(k.Ctx, "DropDown", "selected", v)
				k.KsctlConfig.PreferedStateStore = consts.KsctlStore(v)
				_ = config.SaveConfig(k.KsctlConfig)
			}

			v, err := cli.Confirmation("Do you want to add/modify credentials", "no")
			if err != nil {
				k.l.Error("Failed to get the confirmation", "Reason", err)
				os.Exit(1)
			}
			k.l.Note(k.Ctx, "Confirmation", "selected", v)
			if v {
				if v, err := cli.DropDown(
					"Credentials",
					map[string]string{
						"Amazon Web Services": string(consts.CloudAws),
						"Azure":               string(consts.CloudAzure),
					},
					"",
				); err != nil {
					k.l.Error("Failed to get the credentials", "Reason", err)
					os.Exit(1)
				} else {
					k.l.Note(k.Ctx, "DropDown", "selected", v)
				}
			}
		},
	}

	return cmd
}
