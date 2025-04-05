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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

func (k *KsctlCommand) getClusterName() (string, bool) {
	v, err := k.menuDriven.TextInput("Enter Cluster Name")
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

func (k *KsctlCommand) getBootstrap() (consts.KsctlKubernetes, bool) {
	v, err := k.menuDriven.DropDown(
		"Select the bootstrap type",
		map[string]string{
			"Kubeadm": string(consts.K8sKubeadm),
			"K3s":     string(consts.K8sK3s),
		},
		cli.WithDefaultValue(string(consts.K8sK3s)),
	)
	if err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	}
	k.l.Debug(k.Ctx, "DropDown input", "bootstrapType", v)
	return consts.KsctlKubernetes(v), true
}

type userInputValidation func(int) bool

func (k *KsctlCommand) getCounterValue(prompt string, validate userInputValidation, defaultVal int) (int, bool) {
	v, err := k.menuDriven.TextInput(prompt, cli.WithDefaultValue(strconv.Itoa(defaultVal)))
	if err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return 0, false
	}
	_v, err := strconv.Atoi(v)
	if err != nil {
		k.l.Error("Invalid input", "Reason", err)
		return 0, false
	}

	if !validate(_v) {
		k.l.Error("Invalid input")
		return 0, false
	}
	k.l.Debug(k.Ctx, "Text input", "counterValue", v)
	return _v, true
}

func (k *KsctlCommand) getSelectedRegion(regions provider.RegionsOutput) (string, bool) {
	k.l.Debug(k.Ctx, "Regions", "regions", regions)

	if v, err := k.menuDriven.DropDown(
		"Select the region",
		regions.S(),
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "region", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedInstanceCategory(categories map[string]provider.MachineCategory) (provider.MachineCategory, bool) {
	k.l.Debug(k.Ctx, "Instance categories", "categories", categories)

	vr := make(map[string]string, len(categories))

	for k, _v := range categories {
		useCases := strings.Join(_v.UseCases(), ", ")
		key := fmt.Sprintf("%s\n   Used for: %s\n", k, useCases)
		vr[key] = string(_v)
	}

	if v, err := k.menuDriven.DropDown(
		"Let us know about your workload type",
		vr,
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "instanceCategory", v)
		return provider.MachineCategory(v), true
	}
}

func (k *KsctlCommand) getSelectedK8sVersion(prompt string, vers []string) (string, bool) {
	k.l.Debug(k.Ctx, "List of k8s versions", "versions", vers)

	if v, err := k.menuDriven.DropDownList(
		prompt,
		vers,
		cli.WithDefaultValue(vers[0]),
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "k8sVersion", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedInstanceType(
	prompt string,
	vms map[string]provider.InstanceRegionOutput,
) (string, bool) {
	vr := make(map[string]string, len(vms))
	for sku, vm := range vms {
		if vm.CpuArch == provider.ArchAmd64 {
			displayName := fmt.Sprintf("%s (vCPUs: %d, Memory: %dGB)",
				vm.Description,
				vm.VCpus,
				vm.Memory,
			)
			displayName += fmt.Sprintf(", Price: %.2f %s/month",
				vm.GetCost(),
				vm.Price.Currency,
			)

			vr[displayName] = sku
		}
	}

	k.l.Debug(k.Ctx, "Instance types", "vms", vr)

	if v, err := k.menuDriven.DropDown(
		prompt,
		vr,
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "instanceType", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedManagedClusterOffering(
	prompt string,
	offerings map[string]provider.ManagedClusterOutput,
) (string, bool) {
	vr := make(map[string]string, len(offerings))
	for _, o := range offerings {
		displayName := fmt.Sprintf("%s, Price: %.2f %s/month",
			o.Description,
			o.GetCost(),
			o.Price.Currency,
		)

		vr[displayName] = o.Sku
	}

	k.l.Debug(k.Ctx, "Offerings", "offerings", vr)

	if v, err := k.menuDriven.DropDown(
		prompt,
		vr,
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "managedClusterOffering", v)
		return v, true
	}
}

func (k *KsctlCommand) getSelectedClusterType() (consts.KsctlClusterType, bool) {
	if v, err := k.menuDriven.DropDown(
		"Select the cluster type",
		map[string]string{
			"Cloud Managed (For ex. EKS, AKS, Kind)":   string(consts.ClusterTypeMang),
			"Self Managed (For example, K3s, Kubeadm)": string(consts.ClusterTypeSelfMang),
		},
		cli.WithDefaultValue(string(consts.ClusterTypeMang)),
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "clusterType", v)
		return consts.KsctlClusterType(v), true
	}
}

func (k *KsctlCommand) getSelectedCloudProvider(v consts.KsctlClusterType) (consts.KsctlCloud, bool) {
	options := map[string]string{
		"Amazon Web Services": string(consts.CloudAws),
		"Azure":               string(consts.CloudAzure),
	}

	if v == consts.ClusterTypeMang {
		options["Kind"] = string(consts.CloudLocal)
	}

	if v, err := k.menuDriven.DropDown(
		"Select the cloud provider",
		options,
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "cloudProvider", v)

		if err := k.loadCloudProviderCreds(consts.KsctlCloud(v)); err != nil {
			return "", false
		}

		return consts.KsctlCloud(v), true
	}
}

func (k *KsctlCommand) loadCloudProviderCreds(v consts.KsctlCloud) error {
	switch v {
	case consts.CloudAws:
		if v, err := k.loadAwsCredentials(); err != nil {
			k.l.Error("Failed to load the AWS credentials", "Reason", err)
			return err
		} else {
			k.Ctx = context.WithValue(k.Ctx, consts.KsctlAwsCredentials, v)
		}

	case consts.CloudAzure:
		if v, err := k.loadAzureCredentials(); err != nil {
			k.l.Error("Failed to load the Azure credentials", "Reason", err)
			return err
		} else {
			k.Ctx = context.WithValue(k.Ctx, consts.KsctlAzureCredentials, v)
		}
	}
	return nil
}

func (k *KsctlCommand) getSelectedStorageDriver() (consts.KsctlStore, bool) {
	if k.KsctlConfig.PreferedStateStore != consts.StoreExtMongo && k.KsctlConfig.PreferedStateStore != consts.StoreLocal {
		k.l.Error("Failed to determine StorageDriver", "message", "Please use $ksctl configure to set the storage driver", "currentSetValue", k.KsctlConfig.PreferedStateStore)
		return "", false
	}

	if k.KsctlConfig.PreferedStateStore == consts.StoreExtMongo {
		if errS := k.loadMongoCredentials(); errS != nil {
			k.l.Error("Failed to load the MongoDB credentials", "Reason", errS)
			return "", false
		}
	}

	return k.KsctlConfig.PreferedStateStore, true
}
