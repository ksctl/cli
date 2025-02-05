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

	"github.com/ksctl/cli/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	controllerMeta "github.com/ksctl/ksctl/v2/pkg/handler/cluster/metadata"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

func (k *KsctlCommand) BaseMetadataFields(m *controller.Metadata) {
	if v, ok := k.getClusterName(); !ok {
		os.Exit(1)
	} else {
		m.ClusterName = v
	}

	if v, ok := k.getSelectedClusterType(); !ok {
		os.Exit(1)
	} else {
		m.ClusterType = v
	}

	if v, ok := k.getSelectedCloudProvider(); !ok {
		os.Exit(1)
	} else {
		m.Provider = v
	}

	if v, ok := k.getSelectedStorageDriver(); !ok {
		os.Exit(1)
	} else {
		m.StateLocation = consts.KsctlStore(v)
	}
}

func (k *KsctlCommand) handleRegionSelection(meta *controllerMeta.Controller, m *controller.Metadata) {
	ss := cli.GetSpinner()
	ss.Start("Fetching the region list")

	listOfRegions, err := meta.ListAllRegions()
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	if v, ok := k.getSelectedRegion(listOfRegions); !ok {
		os.Exit(1)
	} else {
		m.Region = v
	}
}

func (k *KsctlCommand) handleInstanceTypeSelection(meta *controllerMeta.Controller, m *controller.Metadata) provider.InstanceRegionOutput {
	ss := cli.GetSpinner()
	ss.Start("Fetching the instance type list")

	listOfVMs, err := meta.ListAllInstances(m.Region)
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()

	v, ok := k.getSelectedInstanceType("Select instance_type for Managed Nodes", listOfVMs)
	if !ok {
		k.l.Error("Failed to get the instance type")
		os.Exit(1)
	}
	return listOfVMs[v]
}

func (k *KsctlCommand) handleManagedK8sVersion(meta *controllerMeta.Controller, m *controller.Metadata) {
	ss := cli.GetSpinner()
	ss.Start("Fetching the managed cluster k8s versions")

	listOfK8sVersions, err := meta.ListAllManagedClusterK8sVersions(m.Region)
	if err != nil {
		ss.Stop()
		k.l.Error("Failed to sync the metadata", "Reason", err)
		os.Exit(1)
	}
	ss.Stop()
	k.l.Debug(k.Ctx, "List of k8s versions", "versions", listOfK8sVersions)

	if v, ok := k.getSelectedK8sVersion("Select the k8s version for Managed Cluster", listOfK8sVersions); !ok {
		k.l.Error("Failed to get the k8s version")
		os.Exit(1)
	} else {
		m.K8sVersion = v
	}
}
