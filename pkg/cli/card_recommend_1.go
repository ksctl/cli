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
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/provider/optimizer"
	"strings"
)

type cardRecommendation struct {
	rr          optimizer.RegionRecommendation
	oo          *optimizer.RecommendationAcrossRegions
	clusterType consts.KsctlClusterType
}

const (
	RenewableLowThreshold             = 10.0
	RenewableMediumThreshold          = 40.0
	LowCarbonLowThreshold             = 10.0
	LowCarbonMediumThreshold          = 40.0
	DirectCo2LowThreshold             = 200.0
	DirectCo2MediumThreshold          = 400.0
	LCACarbonIntensityMediumThreshold = 100.0
	LCACarbonIntensityHighThreshold   = 200.0
)

func (c cardRecommendation) GetUpper() string {
	priceStr := strings.Builder{}
	priceDrop := (c.oo.CurrentTotalCost - c.rr.TotalCost) / c.oo.CurrentTotalCost * 100
	priceStr.WriteString(fmt.Sprintf(
		"Price: %s %s\n",
		color.MagentaString(fmt.Sprintf("$%.2f", c.rr.TotalCost)),
		color.HiGreenString(fmt.Sprintf("↓ %.0f%%", priceDrop)),
	))

	if c.rr.Region.Emission != nil {
		dco2Val := fmt.Sprintf("%.2f %s", c.rr.Region.Emission.DirectCarbonIntensity, c.rr.Region.Emission.Unit)
		if c.rr.Region.Emission.DirectCarbonIntensity > DirectCo2MediumThreshold {
			dco2Val = color.HiRedString(dco2Val)
		} else if c.rr.Region.Emission.DirectCarbonIntensity > DirectCo2LowThreshold {
			dco2Val = color.HiYellowString(dco2Val)
		} else {
			dco2Val = color.HiGreenString(dco2Val)
		}

		renewableVal := fmt.Sprintf("%.1f%%", c.rr.Region.Emission.RenewablePercentage)
		if c.rr.Region.Emission.RenewablePercentage < RenewableLowThreshold {
			renewableVal = color.HiRedString(renewableVal)
		} else if c.rr.Region.Emission.RenewablePercentage < RenewableMediumThreshold {
			renewableVal = color.HiYellowString(renewableVal)
		} else {
			renewableVal = color.HiGreenString(renewableVal)
		}

		lowco2Val := fmt.Sprintf("%.1f%%", c.rr.Region.Emission.LowCarbonPercentage)
		if c.rr.Region.Emission.LowCarbonPercentage < LowCarbonLowThreshold {
			lowco2Val = color.HiRedString(lowco2Val)
		} else if c.rr.Region.Emission.LowCarbonPercentage < LowCarbonMediumThreshold {
			lowco2Val = color.HiYellowString(lowco2Val)
		} else {
			lowco2Val = color.HiGreenString(lowco2Val)
		}

		lcaintensityVal := fmt.Sprintf("%.1f %s", c.rr.Region.Emission.LCACarbonIntensity, c.rr.Region.Emission.Unit)
		if c.rr.Region.Emission.LCACarbonIntensity > LCACarbonIntensityHighThreshold {
			lcaintensityVal = color.HiRedString(lcaintensityVal)
		} else if c.rr.Region.Emission.LCACarbonIntensity > LCACarbonIntensityMediumThreshold {
			lcaintensityVal = color.HiYellowString(lcaintensityVal)
		} else {
			lcaintensityVal = color.HiGreenString(lcaintensityVal)
		}

		e_r := c.rr.Region.Emission
		e_R := c.oo.CurrentRegion.Emission

		if e_R != nil {
			dco2Change := (e_r.DirectCarbonIntensity - e_R.DirectCarbonIntensity) / e_R.DirectCarbonIntensity * 100
			if dco2Change > 0 {
				dco2Val += color.HiRedString(fmt.Sprintf(" ↑ %.0f%%", dco2Change))
			} else {
				dco2Val += color.HiGreenString(fmt.Sprintf(" ↓ %.0f%%", -dco2Change))
			}

			renewableChange := (e_r.RenewablePercentage - e_R.RenewablePercentage) / e_R.RenewablePercentage * 100
			if renewableChange > 0 {
				renewableVal += color.HiGreenString(fmt.Sprintf(" ↑ %.0f%%", renewableChange))
			} else {
				renewableVal += color.HiRedString(fmt.Sprintf(" ↓ %.0f%%", -renewableChange))
			}

			lowco2Change := (e_r.LowCarbonPercentage - e_R.LowCarbonPercentage) / e_R.LowCarbonPercentage * 100
			if lowco2Change > 0 {
				lowco2Val += color.HiGreenString(fmt.Sprintf(" ↑ %.0f%%", lowco2Change))
			} else {
				lowco2Val += color.HiRedString(fmt.Sprintf(" ↓ %.0f%%", -lowco2Change))
			}

			lcaintensityChange := (e_r.LCACarbonIntensity - e_R.LCACarbonIntensity) / e_R.LCACarbonIntensity * 100
			if lcaintensityChange > 0 {
				lcaintensityVal += color.HiRedString(fmt.Sprintf(" ↑ %.0f%%", lcaintensityChange))
			} else {
				lcaintensityVal += color.HiGreenString(fmt.Sprintf(" ↓ %.0f%%", -lcaintensityChange))
			}
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

func (c cardRecommendation) GetLower() string {
	specsStr := strings.Builder{}
	specsStr.WriteString(fmt.Sprintf("Region: %s\n", color.HiCyanString(c.rr.Region.Name)))
	if c.clusterType == consts.ClusterTypeSelfMang {
		specsStr.WriteString(fmt.Sprintf("ControlPlane: %s x %d\n", c.oo.InstanceTypeCP, c.oo.ControlPlaneCount))
		specsStr.WriteString(fmt.Sprintf("Worker: %s x %d\n", c.oo.InstanceTypeWP, c.oo.WorkerPlaneCount))
		specsStr.WriteString(fmt.Sprintf("Etcd: %s x %d\n", c.oo.InstanceTypeDS, c.oo.DataStoreCount))
		specsStr.WriteString(fmt.Sprintf("LoadBalancer: %s\n", c.oo.InstanceTypeLB))
	} else {
		specsStr.WriteString(fmt.Sprintf("ManagedOffering: %s\n", c.oo.ManagedOffering))
		specsStr.WriteString(fmt.Sprintf("Worker: %s x %d\n", c.oo.InstanceTypeWP, c.oo.WorkerPlaneCount))
	}

	return specsStr.String()
}

type cardRecommendations struct {
	mm         consts.KsctlClusterType
	oo         *optimizer.RecommendationAcrossRegions
	tt         []cardRecommendation
	lenOfItems int
}

func (c cardRecommendations) LenOfItems() int {
	return c.lenOfItems
}

func (c cardRecommendations) GetItem(i int) CardItem {
	return c.tt[i]
}

func (c cardRecommendations) GetInstruction() string {
	instructions := "← → to navigate • enter to select plan • q to skip changing region"
	instructions += " • Currently it costs " + fmt.Sprintf("`$%.2f`", c.oo.CurrentTotalCost) + " in " + color.HiCyanString(c.oo.CurrentRegion.Name)

	return instructions
}

func (c cardRecommendations) GetResult(i int) string {
	return c.oo.RegionRecommendations[i].Region.Sku
}

func (c cardRecommendations) GetCardConfiguration() (cardWidth, noOfVisibleItems int) {
	return 45, 2
}

func ConverterForRecommendationIOutputForCards(oo *optimizer.RecommendationAcrossRegions, tt consts.KsctlClusterType) CardPack {
	res := new(cardRecommendations)
	res.oo = oo
	res.mm = tt
	res.lenOfItems = len(oo.RegionRecommendations)

	for i, _ := range oo.RegionRecommendations {
		res.tt = append(res.tt, cardRecommendation{
			oo.RegionRecommendations[i],
			oo,
			tt,
		})
	}

	return res
}
