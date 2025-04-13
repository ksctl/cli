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
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/provider"
)

type cardRegion struct {
	b provider.RegionOutput
}

func (c cardRegion) GetUpper() string {
	return "TODO"
}

func (c cardRegion) GetLower() string {
	return "Name: " + color.HiMagentaString(c.b.Name) + "\nCode: " + color.HiCyanString(c.b.Sku)
}

type cardRegions struct {
	rr       []provider.RegionOutput
	tt       []cardRegion
	lenItems int
}

func (c cardRegions) LenOfItems() int {
	return c.lenItems
}

func (c cardRegions) GetItem(i int) CardItem {
	return c.tt[i]
}

func (c cardRegions) GetInstruction() string {
	return "← → to navigate • enter to select region • q to quit"
}

func (c cardRegions) GetResult(i int) string {
	return c.rr[i].Sku
}

func (c cardRegions) GetCardConfiguration() (cardWidth, noOfVisibleItems int) {
	return 45, 2
}

func ConverterForRegionOutputForCards(regions provider.RegionsOutput) CardPack {
	res := new(cardRegions)
	res.lenItems = len(regions)

	for i, _ := range regions {
		res.tt = append(res.tt, cardRegion{regions[i]})
	}

	return res
}
