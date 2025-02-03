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
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"os"

	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
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
			// we need to collect the cloud provider, and clusterType for the metadata to work

			meta := controller.Metadata{}

			if v, ok := k.getClusterName(); !ok {
				os.Exit(1)
			} else {
				meta.ClusterName = v
			}

			if v, ok := k.getSelectedClusterType(); !ok {
				os.Exit(1)
			} else {
				meta.ClusterType = v
			}

			if v, ok := k.getSelectedCloudProvider(); !ok {
				os.Exit(1)
			} else {
				meta.Provider = v
			}

			if v, ok := k.getSelectedStorageDriver(); !ok {
				os.Exit(1)
			} else {
				k.l.Debug(k.Ctx, "DropDown input", "storageDriver", v)
				meta.StateLocation = consts.KsctlStore(v)
			}

			managerClient, err := controllerMeta.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: meta,
				},
			)
			if err != nil {
				k.l.Error("Failed to create the controller", "Reason", err)
				os.Exit(1)
			}

			ss := cli.GetSpinner()
			ss.Start("Fetching the region list")

			listOfRegions, err := managerClient.ListAllRegions()
			if err != nil {
				ss.StopWithFailure("Failed to fetch the region list", "Reason", err)
				k.l.Error("Failed to sync the metadata", "Reason", err)
				os.Exit(1)
			}
			ss.Stop()

			if v, ok := k.getSelectedRegion(listOfRegions); !ok {
				os.Exit(1)
			} else {
				meta.Region = v
			}
			ss = cli.GetSpinner()
			ss.Start("Fetching the instance type list")

			listOfVMs, err := managerClient.ListAllInstances(meta.Region)
			if err != nil {
				ss.StopWithFailure("Failed to fetch the instance type list", "Reason", err)
				k.l.Error("Failed to sync the metadata", "Reason", err)
				os.Exit(1)
			}
			ss.Stop()

			if meta.ClusterType == consts.ClusterTypeMang {
				if !k.handleManagedCluster(managerClient, &meta, listOfVMs) {
					os.Exit(1)
				}
			}

			if ok, _ := cli.Confirmation("Do you want to proceed with the cluster creation", "no"); !ok {
				os.Exit(1)
			}

			k.l.Success(k.Ctx, "Created the cluster", "Name", meta.ClusterName)
		},
	}

	return cmd
}

func (k *KsctlCommand) handleManagedCluster(
	managerClient *controllerMeta.Controller,
	meta *controller.Metadata,
	listOfVMs map[string]provider.InstanceRegionOutput,
) bool {
	if v, ok := k.getSelectedInstanceType("Select instance_type for Managed Nodes", listOfVMs); !ok {
		return false
	} else {
		meta.ManagedNodeType = v
	}

	if v, ok := k.getCounterValue("Enter the number of Managed Nodes", func(v int) bool {
		return v > 0
	}); !ok {
		return false
	} else {
		meta.NoMP = v
	}

	ss := cli.GetSpinner()
	ss.Start("Fetching the managed cluster offerings")

	listOfOfferings, err := managerClient.ListAllManagedClusterManagementOfferings(meta.Region)
	if err != nil {
		ss.StopWithFailure("Failed to fetch the managed cluster offerings", "Reason", err)
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	var offeringSelected provider.ManagedClusterOutput
	for _, v := range listOfOfferings {
		if v.Tier == "Standard" {
			offeringSelected = v
			break
		}
	}
	k.l.Print(k.Ctx, "Managed Cluster Offering", "Name", offeringSelected.Sku, "Cost", offeringSelected.GetCost())

	//if v, ok := k.getSelectedManagedClusterOffering("Select the managed cluster offering", listOfOfferings); !ok {
	//	return false
	//} else {
	//	meta.ManagedNodeType = v
	//}
	priceCalculator, err := managerClient.PriceCalculator(
		controllerMeta.PriceCalculatorInput{
			ManagedControlPlaneMachine: offeringSelected,
			NoOfWorkerNodes:            meta.NoMP,
			WorkerMachine:              listOfVMs[meta.ManagedNodeType],
		})
	if err != nil {
		k.l.Error("Failed to calculate the price", "Reason", err)
		return false
	}

	priceOfVM := listOfVMs[meta.ManagedNodeType].GetCost() * float64(meta.NoMP)
	curr := listOfVMs[meta.ManagedNodeType].Price.Currency

	k.l.Box(k.Ctx, "Cost Summary", fmt.Sprintf(`
Managed Node(s) Cost = %.2f X %d = %.2f %s
Management Offering = %.2f %s
Total Cost = %.2f %s
`,
		listOfVMs[meta.ManagedNodeType].GetCost(), meta.NoMP, priceOfVM, curr,
		offeringSelected.GetCost(), curr,
		priceCalculator, curr,
	),
	)

	return true
}
