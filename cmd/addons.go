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
	"strconv"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"

	addonsHandler "github.com/ksctl/ksctl/v2/pkg/handler/addons"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Addons() *cobra.Command {

	cmd := &cobra.Command{
		Use: "addons",
		Example: `
ksctl addons --help
`,
		Short: "Use to work with addons",
		Long:  "It is used to work with addons",
	}

	return cmd
}

func (k *KsctlCommand) EnableAddon() *cobra.Command {

	cmd := &cobra.Command{
		Use: "enable",
		Example: `
ksctl addons enable --help
`,
		Short: "Use to enable an addon",
		Long:  "It is used to enable an addon",
		Run: func(cmd *cobra.Command, args []string) {
			m, ok := k.addonClientSetup()
			if !ok {
				os.Exit(1)
			}

			c, err := addonsHandler.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: *m,
				},
			)
			if err != nil {
				k.l.Error("Error in creating the controller", "Error", err)
				os.Exit(1)
			}

			addons, err := c.ListAllAddons()
			if err != nil {
				k.l.Error("Error in listing the addons", "Error", err)
				os.Exit(1)
			}

			addonSku, err := k.menuDriven.DropDownList(
				"Select the addon to enable",
				addons,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			addonVers, err := c.ListAvailableVersions(addonSku)
			if err != nil {
				k.l.Error("Error in listing the versions", "Error", err)
				os.Exit(1)
			}

			addonVer, err := k.menuDriven.DropDownList(
				"Select the version to enable",
				addonVers,
				cli.WithDefaultValue(addonVers[0]),
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			if cc, err := c.GetAddon(addonSku); err != nil {
				k.l.Error("Error in getting the addon", "Error", err)
				os.Exit(1)
			} else {
				if _err := cc.Install(addonVer); _err != nil {
					k.l.Error("Error in enabling the addon", "Error", _err)
					os.Exit(1)
				}
			}

			k.l.Success(k.Ctx, "Addon enabled successfully", "sku", addonSku, "version", addonVer)
		},
	}
	return cmd
}

func (k *KsctlCommand) DisableAddon() *cobra.Command {

	cmd := &cobra.Command{
		Use: "disable",
		Example: `
ksctl addons disable --help
`,
		Short: "Use to disable an addon",
		Long:  "It is used to disable an addon",
		Run: func(cmd *cobra.Command, args []string) {
			m, ok := k.addonClientSetup()
			if !ok {
				os.Exit(1)
			}

			c, err := addonsHandler.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: *m,
				},
			)
			if err != nil {
				k.l.Error("Error in creating the controller", "Error", err)
				os.Exit(1)
			}

			addons, err := c.ListInstalledAddons()
			if err != nil {
				k.l.Error("Error in listing the installed addons", "Error", err)
				os.Exit(1)
			}

			vals := make(map[string]string, len(addons))
			for _, addon := range addons {
				ver := "NaN"
				if addon.Version != "" {
					ver = "@" + addon.Version
				}
				vals[fmt.Sprintf("%s%s", addon.Name, ver)] = addon.Name
			}

			selectedAddon, err := k.menuDriven.DropDown(
				"Select the addon to disable",
				vals,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			if cc, err := c.GetAddon(selectedAddon); err != nil {
				k.l.Error("Error in getting the addon", "Error", err)
				os.Exit(1)
			} else {
				if _err := cc.Uninstall(); _err != nil {
					k.l.Error("Error in disabling the addon", "Error", _err)
					os.Exit(1)
				}
			}

			k.l.Success(k.Ctx, "Addon disabled successfully", "sku", selectedAddon)
		},
	}
	return cmd
}

func (k *KsctlCommand) addonClientSetup() (*controller.Metadata, bool) {
	clusters, err := k.fetchAllClusters()
	if err != nil {
		k.l.Error("Error in fetching the clusters", "Error", err)
		return nil, false
	}

	if len(clusters) == 0 {
		k.l.Error("No clusters found to delete")
		return nil, false
	}

	selectDisplay := make(map[string]string, len(clusters))
	valueMaping := make(map[string]controller.Metadata, len(clusters))

	for idx, cluster := range clusters {
		selectDisplay[makeHumanReadableList(cluster)] = strconv.Itoa(idx)
		valueMaping[strconv.Itoa(idx)] = controller.Metadata{
			ClusterName:   cluster.Name,
			ClusterType:   cluster.ClusterType,
			Provider:      cluster.CloudProvider,
			Region:        cluster.Region,
			StateLocation: k.KsctlConfig.PreferedStateStore,
			K8sDistro:     cluster.K8sDistro,
		}
	}

	selectedCluster, err := k.menuDriven.DropDown(
		"Select the cluster for addon operation",
		selectDisplay,
	)
	if err != nil {
		k.l.Error("Failed to get userinput", "Reason", err)
		return nil, false
	}

	m := valueMaping[selectedCluster]
	return &m, true
}
