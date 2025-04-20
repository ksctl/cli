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
	"encoding/json"
	"fmt"
	"os"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/addons"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/errors"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

func (k *KsctlCommand) baseMetadataFields(m *controller.Metadata) {
	if v, ok := k.getClusterName(); !ok {
		os.Exit(1)
	} else {
		m.ClusterName = v
	}

	if v, ok := k.getSelectedClusterType(); !ok {
		os.Exit(1)
	} else {
		m.ClusterType = v
	}

	if v, ok := k.getSelectedCloudProvider(m.ClusterType); !ok {
		os.Exit(1)
	} else {
		m.Provider = v
	}

	if v, ok := k.getSelectedStorageDriver(); !ok {
		os.Exit(1)
	} else {
		m.StateLocation = consts.KsctlStore(v)
	}

	if m.ClusterType == consts.ClusterTypeSelfMang {
		if v, ok := k.getBootstrap(); ok {
			m.K8sDistro = v
		} else {
			os.Exit(1)
		}
	}
}

func (k *KsctlCommand) handleRegionSelection(meta *controllerMeta.Controller, m *controller.Metadata) []provider.RegionOutput {
	ss := k.menuDriven.GetProgressAnimation()
	ss.Start("Fetching the region list")

	listOfRegions, err := meta.ListAllRegions()
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	k.l.Note(k.Ctx, "Carbon emission data shown represents monthly averages calculated over a one-year period")
	k.l.Note(k.Ctx, "Select the region for the cluster")

	if v, err := k.menuDriven.CardSelection(
		cli.ConverterForRegionOutputForCards(listOfRegions),
	); err != nil {
		k.l.Error("Failed to get the region", "Reason", err)
		os.Exit(1)
	} else {
		if v == "" {
			k.l.Error("Region not selected")
			os.Exit(1)
		}
		k.l.Debug(k.Ctx, "Selected region", "Region", v)
		m.Region = v
	}

	return listOfRegions
}

func (k *KsctlCommand) handleInstanceCategorySelection() provider.MachineCategory {
	v := provider.GetAvailableMachineCategories()

	_v, ok := k.getSelectedInstanceCategory(v)
	if !ok {
		k.l.Error("Failed to get the instance category")
		os.Exit(1)
	}
	return _v
}

func (k *KsctlCommand) handleInstanceTypeSelection(
	meta *controllerMeta.Controller,
	m *controller.Metadata,
	category provider.MachineCategory,
	prompt string,
) provider.InstanceRegionOutput {

	if len(k.inMemInstanceTypesInReg) == 0 {
		if len(category) == 0 {
			k.l.Error("Machine category is not provided")
			os.Exit(1)
		}
		ss := k.menuDriven.GetProgressAnimation()
		ss.Start("Fetching the instance type list")

		listOfVMs, err := meta.ListAllInstances(m.Region)
		if err != nil {
			ss.Stop()
			k.l.Error("Failed to sync the metadata", "Reason", err)
			os.Exit(1)
		}
		ss.Stop()
		k.inMemInstanceTypesInReg = listOfVMs
	}

	availableOptions := make(provider.InstancesRegionOutput, 0, len(k.inMemInstanceTypesInReg))

	k.l.Note(k.Ctx, prompt)

	for _, v := range k.inMemInstanceTypesInReg {
		if v.Category == category && v.CpuArch == provider.ArchAmd64 {
			availableOptions = append(availableOptions, v)
		}
	}

	v, err := k.menuDriven.CardSelection(
		cli.ConverterForInstanceTypesForCards(availableOptions),
	)
	if err != nil {
		k.l.Error("Failed to get the instance type from user", "Reason", err)
		os.Exit(1)
	}
	if v == "" {
		k.l.Error("Instance type not selected")
		os.Exit(1)
	}

	_v, ok := availableOptions.Get(v)
	if !ok {
		k.l.Error("Failed to get the instance type")
		os.Exit(1)
	}

	return *_v
}

func (k *KsctlCommand) getSpecificInstanceForScaledown(
	meta *controllerMeta.Controller,
	region string,
	instanceSku string,
) provider.InstanceRegionOutput {

	if len(k.inMemInstanceTypesInReg) == 0 {
		ss := k.menuDriven.GetProgressAnimation()
		ss.Start("Fetching the instance type list")

		listOfVMs, err := meta.ListAllInstances(region)
		if err != nil {
			ss.Stop()
			k.l.Error("Failed to sync the metadata", "Reason", err)
			os.Exit(1)
		}
		ss.Stop()
		k.inMemInstanceTypesInReg = listOfVMs
	}

	v, ok := k.inMemInstanceTypesInReg.Get(instanceSku)
	if !ok {
		k.l.Error("Failed to get the instance type")
		os.Exit(1)
	}
	return *v
}

func (k *KsctlCommand) handleManagedK8sVersion(meta *controllerMeta.Controller, m *controller.Metadata) {
	ss := k.menuDriven.GetProgressAnimation()
	ss.Start("Fetching the managed cluster k8s versions")

	listOfK8sVersions, err := meta.ListAllManagedClusterK8sVersions(m.Region)
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	if v, ok := k.getSelectedK8sVersion("Select the k8s version for Managed Cluster", listOfK8sVersions); !ok {
		k.l.Error("Failed to get the k8s version")
		os.Exit(1)
	} else {
		m.K8sVersion = v
	}
}

