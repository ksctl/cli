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

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/errors"
	controllerCommon "github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/selfmanaged"
	"github.com/ksctl/ksctl/v2/pkg/provider"
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
			clusters, err := k.fetchSelfManagedClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Error("There is no SelfManaged cluster")
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
						K8sVersion:    cluster.K8sVersion,
						NoWP:          cluster.NoWP,
					}
				}
			}

			selectedCluster, err := k.menuDriven.DropDown(
				"Select the cluster to scaleup",
				selectDisplay,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterScaleUp, telemetry.TelemetryMeta{
				CloudProvider:     m.Provider,
				StorageDriver:     m.StateLocation,
				Region:            m.Region,
				ClusterType:       m.ClusterType,
				BootstrapProvider: m.K8sDistro,
				K8sVersion:        m.K8sVersion,
			}); err != nil {
				k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
			}

			currWP := m.NoWP

			v, ok := k.getCounterValue(
				"Enter the desired number of worker nodes",
				func(i int) bool {
					return i > currWP
				},
				currWP,
			)
			if !ok {
				k.l.Warn(k.Ctx, "Make sure the no of workernodes should be more than the current workernodes")
				os.Exit(1)
			}

			m.NoWP = v

			if err := k.loadCloudProviderCreds(m.Provider); err != nil {
				k.l.Error("Error in loading the cloud provider creds", "Error", err)
				os.Exit(1)
			}

			metaClient, err := controllerMeta.NewController(
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

			wp := k.handleInstanceTypeSelection(metaClient, &m, "Select instance_type for Worker Nodes")

			m.WorkerPlaneNodeType = wp.Sku

			curr := "$"
			if wp.Price.Currency == "USD" {
				curr = "$"
			} else if wp.Price.Currency == "INR" {
				curr = "₹"
			} else if wp.Price.Currency == "EUR" {
				curr = "€"
			}

			k.l.Box(k.Ctx, "Updated Cost", fmt.Sprintf("Cost of the cluster will +%s%.2f (%d X %s)", curr, float64(m.NoWP-currWP)*wp.GetCost(), m.NoWP-currWP, wp.Sku))

			// {
			// 	cc := m
			// 	cc.NoWP -= currWP
			// 	k.metadataSummary(cc)
			// }

			if ok, _ := k.menuDriven.Confirmation("Do you want to proceed with the cluster scaleup", cli.WithDefaultValue("no")); !ok {
				os.Exit(1)
			}

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
			clusters, err := k.fetchSelfManagedClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Error("There is no SelfManaged cluster")
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
						K8sVersion:    cluster.K8sVersion,
						NoWP:          cluster.NoWP,
						WorkerPlaneNodeType: func() string {
							g := []string{}
							for _, v := range cluster.WP {
								g = append(g, v.VMSize)
							}
							return strings.Join(g, ",")
						}(),
					}
				}
			}

			selectedCluster, err := k.menuDriven.DropDown(
				"Select the cluster to scaledown",
				selectDisplay,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterScaleDown, telemetry.TelemetryMeta{
				CloudProvider:     m.Provider,
				StorageDriver:     m.StateLocation,
				Region:            m.Region,
				ClusterType:       m.ClusterType,
				BootstrapProvider: m.K8sDistro,
				K8sVersion:        m.K8sVersion,
			}); err != nil {
				k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
			}

			if err := k.loadCloudProviderCreds(m.Provider); err != nil {
				k.l.Error("Error in loading the cloud provider creds", "Error", err)
				os.Exit(1)
			}

			currWP := m.NoWP
			if currWP == 0 {
				k.l.Error("There is no worker node to scale down")
				os.Exit(1)
			}

			v, ok := k.getCounterValue(
				"Enter the desired number of worker nodes",
				func(i int) bool {
					return i < currWP && i >= 0
				},
				currWP,
			)
			if !ok {
				k.l.Warn(k.Ctx, "Make sure the no of workernodes should be less than the current workernodes and not less than 0")
				os.Exit(1)
			}

			m.NoWP = v

			{
				// for just showing the costs changes
				cc := m

				metaClient, err := controllerMeta.NewController(
					k.Ctx,
					k.l,
					&controller.Client{
						Metadata: cc,
					},
				)
				if err != nil {
					k.l.Error("Failed to create the controller", "Reason", err)
					os.Exit(1)
				}

				vms := strings.Split(cc.WorkerPlaneNodeType, ",")

				g := map[string]struct {
					Count int
					VM    provider.InstanceRegionOutput
				}{}

				for i := cc.NoWP; i < len(vms); i++ {
					vm := vms[i]

					wp := k.getSpecificInstance(metaClient, cc.Region, vm)
					if _, ok := g[wp.Sku]; ok {
						g[wp.Sku] = struct {
							Count int
							VM    provider.InstanceRegionOutput
						}{
							Count: g[wp.Sku].Count + 1,
							VM:    wp,
						}
					} else {
						g[wp.Sku] = struct {
							Count int
							VM    provider.InstanceRegionOutput
						}{
							Count: 1,
							VM:    wp,
						}
					}
				}

				curr := "$"
				total := 0.0
				vmSize := []string{}
				for k, x := range g {
					if x.VM.Price.Currency == "USD" {
						curr = "$"
					} else if x.VM.Price.Currency == "INR" {
						curr = "₹"
					} else if x.VM.Price.Currency == "EUR" {
						curr = "€"
					}
					total += float64(x.Count) * x.VM.GetCost()
					vmSize = append(vmSize, fmt.Sprintf("(%d X %s)", x.Count, k))
				}

				k.l.Box(k.Ctx, "Updated Cost", fmt.Sprintf("Cost of the cluster will -%s%.2f <%s>", curr, total, strings.Join(vmSize, ",")))

				cc.NoWP -= currWP

				cc.WorkerPlaneNodeType = strings.Join(vmSize, ",")

				// k.metadataSummary(cc)
			}

			if ok, _ := k.menuDriven.Confirmation("Do you want to proceed with the cluster scaledown", cli.WithDefaultValue("no")); !ok {
				os.Exit(1)
			}

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

func (k *KsctlCommand) fetchSelfManagedClusters() ([]provider.ClusterData, error) {
	m := controller.Metadata{}
	if v, ok := k.getSelectedStorageDriver(); !ok {
		return nil, errors.NewError(errors.ErrInvalidStorageProvider)
	} else {
		m.StateLocation = v
	}

	m.Provider = consts.CloudAll
	m.ClusterType = consts.ClusterTypeSelfMang

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
