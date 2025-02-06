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
	"strconv"

	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

func (k *KsctlCommand) getClusterName() (string, bool) {
	v, err := cli.TextInput("Enter Cluster Name", "")
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
	v, err := cli.DropDown(
		"Select the bootstrap type",
		map[string]string{
			"Kubeadm": string(consts.K8sKubeadm),
			"K3s":     string(consts.K8sK3s),
		},
		string(consts.K8sK3s),
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
	v, err := cli.TextInput(prompt, strconv.Itoa(defaultVal))
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

func (k *KsctlCommand) getSelectedRegion(regions []provider.RegionOutput) (string, bool) {
	vr := make(map[string]string, len(regions))
	for _, r := range regions {
		vr[r.Name] = r.Sku
	}

	k.l.Debug(k.Ctx, "Regions", "regions", vr)

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

func (k *KsctlCommand) getSelectedK8sVersion(prompt string, vers []string) (string, bool) {
	k.l.Debug(k.Ctx, "List of k8s versions", "versions", vers)

	if v, err := cli.DropDownList(
		prompt,
		vers,
		vers[0],
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

	if v, err := cli.DropDown(
		prompt,
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

	if v, err := cli.DropDown(
		prompt,
		vr,
		"",
	); err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return "", false
	} else {
		k.l.Debug(k.Ctx, "DropDown input", "managedClusterOffering", v)
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

func (k *KsctlCommand) getSelectedCloudProvider(v consts.KsctlClusterType) (consts.KsctlCloud, bool) {
	options := map[string]string{
		"Amazon Web Services": string(consts.CloudAws),
		"Azure":               string(consts.CloudAzure),
	}

	if v == consts.ClusterTypeMang {
		options["Kind"] = string(consts.CloudLocal)
	}

	if v, err := cli.DropDown(
		"Select the cloud provider",
		options,
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
