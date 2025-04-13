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

type cardRegion struct {
	b provider.RegionOutput
}

func (c cardRegion) GetUpper() string {
	priceStr := strings.Builder{}

	if c.b.Emission != nil {
		dco2Val := fmt.Sprintf("%.2f %s", c.b.Emission.DirectCarbonIntensity, c.b.Emission.Unit)
		if c.b.Emission.DirectCarbonIntensity > DirectCo2MediumThreshold {
			dco2Val = color.HiRedString(dco2Val)
		} else if c.b.Emission.DirectCarbonIntensity > DirectCo2LowThreshold {
			dco2Val = color.HiYellowString(dco2Val)
		} else {
			dco2Val = color.HiGreenString(dco2Val)
		}

		renewableVal := fmt.Sprintf("%.1f%%", c.b.Emission.RenewablePercentage)
		if c.b.Emission.RenewablePercentage < RenewableLowThreshold {
			renewableVal = color.HiRedString(renewableVal)
		} else if c.b.Emission.RenewablePercentage < RenewableMediumThreshold {
			renewableVal = color.HiYellowString(renewableVal)
		} else {
			renewableVal = color.HiGreenString(renewableVal)
		}

		lowco2Val := fmt.Sprintf("%.1f%%", c.b.Emission.LowCarbonPercentage)
		if c.b.Emission.LowCarbonPercentage < LowCarbonLowThreshold {
			lowco2Val = color.HiRedString(lowco2Val)
		} else if c.b.Emission.LowCarbonPercentage < LowCarbonMediumThreshold {
			lowco2Val = color.HiYellowString(lowco2Val)
		} else {
			lowco2Val = color.HiGreenString(lowco2Val)
		}

		lcaintensityVal := fmt.Sprintf("%.1f %s", c.b.Emission.LCACarbonIntensity, c.b.Emission.Unit)
		if c.b.Emission.LCACarbonIntensity > LCACarbonIntensityHighThreshold {
			lcaintensityVal = color.HiRedString(lcaintensityVal)
		} else if c.b.Emission.LCACarbonIntensity > LCACarbonIntensityMediumThreshold {
			lcaintensityVal = color.HiYellowString(lcaintensityVal)
		} else {
			lcaintensityVal = color.HiGreenString(lcaintensityVal)
		}
		priceStr.WriteString(fmt.Sprintf("🌍 Direct Co2: %s\n", dco2Val))
		priceStr.WriteString(fmt.Sprintf("🌱 Renewable: %s\n", renewableVal))
		priceStr.WriteString(fmt.Sprintf("💨 Low Carbon: %s\n", lowco2Val))
		priceStr.WriteString(fmt.Sprintf("🔄 Lifecycle Co2: %s\n", lcaintensityVal))
	} else {
		priceStr.WriteString(color.HiYellowString("Emissions data is currently unavailable 🌍\n"))
	}

	return priceStr.String()
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
	res.rr = regions

	for i, _ := range regions {
		res.tt = append(res.tt, cardRegion{regions[i]})
	}

	return res
}
