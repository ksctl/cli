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
	"github.com/ksctl/ksctl/v2/pkg/errors"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/ksctl/ksctl/v2/pkg/validation"
	"github.com/spf13/cobra"

	controllerCommon "github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
)

func (k *KsctlCommand) List() *cobra.Command {

	var clusterType = ""

	cmd := &cobra.Command{
		Use: "list",
		Example: `
ksctl list --help
`,
		Short: "Use to list the clusters",
		Long:  "It is used to list the clusters created by the user",
		Run: func(cmd *cobra.Command, args []string) {
			if !validation.ValidateClusterType(consts.KsctlClusterType(clusterType)) {
				k.l.Error("Invalid cluster type", "Type", clusterType)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&clusterType, "cluster-type", "", "Type of cluster to list")

	return cmd
}

func (k *KsctlCommand) ListAll() *cobra.Command {

	cmd := &cobra.Command{
		Use: "all",
		Example: `
ksctl list all --help
`,
		Short: "Use to list all the clusters",
		Long:  "It is used to list all the clusters created by the user",
		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Print(k.Ctx, "No clusters found")
				return
			}

			HandleTableOutputListAll(k.Ctx, k.l, clusters)
		},
	}

	return cmd
}

func (k *KsctlCommand) fetchAllClusters() ([]provider.ClusterData, error) {
	m := controller.Metadata{}
	if v, ok := k.getSelectedStorageDriver(); !ok {
		return nil, errors.NewError(errors.ErrInvalidStorageProvider)
	} else {
		m.StateLocation = v
	}

	m.Provider = consts.CloudAll

	managerClient, err := controllerCommon.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: m,
		},
	)
	if err != nil {
		k.l.Error("unable to initialize the ksctl manager", "Reason", err)
		return nil, err
	}

	clusters, err := managerClient.ListClusters()
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func HandleTableOutputListAll(ctx context.Context, l logger.Logger, data []provider.ClusterData) {
	headers := []string{"Name", "Type", "Cloud", "Region", "BootstrapProvider"}
	var dataToPrint [][]string = make([][]string, 0, len(data))
	for _, v := range data {
		var row []string
		row = append(row, v.Name, string(v.ClusterType), string(v.CloudProvider))
		if v.Region == "" || v.Region == "LOCAL" {
			row = append(row, "")
		} else {
			row = append(row, v.Region)
		}
		row = append(
			row,
			string(v.K8sDistro),
		)
		dataToPrint = append(dataToPrint, row)
	}

	l.Table(ctx, headers, dataToPrint)
}

func handleTableOutputGet(ctx context.Context, l logger.Logger, data provider.ClusterData) {

	headers := []string{"Attributes", "Values"}
	dataToPrint := [][]string{
		{"ClusterName", data.Name},
		{"Region", data.Region},
		{"CloudProvider", string(data.CloudProvider)},
		{"ClusterType", string(data.ClusterType)},
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
			[]string{"HaProxyVersion", data.LB.VMSize},
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
