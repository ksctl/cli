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

	"github.com/gookit/goutil/dump"
	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	controllerCommon "github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/provider"

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

			meta := controller.Metadata{}

			if v, ok := k.getClusterName(); !ok {
				os.Exit(1)
			} else {
				meta.ClusterName = v
			}

			if v, ok := k.getSelectedCloudProvider(); !ok {
				os.Exit(1)
			} else {
				meta.Provider = v
			}

			if v, ok := k.getSelectedStorageDriver(); !ok {
				os.Exit(1)
			} else {
				k.l.Debug(k.Ctx, "DropDown input", "storageDriver", v)
				meta.StateLocation = consts.KsctlStore(v)
			}

			managerClient, err := controllerCommon.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: meta,
				},
			)
			if err != nil {
				k.l.Error("Failed to create the controller", "Reason", err)
				os.Exit(1)
			}

			regions, err := managerClient.SyncMetadata()
			if err != nil {
				k.l.Error("Failed to sync the metadata", "Reason", err)
				os.Exit(1)
			}

			if v, ok := k.getSelectedRegion(regions); !ok {
				os.Exit(1)
			} else {
				meta.Region = v
			}

			dump.Println(meta)

		},
	}

	return cmd
}

func (k *KsctlCommand) getClusterName() (string, bool) {
	v, err := cli.TextInput("Enter Cluster Name")
	if err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	}
	if len(v) == 0 {
		k.l.Error("Cluster name cannot be empty")
		return "", false
	}
	k.l.Debug(k.Ctx, "Text input", "clusterName", v)
	return v, true
}

func (k *KsctlCommand) getSelectedRegion(regions []provider.RegionOutput) (string, bool) {
	vr := make(map[string]string, len(regions))
	for _, r := range regions {
		vr[r.Name] = r.Sku
	}
	dump.Println(vr)

	if v, err := cli.DropDown(
		"Select the region",
		vr,
		"",
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "region", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedCloudProvider() (consts.KsctlCloud, bool) {
	if v, err := cli.DropDown(
		"Select the cloud provider",
		map[string]string{
			"Amazon Web Services": string(consts.CloudAws),
			"Azure":               string(consts.CloudAzure),
			"Kind":                string(consts.CloudLocal),
		},
		"",
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "cloudProvider", v)

		switch consts.KsctlCloud(v) {
		case consts.CloudAws:
			if errC := k.loadAwsCredentials(); errC != nil {
				k.l.Error("Failed to load the AWS credentials", "Reason", errC)
				return "", false
			}
		case consts.CloudAzure:
			if errC := k.loadAzureCredentials(); errC != nil {
				k.l.Error("Failed to load the Azure credentials", "Reason", errC)
				return "", false
			}
		}

		return consts.KsctlCloud(v), true
	}
}

func (k *KsctlCommand) getSelectedStorageDriver() (consts.KsctlStore, bool) {
	if v, err := cli.DropDown(
		"Select the Storage Driver",
		map[string]string{
			"MongoDb": string(consts.StoreExtMongo),
			"Local":   string(consts.StoreLocal),
		},
		string(k.KsctlConfig.PreferedStateStore),
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "storageDriver", v)
		if errS := k.loadMongoCredentials(); errS != nil {
			k.l.Error("Failed to load the MongoDB credentials", "Reason", errS)
			return "", false
		}

		return consts.KsctlStore(v), true
	}
}
