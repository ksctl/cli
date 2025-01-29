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
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Create() *cobra.Command {

	cmd := &cobra.Command{
		Use: "create",
		Example: `
ksctl create --help
		`,
		Short: "Use to create a cluster",
		Long:  "It is used to create cluster with the given name from user",

		Run: func(cmd *cobra.Command, args []string) {
			if v, err := cli.TextInput("Enter Cluster Name"); err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			} else {
				k.l.Note(k.Ctx, "Text input", "clusterName", v)
			}

			if v, err := cli.DropDown(
				"Select the cloud provider",
				map[string]string{
					"Amazon Web Services": string(consts.CloudAws),
					"Azure":               string(consts.CloudAzure),
					"Kind":                string(consts.CloudLocal),
				},
				string(k.KsctlConfig.DefaultProvider),
			); err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			} else {
				k.l.Note(k.Ctx, "DropDown input", "cloudProvider", v)
			}

			if v, err := cli.DropDown(
				"Select the Storage Driver",
				map[string]string{
					"MongoDb": string(consts.StoreExtMongo),
					"Local":   string(consts.StoreLocal),
				},
				string(k.KsctlConfig.DefaultProvider),
			); err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			} else {
				k.l.Note(k.Ctx, "DropDown input", "storageDriver", v)
			}
		},
	}

	return cmd
}
