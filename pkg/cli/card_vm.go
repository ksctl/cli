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

package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"strings"
)

type cardVM struct {
	rr provider.InstanceRegionOutput
}

func (c cardVM) GetUpper() string {
	resp := strings.Builder{}

	resp.WriteString(fmt.Sprintf(
		"Price: %s\n",
		color.HiMagentaString("%s %.2f", c.rr.Price.Currency, c.rr.GetCost()),
	))
	if c.rr.EmboddedEmissions != nil {
		resp.WriteString(fmt.Sprintf("🏭 Embodied: %.2f %s\n", c.rr.EmboddedEmissions.EmboddedCo2, c.rr.EmboddedEmissions.Co2Unit))
	} else {
		resp.WriteString(color.HiYellowString("Emissions data is currently unavailable 🌍\n"))
	}

	return resp.String()
}

func (c cardVM) GetLower() string {
	resp := strings.Builder{}

	resp.WriteString(fmt.Sprintf("Code: %s\n", color.HiMagentaString(c.rr.Sku)))
	resp.WriteString(fmt.Sprintf("vCPUs: %d\n", c.rr.VCpus))
	resp.WriteString(fmt.Sprintf("Memory: %d GB\n", c.rr.Memory))
	arch := string(c.rr.CpuArch)
	if arch == string(provider.ArchArm64) {
		arch = color.HiCyanString(arch)
	} else {
		arch = color.HiBlueString(arch)
	}
	resp.WriteString(fmt.Sprintf("Arch: %s\n", arch))

	diskSku := c.rr.Disk.Sku
	diskSize := c.rr.Disk.Size
	diskTier := c.rr.Disk.Tier

	if diskSku != nil {
		resp.WriteString(fmt.Sprintf("Disk: %s\n", *diskSku))
	}

	if diskTier != nil {
		resp.WriteString(fmt.Sprintf("DiskTier: %s\n", *diskTier))
	}

	if diskSize != nil {
		resp.WriteString(fmt.Sprintf("DiskSize: %d GB", *diskSize))
	}

	return resp.String()
}

type cardVMs struct {
	rr         provider.InstancesRegionOutput
	tt         []cardVM
	lenOfItems int
}

func (c cardVMs) LenOfItems() int {
	return c.lenOfItems
}

func (c cardVMs) GetItem(i int) CardItem {
	return c.tt[i]
}

func (c cardVMs) GetInstruction() string {
	return "← → to navigate • enter to select instanceType • q to quit"
}

func (c cardVMs) GetResult(i int) string {
	return c.rr[i].Sku
}

func (c cardVMs) GetCardConfiguration() (cardWidth, noOfVisibleItems int) {
	return 33, 3
}

func ConverterForInstanceTypesForCards(vms provider.InstancesRegionOutput) CardPack {
	res := new(cardVMs)
	res.lenOfItems = len(vms)
	res.rr = vms

	for i, _ := range vms {
		res.tt = append(res.tt, cardVM{vms[i]})
	}

	return res
}
