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
	"os"
	"strconv"
	"strings"

	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Get() *cobra.Command {

	cmd := &cobra.Command{
		Use: "get",
		Example: `
ksctl get --help
`,
		Short: "Use to get the cluster",
		Long:  "It is used to get the cluster created by the user",
		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Print(k.Ctx, "No clusters found")
				os.Exit(1)
			}

			selectDisplay := make(map[string]string, len(clusters))
			valueMaping := make(map[string]provider.ClusterData, len(clusters))

			for idx, cluster := range clusters {
				selectDisplay[makeHumanReadableList(cluster)] = strconv.Itoa(idx)
				valueMaping[strconv.Itoa(idx)] = cluster
			}

			selectedCluster, err := k.menuDriven.DropDown(
				"Select the cluster to delete",
				selectDisplay,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			cluster := valueMaping[selectedCluster]

			handleTableOutputGet(k.Ctx, k.l, cluster)

		},
	}

	return cmd
}

// TODO: get the addons in the cluster

func handleTableOutputGet(ctx context.Context, l logger.Logger, data provider.ClusterData) {

	headers := []string{"Attributes", "Values"}
	dataToPrint := [][]string{
		{"ClusterName", data.Name},
		{"CloudProvider", string(data.CloudProvider)},
		{"ClusterType", string(data.ClusterType)},
	}
	if data.CloudProvider != consts.CloudLocal {
		dataToPrint = append(dataToPrint,
			[]string{"Region", data.Region},
		)
	}

	if data.ClusterType == consts.ClusterTypeSelfMang {
		nodes := func(vm []provider.VMData) string {
			slice := make([]string, 0, len(vm))
			for _, v := range vm {
				slice = append(slice, v.VMSize)
			}
			return strings.Join(slice, ",")
		}

		dataToPrint = append(
			dataToPrint,
			[]string{"BootstrapProvider", string(data.K8sDistro)},
			[]string{"BootstrapKubernetesVersion", data.K8sVersion},
			[]string{"ControlPlaneNodes", nodes(data.CP)},
			[]string{"WorkerPlaneNodes", nodes(data.WP)},
			[]string{"EtcdNodes", nodes(data.DS)},
			[]string{"LoadBalancer", data.LB.VMSize},
			[]string{"EtcdVersion", data.EtcdVersion},
			[]string{"HaProxyVersion", data.HAProxyVersion},
		)
	} else {
		dataToPrint = append(
			dataToPrint,
			[]string{"ManagedNodes", strconv.Itoa(data.NoMgt) + " X " + data.Mgt.VMSize},
			[]string{"ManagedK8sVersion", data.K8sVersion},
		)
	}

	dataToPrint = append(dataToPrint,
		[]string{"Addons", strings.Join(data.Apps, ",")},
		[]string{"CNI", data.Cni},
	)

	l.Table(ctx, headers, dataToPrint)
}
