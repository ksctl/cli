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
	"strconv"

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

func (k *KsctlCommand) handleRegionSelection(meta *controllerMeta.Controller, m *controller.Metadata) {
	ss := k.menuDriven.GetProgressAnimation()
	ss.Start("Fetching the region list")

	listOfRegions, err := meta.ListAllRegions()
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	if v, ok := k.getSelectedRegion(listOfRegions); !ok {
		os.Exit(1)
	} else {
		m.Region = v
	}
}

func (k *KsctlCommand) handleInstanceTypeSelection(meta *controllerMeta.Controller, m *controller.Metadata, prompt string) provider.InstanceRegionOutput {
	if len(k.inMemInstanceTypesInReg) == 0 {
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

	v, ok := k.getSelectedInstanceType(prompt, k.inMemInstanceTypesInReg)
	if !ok {
		k.l.Error("Failed to get the instance type")
		os.Exit(1)
	}
	return k.inMemInstanceTypesInReg[v]
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
	k.l.Box(k.Ctx, "Cluster Blueprint", "Here is the blueprint of the cluster")

	headers := []string{"Attributes", "Values"}

	{
		k.l.Box(k.Ctx, "Ksctl Cluster Summary", "Key attributes of the cluster")
		dd := [][]string{}
		dd = append(dd,
			[]string{"ClusterName", meta.ClusterName},
			[]string{"Region", meta.Region},
			[]string{"CloudProvider", string(meta.Provider)},
			[]string{"ClusterType", string(meta.ClusterType)},
		)

		k.l.Table(k.Ctx, headers, dd)
	}

	println()

	{
		dd := [][]string{}

		if meta.NoCP > 0 {
			dd = append(dd, []string{"ControlPlaneNodes", strconv.Itoa(meta.NoCP) + " X " + meta.ControlPlaneNodeType})
		}
		if meta.NoWP > 0 {
			dd = append(dd, []string{"WorkerPlaneNodes", strconv.Itoa(meta.NoWP) + " X " + meta.WorkerPlaneNodeType})
		}
		if meta.NoDS > 0 {
			dd = append(dd, []string{"EtcdNodes", strconv.Itoa(meta.NoDS) + " X " + meta.DataStoreNodeType})
		}
		if meta.LoadBalancerNodeType != "" {
			dd = append(dd, []string{"LoadBalancer", meta.LoadBalancerNodeType})
		}
		if len(meta.ManagedNodeType) > 0 {
			dd = append(dd, []string{"ManagedNodes", strconv.Itoa(meta.NoMP) + " X " + meta.ManagedNodeType})
		}

		if len(dd) > 0 {
			k.l.Box(k.Ctx, "Ksctl Cluster Summary", "Infrastructure details of the cluster")
			k.l.Table(k.Ctx, headers, dd)
		}
	}
	println()

	{
		dd := [][]string{}

		if meta.K8sDistro != "" {
			dd = append(dd, []string{"BootstrapProvider", string(meta.K8sDistro)})
		}
		if meta.EtcdVersion != "" {
			dd = append(dd, []string{"EtcdVersion", meta.EtcdVersion})
		}
		if meta.K8sVersion != "" {
			dd = append(dd, []string{"BootstrapKubernetesVersion", meta.K8sVersion})
		}

		if len(dd) > 0 {
			k.l.Box(k.Ctx, "Ksctl Cluster Summary", "Bootstrap details of the cluster")
			k.l.Table(k.Ctx, headers, dd)
		}

	}
	println()

	{
		// Addons Summary
		if len(meta.Addons) > 0 {
			dd := [][]string{}

			k.l.Box(k.Ctx, "Ksctl Cluster Summary", "Addons details of the cluster")

			v, _ := json.MarshalIndent(meta.Addons, "", "  ")
			dd = append(dd, []string{"Addons", string(v)})
			k.l.Table(k.Ctx, headers, dd)
			println()
		}
	}
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
