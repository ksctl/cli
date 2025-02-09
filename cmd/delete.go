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
	"strings"

	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/managed"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/selfmanaged"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Delete() *cobra.Command {

	cmd := &cobra.Command{
		Use: "delete",
		Example: `
ksctl delete --help
		`,
		Short: "Use to delete a cluster",
		Long:  "It is used to delete cluster with the given name from user",

		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Error("No clusters found to delete")
				os.Exit(1)
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

			selectedCluster, err := cli.DropDown(
				"Select the cluster to delete",
				selectDisplay,
				"",
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			if ok, _ := cli.Confirmation("Do you want to proceed with the cluster deletion", "no"); !ok {
				os.Exit(1)
			}

			if k.loadCloudProviderCreds(m.Provider) != nil {
				os.Exit(1)
			}

			if m.ClusterType == consts.ClusterTypeMang {
				c, err := managed.NewController(
					k.Ctx,
					k.l,
					&controller.Client{
						Metadata: m,
					},
				)
				if err != nil {
					k.l.Error("Failed to create the controller", "Reason", err)
					os.Exit(1)
				}

				if err := c.Delete(); err != nil {
					k.l.Error("Failed to delete your managed cluster", "Reason", err)
					os.Exit(1)
				}

			} else {
				c, err := selfmanaged.NewController(
					k.Ctx,
					k.l,
					&controller.Client{
						Metadata: m,
					},
				)
				if err != nil {
					k.l.Error("Failed to create the controller", "Reason", err)
					os.Exit(1)
				}

				if err := c.Delete(); err != nil {
					k.l.Error("Failed to delete your selfmanaged cluster", "Reason", err)
					os.Exit(1)
				}
			}

			k.l.Success(k.Ctx, "Deleted your cluster", "Name", m.ClusterName)
		},
	}

	return cmd
}

func makeHumanReadableList(m provider.ClusterData) string {
	fields := []string{"%s", "[%s]", "=>"}
	args := []any{m.Name, m.ClusterType}
	if m.CloudProvider == consts.CloudLocal {
		fields = append(fields, "local")
	} else {
		fields = append(fields, "%s ⟨%s⟩")
		args = append(args, m.CloudProvider, m.Region)
	}

	return fmt.Sprintf(strings.Join(fields, " "), args...)
}
