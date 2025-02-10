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
	"strconv"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/selfmanaged"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) ScaleUp() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaleup",
		Example: `
ksctl update scaleup --help
		`,
		Short: "Use to manually scaleup a selfmanaged cluster",
		Long:  "It is used to manually scaleup a selfmanaged cluster",

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
				if cluster.ClusterType == consts.ClusterTypeSelfMang {
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
			}

			selectedCluster, err := cli.DropDown(
				"Select the cluster to scaleup",
				selectDisplay,
				"",
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			c, err := selfmanaged.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: m,
				},
			)
			if err != nil {
				k.l.Error("Error in creating the controller", "Error", err)
				os.Exit(1)
			}

			if err := c.AddWorkerNodes(); err != nil {
				k.l.Error("Error in scaling up the cluster", "Error", err)
				os.Exit(1)
			}

			k.l.Success(k.Ctx, "Cluster workernode scaled up successfully")
		},
	}
	return cmd
}

func (k *KsctlCommand) ScaleDown() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaledown",
		Example: `
ksctl update scaledown --help
		`,
		Short: "Use to manually scaledown a selfmanaged cluster",
		Long:  "It is used to manually scaledown a selfmanaged cluster",

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
				if cluster.ClusterType == consts.ClusterTypeSelfMang {
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
			}

			selectedCluster, err := cli.DropDown(
				"Select the cluster to scaledown",
				selectDisplay,
				"",
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			c, err := selfmanaged.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: m,
				},
			)
			if err != nil {
				k.l.Error("Error in creating the controller", "Error", err)
				os.Exit(1)
			}

			if err := c.DeleteWorkerNodes(); err != nil {
				k.l.Error("Error in scaling down the cluster", "Error", err)
				os.Exit(1)
			}

			k.l.Success(k.Ctx, "Cluster workernode scaled down successfully")
		},
	}
	return cmd
}
