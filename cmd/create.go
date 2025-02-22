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

	"github.com/ksctl/ksctl/v2/pkg/consts"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"

	controllerManaged "github.com/ksctl/ksctl/v2/pkg/handler/cluster/managed"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	controllerSelfManaged "github.com/ksctl/ksctl/v2/pkg/handler/cluster/selfmanaged"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Create() *cobra.Command {

	cmd := &cobra.Command{
		Use: "create",
		Example: `
ksctl create --help
		`,
		Short: "Use to create a cluster",
		Long:  "It is used to create cluster with the given name from user",

		Run: func(cmd *cobra.Command, args []string) {
			meta := controller.Metadata{}

			k.baseMetadataFields(&meta)

			if meta.ClusterType == consts.ClusterTypeMang {
				k.metadataForManagedCluster(&meta)
			} else {
				k.metadataForSelfManagedCluster(&meta)
			}

			k.l.Success(k.Ctx, "Created the cluster", "Name", meta.ClusterName)
		},
	}

	return cmd
}

func (k *KsctlCommand) metadataForSelfManagedCluster(
	meta *controller.Metadata,
) {
	metaClient, err := controllerMeta.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: *meta,
		},
	)
	if err != nil {
		k.l.Error("Failed to create the controller", "Reason", err)
		os.Exit(1)
	}

	k.handleRegionSelection(metaClient, meta)

	cp := k.handleInstanceTypeSelection(metaClient, meta, "Select instance_type for Control Plane")
	wp := k.handleInstanceTypeSelection(metaClient, meta, "Select instance_type for Worker Nodes")
	etcd := k.handleInstanceTypeSelection(metaClient, meta, "Select instance_type for Etcd Nodes")
	lb := k.handleInstanceTypeSelection(metaClient, meta, "Select instance_type for Load Balancer")

	meta.ControlPlaneNodeType = cp.Sku
	meta.WorkerPlaneNodeType = wp.Sku
	meta.DataStoreNodeType = etcd.Sku
	meta.LoadBalancerNodeType = lb.Sku

	if v, ok := k.getCounterValue("Enter the number of Control Plane Nodes", func(v int) bool {
		return v >= 3
	}, 3); !ok {
		k.l.Error("Failed to get the number of control plane nodes")
		os.Exit(1)
	} else {
		meta.NoCP = v
	}

	if v, ok := k.getCounterValue("Enter the number of Worker Nodes", func(v int) bool {
		return v > 0
	}, 1); !ok {
		k.l.Error("Failed to get the number of worker nodes")
		os.Exit(1)
	} else {
		meta.NoWP = v
	}

	if v, ok := k.getCounterValue("Enter the number of Etcd Nodes", func(v int) bool {
		return v >= 3
	}, 3); !ok {
		k.l.Error("Failed to get the number of etcd nodes")
		os.Exit(1)
	} else {
		meta.NoDS = v
	}

	bootstrapVers, err := metaClient.ListAllBootstrapVersions()
	if err != nil {
		k.l.Error("Failed to get the list of bootstrap versions", "Reason", err)
		os.Exit(1)
	}
	if v, err := k.menuDriven.DropDownList("Select the bootstrap version", bootstrapVers, cli.WithDefaultValue(bootstrapVers[0])); err != nil {
		k.l.Error("Failed to get the bootstrap version", "Reason", err)
		os.Exit(1)
	} else {
		k.l.Debug(k.Ctx, "Selected bootstrap version", "Version", v)
		meta.K8sVersion = v
	}

	etcdVers, err := metaClient.ListAllEtcdVersions()
	if err != nil {
		k.l.Error("Failed to get the list of etcd versions", "Reason", err)
		os.Exit(1)
	}
	if v, err := k.menuDriven.DropDownList("Select the etcd version", etcdVers, cli.WithDefaultValue(etcdVers[0])); err != nil {
		k.l.Error("Failed to get the etcd version", "Reason", err)
		os.Exit(1)
	} else {
		k.l.Debug(k.Ctx, "Selected etcd version", "Version", v)
		meta.EtcdVersion = v
	}

	_, err = metaClient.PriceCalculator(
		controllerMeta.PriceCalculatorInput{
			Currency:              cp.Price.Currency,
			NoOfWorkerNodes:       meta.NoWP,
			NoOfControlPlaneNodes: meta.NoCP,
			NoOfEtcdNodes:         meta.NoDS,
			ControlPlaneMachine:   cp,
			WorkerMachine:         wp,
			EtcdMachine:           etcd,
			LoadBalancerMachine:   lb,
		})
	if err != nil {
		k.l.Error("Failed to calculate the price", "Reason", err)
		os.Exit(1)
	}

	managedCNI, defaultCNI, ksctlCNI, defaultKsctl, err := metaClient.ListBootstrapCNIs()
	if err != nil {
		k.l.Error("Failed to get the list of self managed CNIs", "Reason", err)
		os.Exit(1)
	}

	v, err := k.handleCNI(managedCNI, defaultCNI, ksctlCNI, defaultKsctl)
	if err != nil {
		k.l.Error("Failed to get the CNI", "Reason", err)
		os.Exit(1)
	}

	meta.Addons = v

	k.metadataSummary(*meta)

	if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterCreate, telemetry.TelemetryMeta{
		CloudProvider:     meta.Provider,
		StorageDriver:     meta.StateLocation,
		Region:            meta.Region,
		ClusterType:       meta.ClusterType,
		BootstrapProvider: meta.K8sDistro,
		K8sVersion:        meta.K8sVersion,
		Addons:            telemetry.TranslateMetadata(meta.Addons),
	}); err != nil {
		k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
	}

	if ok, _ := k.menuDriven.Confirmation("Do you want to proceed with the cluster creation", cli.WithDefaultValue("no")); !ok {
		os.Exit(1)
	}

	c, err := controllerSelfManaged.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: *meta,
		},
	)
	if err != nil {
		k.l.Error("Failed to create the controller", "Reason", err)
		os.Exit(1)
	}

	if err := c.Create(); err != nil {
		k.l.Error("Failed to create the cluster", "Reason", err)
		os.Exit(1)
	}

	return
}

