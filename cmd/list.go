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

	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/errors"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/spf13/cobra"

	controllerCommon "github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
)

func (k *KsctlCommand) List() *cobra.Command {

	cmd := &cobra.Command{
		Use: "list",
		Example: `
ksctl list --help
`,
		Short: "Use to list all the clusters",
		Long:  "It is used to list all the clusters created by the user",
		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterList, telemetry.TelemetryMeta{}); err != nil {
				k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
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
		if v.CloudProvider == consts.CloudLocal {
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
