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
	"fmt"
	"os"

	"github.com/gookit/goutil/dump"
	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/pterm/pterm"

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
			// we need to collect the cloud provider, and clusterType for the metadata to work

			meta := controller.Metadata{}

			if v, ok := k.getClusterName(); !ok {
				os.Exit(1)
			} else {
				meta.ClusterName = v
			}

			if v, ok := k.getSelectedClusterType(); !ok {
				os.Exit(1)
			} else {
				meta.ClusterType = v
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

			managerClient, err := controllerMeta.NewController(
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

			regions, err := managerClient.ListAllRegions()
			if err != nil {
				k.l.Error("Failed to sync the metadata", "Reason", err)
				os.Exit(1)
			}

			if v, ok := k.getSelectedRegion(regions); !ok {
				os.Exit(1)
			} else {
				meta.Region = v
			}
			spinner := pterm.DefaultSpinner
			// give me the snake spinnner here
			spinner.Sequence = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			ss, err := spinner.Start("Please wait, fetching the instance types")
			if err != nil {
				k.l.Error("Failed to start the spinner", "Reason", err)
				os.Exit(1)
			}

			vms, err := managerClient.ListAllInstances(meta.Region)
			if err != nil {
				ss.Fail(fmt.Sprintf("Unable to get instance_type list: %v", err))
				k.l.Error("Failed to sync the metadata", "Reason", err)
				os.Exit(1)
			}
			ss.Success("Fetched the instance type list")
			if v, ok := k.getSelectedInstanceType(vms); !ok {
				os.Exit(1)
			} else {
				meta.ManagedNodeType = v
				dump.Println(vms[v])
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

func (k *KsctlCommand) getSelectedInstanceType(vms map[string]provider.InstanceRegionOutput) (string, bool) {
	vr := make(map[string]string, len(vms))
	for sku, vm := range vms {
		if vm.CpuArch == provider.ArchAmd64 {
			displayName := fmt.Sprintf("%s (vCPUs: %d, Memory: %dGB)",
				vm.Description,
				vm.VCpus,
				vm.Memory,
			)
			cost := 0.0
			if vm.Price.HourlyPrice != nil {
				cost = *vm.Price.HourlyPrice * 730
			}
			if vm.Price.MonthlyPrice != nil { // it overrides the hourly rate if its there
				cost = *vm.Price.MonthlyPrice
			}
			displayName += fmt.Sprintf(", Price: %.2f %s/month",
				cost,
				vm.Price.Currency,
			)

			vr[displayName] = sku
		}
	}

	if v, err := cli.DropDown(
		"Select the instance type",
		vr,
		"",
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "instanceType", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedClusterType() (consts.KsctlClusterType, bool) {
	if v, err := cli.DropDown(
		"Select the cluster type",
		map[string]string{
			"Cloud Managed (For ex. EKS, AKS, GKE)":    string(consts.ClusterTypeMang),
			"Self Managed (For example, K3s, Kubeadm)": string(consts.ClusterTypeSelfMang),
		},
		string(consts.ClusterTypeMang),
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "clusterType", v)
		return consts.KsctlClusterType(v), true
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