func (k *KsctlCommand) metadataForManagedCluster(
	meta *controller.Metadata,
) {
	metaClient, err := controllerMeta.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: *meta,
		},
	)
	if err != nil {
		k.l.Error("Failed to create the controller", "Reason", err)
		os.Exit(1)
	}

	if v, ok := k.getCounterValue("Enter the number of Managed Nodes", func(v int) bool {
		return v > 0
	}, 1); !ok {
		k.l.Error("Failed to get the number of managed nodes")
		os.Exit(1)
	} else {
		meta.NoMP = v
	}

	if meta.Provider != consts.CloudLocal {
		k.handleRegionSelection(metaClient, meta)
		vm := k.handleInstanceTypeSelection(metaClient, meta, "Select instance_type for Managed Nodes")
		meta.ManagedNodeType = vm.Sku

		k.menuDriven.GetProgressAnimation().Start("Fetching the managed cluster offerings")

		listOfOfferings, err := metaClient.ListAllManagedClusterManagementOfferings(meta.Region, nil)
		if err != nil {
			k.menuDriven.GetProgressAnimation().Stop()
			k.l.Error("Failed to sync the metadata", "Reason", err)
			os.Exit(1)
		}
		k.menuDriven.GetProgressAnimation().Stop()

		offeringSelected := ""

		if v, ok := k.getSelectedManagedClusterOffering("Select the managed cluster offering", listOfOfferings); !ok {
			k.l.Error("Failed to get the managed cluster offering")
			os.Exit(1)
		} else {
			offeringSelected = v
		}

		_, err = metaClient.PriceCalculator(
			controllerMeta.PriceCalculatorInput{
				ManagedControlPlaneMachine: listOfOfferings[offeringSelected],
				NoOfWorkerNodes:            meta.NoMP,
				WorkerMachine:              vm,
			})
		if err != nil {
			k.l.Error("Failed to calculate the price", "Reason", err)
			os.Exit(1)
		}
	}

	managedCNI, defaultCNI, ksctlCNI, defaultKsctl, err := metaClient.ListManagedCNIs()
	if err != nil {
		k.l.Error("Failed to get the list of managed CNIs", "Reason", err)
		os.Exit(1)
	}

	if v, err := k.handleCNI(managedCNI, defaultCNI, ksctlCNI, defaultKsctl); err != nil {
		k.l.Error("Failed to get the CNI", "Reason", err)
		os.Exit(1)
	} else {
		meta.Addons = v
	}

	k.handleManagedK8sVersion(metaClient, meta)

	k.metadataSummary(*meta)

	if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterCreate, telemetry.TelemetryMeta{
		CloudProvider:     meta.Provider,
		StorageDriver:     meta.StateLocation,
		Region:            meta.Region,
		ClusterType:       meta.ClusterType,
		BootstrapProvider: meta.K8sDistro,
		K8sVersion:        meta.K8sVersion,
		Addons:            telemetry.TranslateMetadata(meta.Addons),
	}); err != nil {
		k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
	}

	if ok, _ := k.menuDriven.Confirmation("Do you want to proceed with the cluster creation", cli.WithDefaultValue("no")); !ok {
		os.Exit(1)
	}

	c, err := controllerManaged.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: *meta,
		},
	)
	if err != nil {
		k.l.Error("Failed to create the controller", "Reason", err)
		os.Exit(1)
	}

	if err := c.Create(); err != nil {
		k.l.Error("Failed to create the cluster", "Reason", err)
		os.Exit(1)
	}

	return
}