func (k *KsctlCommand) metadataSummary(meta controller.Metadata) {
	// Add more space before the first box for better visual separation
	println()

	k.l.Box(k.Ctx, "✨ Cluster Blueprint", "Summary of your planned Kubernetes cluster")

	// Add consistent spacing between sections
	println()
	println()

	// Main cluster attributes section
	{
		k.l.Box(k.Ctx, "🔑 Key Attributes", "Core configuration details")

		headers := []string{"Attribute", "Value"}
		dd := [][]string{
			{"🔖 Cluster Name", meta.ClusterName},
			{"📍 Region", meta.Region},
			{"☁️  Cloud Provider", string(meta.Provider)},
			{"🏗️  Cluster Type", string(meta.ClusterType)},
		}

		k.l.Table(k.Ctx, headers, dd)
	}

	// Add consistent spacing between sections
	println()
	println()

	// Infrastructure details section
	{
		if meta.NoCP > 0 || meta.NoWP > 0 || meta.NoDS > 0 || len(meta.ManagedNodeType) > 0 {
			k.l.Box(k.Ctx, "🖥️  Infrastructure", "Node configuration and topology")

			headers := []string{"Component", "Specification"}
			dd := [][]string{}

			if meta.NoCP > 0 {
				dd = append(dd, []string{"🎮 Control Plane", fmt.Sprintf("%d × %s", meta.NoCP, meta.ControlPlaneNodeType)})
			}
			if meta.NoWP > 0 {
				dd = append(dd, []string{"🔋 Worker Nodes", fmt.Sprintf("%d × %s", meta.NoWP, meta.WorkerPlaneNodeType)})
			}
			if meta.NoDS > 0 {
				dd = append(dd, []string{"💾 Etcd Nodes", fmt.Sprintf("%d × %s", meta.NoDS, meta.DataStoreNodeType)})
			}
			if meta.LoadBalancerNodeType != "" {
				dd = append(dd, []string{"⚖️  Load Balancer", meta.LoadBalancerNodeType})
			}
			if len(meta.ManagedNodeType) > 0 {
				dd = append(dd, []string{"🌐 Managed Nodes", fmt.Sprintf("%d × %s", meta.NoMP, meta.ManagedNodeType)})
			}

			k.l.Table(k.Ctx, headers, dd)
		}
	}

	// Add consistent spacing between sections
	println()
	println()

	// Kubernetes configuration section
	{
		if meta.K8sDistro != "" || meta.EtcdVersion != "" || meta.K8sVersion != "" {
			k.l.Box(k.Ctx, "⚙️  Kubernetes Configuration", "Software versions and distributions")

			headers := []string{"Component", "Version/Type"}
			dd := [][]string{}

			if meta.K8sDistro != "" {
				dd = append(dd, []string{"🚀 Bootstrap Provider", string(meta.K8sDistro)})
			}
			if meta.K8sVersion != "" {
				dd = append(dd, []string{"🔄 Kubernetes Version", meta.K8sVersion})
			}
			if meta.EtcdVersion != "" {
				dd = append(dd, []string{"📦 Etcd Version", meta.EtcdVersion})
			}

			k.l.Table(k.Ctx, headers, dd)
		}
	}

	// Add consistent spacing between sections
	println()
	println()

	// Addons section
	{
		if len(meta.Addons) > 0 {
			k.l.Box(k.Ctx, "🧩 Cluster Add-ons", "Additional components to be installed")

			headers := []string{"Add-on", "Details"}
			dd := [][]string{}

			for _, addon := range meta.Addons {
				v := struct {
					Label  string
					IsCNI  bool
					Config *string
				}{
					Label:  addon.Label,
					IsCNI:  addon.IsCNI,
					Config: addon.Config,
				}

				b, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					k.l.Error("Failed to marshal addon config", "Reason", err)
					os.Exit(1)
				}
				dd = append(dd, []string{addon.Name, string(b)})
			}

			k.l.Table(k.Ctx, headers, dd)
		}
	}

	// End summary note with extra spacing
	println()
	println()
	k.l.Note(k.Ctx, "Your cluster will be provisioned with these specifications")
	println()
}

func (k *KsctlCommand) handleCNI(managedCNI addons.ClusterAddons, defaultOptionManaged string, ksctlCNI addons.ClusterAddons, defaultOptionKsctl string) (addons.ClusterAddons, error) {
	var v addons.ClusterAddons

	handleInput := func(vc addons.ClusterAddons, prompt string, defaultOpt string, errorPrompt string) (addons.ClusterAddon, error) {
		cc := map[string]string{}
		cm := map[string]addons.ClusterAddon{}
		for _, c := range vc {
			cc[fmt.Sprintf("%s <%s>", c.Name, c.Label)] = c.Name
			cm[c.Name] = c
		}

		selected, err := k.menuDriven.DropDown(
			prompt,
			cc,
			cli.WithDefaultValue(defaultOpt),
		)
		if err != nil {
			return addons.ClusterAddon{}, errors.WrapError(
				errors.ErrInvalidUserInput,
				k.l.NewError(k.Ctx, errorPrompt, "Reason", err),
			)
		}

		return cm[selected], nil
	}

	_v0, err := handleInput(managedCNI, "Select the CNI addon provided by offering", defaultOptionManaged, "Failed to get the CNI addon provided by managed offering")
	if err != nil {
		return nil, err
	}

	v = append(v, _v0)

	if _v0.Name != string(consts.CNINone) {
		return v, nil
	}

	_v1, err := handleInput(ksctlCNI, "Select the CNI addon provided by ksctl", defaultOptionKsctl, "Failed to get the CNI addon provided by ksctl")
	if err != nil {
		return nil, err
	}

	v = append(v, _v1)

	return v, nil
}
