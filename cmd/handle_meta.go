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
	"os/exec"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/addons"
	"github.com/ksctl/ksctl/v2/pkg/bootstrap/handler/cni"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/errors"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/ksctl/ksctl/v2/pkg/utilities"
	"gopkg.in/yaml.v3"
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
	// Use the new interactive cluster summary
	cli.NewBlueprintUI(os.Stdout).RenderClusterBlueprint(meta)
}

func (k *KsctlCommand) handleCNI(metaClient *controllerMeta.Controller, managedCNI addons.ClusterAddons, defaultOptionManaged string, ksctlCNI addons.ClusterAddons, defaultOptionKsctl string) (addons.ClusterAddons, error) {
	var v addons.ClusterAddons

	handleInput := func(
		vc addons.ClusterAddons,
		prompt string,
		defaultOpt string,
		errorPrompt string,
	) (addons.ClusterAddon, error) {
		cc := map[string]string{}
		cm := map[string]addons.ClusterAddon{}
		for _, c := range vc {
			cc[fmt.Sprintf("%s (From: %s)", c.Name, c.Label)] = c.Name
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

	ss := k.menuDriven.GetProgressAnimation()
	ss.Start("Fetching the CNI version list")

	config := make(map[string]map[string]any)
	if _v1.Name == string(consts.CNIFlannel) {
		vers, errF := metaClient.ListAllFlannelVersions()
		ss.Stop()
		if errF != nil { // Skip further processing if error
			k.l.Warn(k.Ctx, "Failed to get the Flannel version list", "Reason", errF)
		} else {
			if v, err := k.menuDriven.DropDownList("Select the flannel version", vers, cli.WithDefaultValue(vers[0])); err != nil {
				k.l.Error("Failed to get the flannel version", "Reason", err)
			} else {
				k.l.Debug(k.Ctx, "Selected flannel version", "Version", v)
				config[string(cni.FlannelComponentID)] = map[string]any{
					"version": v,
				}
			}
		}
	}

	if _v1.Name == string(consts.CNICilium) {
		vers, errC := metaClient.ListAllCiliumVersions()
		ss.Stop()
		if errC != nil { // Skip further processing if error
			k.l.Warn(k.Ctx, "Failed to get the Cilium version list", "Reason", errC)
		} else {
			if v, err := k.menuDriven.DropDownList("Select the cilium version", vers, cli.WithDefaultValue(vers[0])); err != nil {
				k.l.Error("Failed to get the cilium version", "Reason", err)
			} else {
				k.l.Debug(k.Ctx, "Selected cilium version", "Version", v)
				config[string(cni.CiliumComponentID)] = map[string]any{
					"version": v,
				}
			}
		}

		// Get the Cilium Specific options
		// where 2 modes are there one if guided and another is advance
		ciliumMode, err := k.menuDriven.DropDownList("Select the cilium mode", []string{"guided", "advanced"}, cli.WithDefaultValue("guided"))
		if err != nil {
			return nil, errors.WrapError(
				errors.ErrInvalidUserInput,
				k.l.NewError(k.Ctx, "Failed to get the cilium mode", "Reason", err),
			)
		}
		k.l.Print(k.Ctx, "Selected cilium mode", "Mode", ciliumMode)

		if ciliumMode == "guided" {
			availableGuidedSetup := cni.CiliumGuidedConfigurations()

			input := make(map[string]string, len(availableGuidedSetup))
			for _, v := range availableGuidedSetup {
				input[fmt.Sprintf("%s: %s", v.Name, v.Description)] = v.Name
			}

			selectedOption, err := k.menuDriven.MultiSelect("Select the cilium guided setup", input)
			if err != nil {
				return nil, errors.WrapError(
					errors.ErrInvalidUserInput,
					k.l.NewError(k.Ctx, "Failed to get the cilium guided setup", "Reason", err),
				)
			}
			k.l.Debug(k.Ctx, "Selected cilium guided setup", "Setup", selectedOption)
			config[string(cni.CiliumComponentID)]["guidedConfig"] = selectedOption

		} else {
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}

			tempFile, err := os.CreateTemp("", "cilium_config_*.yaml")
			if err != nil {
				k.l.Error("Failed to create temporary file", "Reason", err)
				return nil, errors.WrapError(
					errors.ErrInvalidUserInput,
					k.l.NewError(k.Ctx, "Failed to create temporary file", "Reason", err),
				)
			}
			defer os.Remove(tempFile.Name())

			cmd := exec.Command(editor, tempFile.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				k.l.Error("Failed to open editor", "Reason", err)
				return nil, errors.WrapError(
					errors.ErrInvalidUserInput,
					k.l.NewError(k.Ctx, "Failed to open editor", "Reason", err),
				)
			}

			content, err := os.ReadFile(tempFile.Name())
			if err != nil {
				k.l.Error("Failed to read temporary file", "Reason", err)
				return nil, errors.WrapError(
					errors.ErrInvalidUserInput,
					k.l.NewError(k.Ctx, "Failed to read temporary file", "Reason", err),
				)
			}

			var customConfig map[string]any
			if err := yaml.Unmarshal(content, &customConfig); err != nil {
				k.l.Error("Failed to parse YAML content", "Reason", err)
				return nil, errors.WrapError(
					errors.ErrInvalidUserInput,
					k.l.NewError(k.Ctx, "Failed to parse YAML content", "Reason", err),
				)
			}

			config[string(cni.CiliumComponentID)]["ciliumChartOverridings"] = customConfig
		}
	}

	_config, err := json.Marshal(config)
	if err != nil {
		return nil, errors.WrapError(
			errors.ErrInvalidUserInput,
			k.l.NewError(k.Ctx, "Failed to marshal the CNI config", "Reason", err),
		)
	}
	_v1.Config = utilities.Ptr(string(_config))

	v = append(v, _v1)
	return v, nil
}
