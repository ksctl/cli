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
	"cmp"
	"fmt"
	"slices"

	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

func (k *KsctlCommand) findManagedOfferingCostAcrossRegions(
	meta controller.Metadata,
	availRegions []provider.RegionOutput,
	managedOfferingSku string,
) (map[string]float64, error) {
	metaClient, err := controllerMeta.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: meta,
		},
	)
	if err != nil {
		return nil, err
	}
	resultChan := make(chan struct {
		region string
		price  float64
		err    error
	}, len(availRegions))

	for _, region := range availRegions {
		regSku := region.Sku
		go func(sku string) {
			_price, err := metaClient.ListAllManagedClusterManagementOfferings(regSku, nil)
			if err == nil {
				v, ok := _price[managedOfferingSku]
				if ok {
					resultChan <- struct {
						region string
						price  float64
						err    error
					}{sku, v.GetCost(), nil}
				} else {
					resultChan <- struct {
						region string
						price  float64
						err    error
					}{sku, 0.0, fmt.Errorf("managed offering not found")}
				}

			} else {
				resultChan <- struct {
					region string
					price  float64
					err    error
				}{sku, 0.0, err}
			}
		}(regSku)
	}

	cost := make(map[string]float64, len(availRegions))
	for i := 0; i < len(availRegions); i++ {
		result := <-resultChan
		if result.err == nil {
			cost[result.region] = result.price
		}
	}

	return cost, nil
}

// findInstanceCostAcrossRegions it returns a map of K[V] where K is the region and V is the cost of the instance
func (k *KsctlCommand) findInstanceCostAcrossRegions(
	meta controller.Metadata,
	availRegions []provider.RegionOutput,
	instanceSku string,
) (map[string]float64, error) {
	metaClient, err := controllerMeta.NewController(
		k.Ctx,
		k.l,
		&controller.Client{
			Metadata: meta,
		},
	)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan struct {
		region string
		price  float64
		err    error
	}, len(availRegions))

	for _, region := range availRegions {
		regSku := region.Sku
		go func(sku string) {
			price, err := metaClient.GetPriceForInstance(sku, instanceSku)
			resultChan <- struct {
				region string
				price  float64
				err    error
			}{sku, price, err}
		}(regSku)
	}

	cost := make(map[string]float64, len(availRegions))
	for i := 0; i < len(availRegions); i++ {
		result := <-resultChan
		if result.err == nil {
			cost[result.region] = result.price
		}
	}

	return cost, nil
}

type RecommendationManagedCost struct {
	region    string
	totalCost float64

	cpCost float64
	wpCost float64
}

func (k *KsctlCommand) getBestRegionsWithTotalCostManaged(
	allAvailRegions []provider.RegionOutput,
	costForCP map[string]float64,
	costForWP map[string]float64,
) []RecommendationManagedCost {

	checkRegion := func(region string, m map[string]float64) bool {
		_, ok := m[region]
		return ok
	}

	var costForCluster []RecommendationManagedCost

	for _, region := range allAvailRegions {
		if !checkRegion(region.Sku, costForCP) ||
			!checkRegion(region.Sku, costForWP) {
			continue
		}

		totalCost := costForCP[region.Sku] + costForWP[region.Sku]

		costForCluster = append(costForCluster, RecommendationManagedCost{
			region:    region.Sku,
			cpCost:    costForCP[region.Sku],
			wpCost:    costForWP[region.Sku],
			totalCost: totalCost,
		})
	}

	slices.SortFunc(costForCluster, func(a, b RecommendationManagedCost) int {
		return cmp.Compare(a.totalCost, b.totalCost)
	})

	return costForCluster
}

type RecommendationSelfManagedCost struct {
	region    string
	totalCost float64

	cpCost   float64
	wpCost   float64
	etcdCost float64
	lbCost   float64
}

func (k *KsctlCommand) getBestRegionsWithTotalCostSelfManaged(
	allAvailRegions []provider.RegionOutput,
	costForCP map[string]float64,
	costForWP map[string]float64,
	costForDS map[string]float64,
	costForLB map[string]float64,
) []RecommendationSelfManagedCost {

	checkRegion := func(region string, m map[string]float64) bool {
		_, ok := m[region]
		return ok
	}

	var costForCluster []RecommendationSelfManagedCost

	for _, region := range allAvailRegions {
		if !checkRegion(region.Sku, costForCP) ||
			!checkRegion(region.Sku, costForWP) ||
			!checkRegion(region.Sku, costForDS) ||
			!checkRegion(region.Sku, costForLB) {
			continue
		}

		totalCost := costForCP[region.Sku] + costForWP[region.Sku] + costForDS[region.Sku] + costForLB[region.Sku]

		costForCluster = append(costForCluster, RecommendationSelfManagedCost{
			region:    region.Sku,
			cpCost:    costForCP[region.Sku],
			wpCost:    costForWP[region.Sku],
			etcdCost:  costForDS[region.Sku],
			lbCost:    costForLB[region.Sku],
			totalCost: totalCost,
		})
	}

	slices.SortFunc(costForCluster, func(a, b RecommendationSelfManagedCost) int {
		return cmp.Compare(a.totalCost, b.totalCost)
	})

	return costForCluster
}

// OptimizeSelfManagedInstanceTypesAcrossRegions it returns a sorted list of regions based on the total cost in ascending order across all the regions
//
//	It is a core function that is used to optimize the cost of the self-managed cluster instanceType across all the regions (Cost Optimization)
func (k *KsctlCommand) OptimizeSelfManagedInstanceTypesAcrossRegions(
	meta *controller.Metadata,
	allAvailRegions []provider.RegionOutput,
	cp provider.InstanceRegionOutput,
	wp provider.InstanceRegionOutput,
	etcd provider.InstanceRegionOutput,
	lb provider.InstanceRegionOutput,
) []RecommendationSelfManagedCost {
	cpInstanceCosts, err := k.findInstanceCostAcrossRegions(*meta, allAvailRegions, cp.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of control plane instances", "Reason", err)
	}

	wpInstanceCosts, err := k.findInstanceCostAcrossRegions(*meta, allAvailRegions, wp.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of worker plane instances", "Reason", err)
	}

	etcdInstanceCosts, err := k.findInstanceCostAcrossRegions(*meta, allAvailRegions, etcd.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of etcd instances", "Reason", err)
	}

	lbInstanceCosts, err := k.findInstanceCostAcrossRegions(*meta, allAvailRegions, lb.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of load balancer instances", "Reason", err)
	}

	return k.getBestRegionsWithTotalCostSelfManaged(
		allAvailRegions,
		cpInstanceCosts,
		wpInstanceCosts,
		etcdInstanceCosts,
		lbInstanceCosts,
	)
}

func (k *KsctlCommand) OptimizeManagedOfferingsAcrossRegions(
	meta *controller.Metadata,
	allAvailRegions []provider.RegionOutput,
	cp provider.ManagedClusterOutput,
	wp provider.InstanceRegionOutput,
) []RecommendationManagedCost {
	wpInstanceCosts, err := k.findInstanceCostAcrossRegions(*meta, allAvailRegions, wp.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of worker plane instances", "Reason", err)
	}

	cpInstanceCosts, err := k.findManagedOfferingCostAcrossRegions(*meta, allAvailRegions, cp.Sku)
	if err != nil {
		k.l.Error("Failed to get the cost of control plane managed offerings", "Reason", err)
	}

	return k.getBestRegionsWithTotalCostManaged(
		allAvailRegions,
		cpInstanceCosts,
		wpInstanceCosts,
	)
}

func (k *KsctlCommand) PrintRecommendationSelfManagedCost(
	costs []RecommendationSelfManagedCost,
	noOfCP int,
	noOfWP int,
	noOfDS int,
	instanceTypeCP string,
	instanceTypeWP string,
	instanceTypeDS string,
	instanceTypeLB string,
) {
	k.l.Print(k.Ctx,
		"Here is your recommendation",
		"Parameter", "Region wise cost",
		"OptimizedRegion", color.HiCyanString(costs[0].region),
	)

	headers := []string{
		"Region",
		fmt.Sprintf("Control Plane (%s)", instanceTypeCP),
		fmt.Sprintf("Worker Plane (%s)", instanceTypeWP),
		fmt.Sprintf("Etcd Nodes (%s)", instanceTypeDS),
		fmt.Sprintf("Load Balancer (%s)", instanceTypeLB),
		"Total Monthly Cost",
	}

	var data [][]string
	for _, cost := range costs {
		total := cost.cpCost*float64(noOfCP) + cost.wpCost*float64(noOfWP) + cost.etcdCost*float64(noOfDS) + cost.lbCost
		data = append(data, []string{
			cost.region,
			fmt.Sprintf("$%.2f X %d", cost.cpCost, noOfCP),
			fmt.Sprintf("$%.2f X %d", cost.wpCost, noOfWP),
			fmt.Sprintf("$%.2f X %d", cost.etcdCost, noOfDS),
			fmt.Sprintf("$%.2f X 1", cost.lbCost),
			fmt.Sprintf("$%.2f", total),
		})
	}

	k.l.Table(k.Ctx, headers, data)
}

func (k *KsctlCommand) PrintRecommendationManagedCost(
	costs []RecommendationManagedCost,
	noOfWP int,
	managedOfferingCP string,
	instanceTypeWP string,
) {
	k.l.Print(k.Ctx,
		"Here is your recommendation",
		"Parameter", "Region wise cost",
		"OptimizedRegion", color.HiCyanString(costs[0].region),
	)

	headers := []string{
		"Region",
		fmt.Sprintf("Control Plane (%s)", managedOfferingCP),
		fmt.Sprintf("Worker Plane (%s)", instanceTypeWP),
		"Total Monthly Cost",
	}

	var data [][]string
	for _, cost := range costs {
		total := cost.cpCost + cost.wpCost*float64(noOfWP)
		data = append(data, []string{
			cost.region,
			fmt.Sprintf("$%.2f X 1", cost.cpCost),
			fmt.Sprintf("$%.2f X %d", cost.wpCost, noOfWP),
			fmt.Sprintf("$%.2f", total),
		})
	}

	k.l.Table(k.Ctx, headers, data)
}
