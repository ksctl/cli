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
	"strings"

	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/provider/optimizer"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	recommendI   *optimizer.RecommendationAcrossRegions
	clusterType  consts.KsctlClusterType
	currentPlan  int
	windowWidth  int
	windowHeight int
	help         help.Model
	keys         keyMap
	quitting     bool
	selectedPlan int // -1 means no selection yet
}

type keyMap struct {
	left     key.Binding
	right    key.Binding
	selected key.Binding
	quit     key.Binding
	help     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.left, k.right, k.selected},
		{k.help, k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "previous plan"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next plan"),
		),
		selected: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "select plan"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}

func newModel(
	clusterType consts.KsctlClusterType,
	recommendI *optimizer.RecommendationAcrossRegions) Model {
	return Model{
		recommendI:   recommendI,
		clusterType:  clusterType,
		currentPlan:  0,
		help:         help.New(),
		keys:         newKeyMap(),
		selectedPlan: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keys.left):
			// Move to previous plan
			if m.currentPlan > 0 {
				m.currentPlan--
			}
			return m, nil

		case key.Matches(msg, m.keys.right):
			// Move to next plan
			if m.currentPlan < len(m.recommendI.RegionRecommendations)-1 {
				m.currentPlan++
			}
			return m, nil

		case key.Matches(msg, m.keys.selected):
			// Select current plan
			m.selectedPlan = m.currentPlan
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	}

	return m, nil
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

func (m Model) View() string {
	if m.quitting {
		if m.selectedPlan >= 0 && m.selectedPlan < len(m.recommendI.RegionRecommendations) {
			plan := m.recommendI.RegionRecommendations[m.selectedPlan]
			return plan.Region.Sku
		}
		return "NULL"
	}

	// Calculate dimensions based on window size
	maxWidth := m.windowWidth
	if maxWidth == 0 {
		maxWidth = 100 // Default width
	}

	// Card styling
	cardWidth := 45   // Increased from 30 to 45 for more information display
	visibleCards := 2 // Show exactly 2 visible cards

	// Add some top padding
	var builder strings.Builder
	builder.WriteString("\n\n")

	// Determine which cards to show
	var startIdx int
	if m.currentPlan == 0 {
		startIdx = 0
	} else if m.currentPlan == len(m.recommendI.RegionRecommendations)-1 {
		startIdx = max(0, len(m.recommendI.RegionRecommendations)-visibleCards)
	} else {
		startIdx = max(0, m.currentPlan-1)
	}
	endIdx := min(len(m.recommendI.RegionRecommendations), startIdx+visibleCards)

	cards := make([]string, endIdx-startIdx)

	for i := startIdx; i < endIdx; i++ {
		plan := m.recommendI.RegionRecommendations[i]
		isActive := i == m.currentPlan

		var borderColor, textColor, separatorColor lipgloss.Color
		var borderStyle lipgloss.Border

		if isActive {
			borderColor = lipgloss.Color("10")
			textColor = lipgloss.Color("#EEEEC7")
			separatorColor = lipgloss.Color("#555555")
			borderStyle = lipgloss.RoundedBorder()
		} else {
			borderColor = lipgloss.Color("#555555")
			textColor = lipgloss.Color("#AAAAAA")
			separatorColor = lipgloss.Color("#555555")
			borderStyle = lipgloss.RoundedBorder()
		}

		priceStyle := lipgloss.NewStyle().
			Foreground(textColor).
			Width(cardWidth - 3).
			Align(lipgloss.Center).
			PaddingTop(1)

		if isActive {
			priceStyle = priceStyle.Bold(true)
		}

		priceStr := strings.Builder{}
		priceDrop := (m.recommendI.CurrentTotalCost - plan.TotalCost) / m.recommendI.CurrentTotalCost * 100
		priceStr.WriteString(fmt.Sprintf("Price: %s %s\n", color.MagentaString(fmt.Sprintf("$%.2f", plan.TotalCost)), color.HiGreenString(fmt.Sprintf("↓ %.0f%%", priceDrop))))
		if plan.Region.Emission != nil {
			dco2_val := fmt.Sprintf("%.2f %s", plan.Region.Emission.DirectCarbonIntensity, plan.Region.Emission.Unit)
			if plan.Region.Emission.DirectCarbonIntensity > DirectCo2MediumThreshold {
				dco2_val = color.HiRedString(dco2_val)
			} else if plan.Region.Emission.DirectCarbonIntensity > DirectCo2LowThreshold {
				dco2_val = color.HiYellowString(dco2_val)
			} else {
				dco2_val = color.HiGreenString(dco2_val)
			}

			reneable_val := fmt.Sprintf("%.1f%%", plan.Region.Emission.RenewablePercentage)
			if plan.Region.Emission.RenewablePercentage < RenewableLowThreshold {
				reneable_val = color.HiRedString(reneable_val)
			} else if plan.Region.Emission.RenewablePercentage < RenewableMediumThreshold {
				reneable_val = color.HiYellowString(reneable_val)
			} else {
				reneable_val = color.HiGreenString(reneable_val)
			}

			lowCo2_val := fmt.Sprintf("%.1f%%", plan.Region.Emission.LowCarbonPercentage)
			if plan.Region.Emission.LowCarbonPercentage < LowCarbonLowThreshold {
				lowCo2_val = color.HiRedString(lowCo2_val)
			} else if plan.Region.Emission.LowCarbonPercentage < LowCarbonMediumThreshold {
				lowCo2_val = color.HiYellowString(lowCo2_val)
			} else {
				lowCo2_val = color.HiGreenString(lowCo2_val)
			}

			lcaIntensity_val := fmt.Sprintf("%.1f %s", plan.Region.Emission.LCACarbonIntensity, plan.Region.Emission.Unit)
			if plan.Region.Emission.LCACarbonIntensity > LCACarbonIntensityHighThreshold {
				lcaIntensity_val = color.HiRedString(lcaIntensity_val)
			} else if plan.Region.Emission.LCACarbonIntensity > LCACarbonIntensityMediumThreshold {
				lcaIntensity_val = color.HiYellowString(lcaIntensity_val)
			} else {
				lcaIntensity_val = color.HiGreenString(lcaIntensity_val)
			}

			priceStr.WriteString(fmt.Sprintf("🌍 Direct Emissions: %s\n", dco2_val))
			priceStr.WriteString(fmt.Sprintf("🌱 Renewable Energy: %s\n", reneable_val))
			priceStr.WriteString(fmt.Sprintf("💨 Low Carbon Energy: %s\n", lowCo2_val))
			priceStr.WriteString(fmt.Sprintf("🔄 Lifecycle Emissions: %s\n", lcaIntensity_val))
		} else {
			priceStr.WriteString(color.HiYellowString("Emissions data is currently unavailable 🌍\n"))
		}

		priceSection := priceStyle.Render(priceStr.String())

		separator := lipgloss.NewStyle().
			Foreground(separatorColor).
			Width(cardWidth - 3).
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			Render(strings.Repeat("─", cardWidth-6))

		specsStyle := lipgloss.NewStyle().
			Foreground(textColor).
			Width(cardWidth - 3).
			Align(lipgloss.Center).
			PaddingLeft(1)

		if isActive {
			specsStyle = specsStyle.Bold(true)
		}

		specsStr := strings.Builder{}
		specsStr.WriteString(fmt.Sprintf("Region: %s\n", color.HiCyanString(plan.Region.Name)))
		if m.clusterType == consts.ClusterTypeSelfMang {
			specsStr.WriteString(fmt.Sprintf("ControlPlane: %s x %d\n", m.recommendI.InstanceTypeCP, m.recommendI.ControlPlaneCount))
			specsStr.WriteString(fmt.Sprintf("Worker: %s x %d\n", m.recommendI.InstanceTypeWP, m.recommendI.WorkerPlaneCount))
			specsStr.WriteString(fmt.Sprintf("Etcd: %s x %d\n", m.recommendI.InstanceTypeDS, m.recommendI.DataStoreCount))
			specsStr.WriteString(fmt.Sprintf("LoadBalancer: %s\n", m.recommendI.InstanceTypeLB))
		} else {
			specsStr.WriteString(fmt.Sprintf("ManagedOffering: %s\n", m.recommendI.ManagedOffering))
			specsStr.WriteString(fmt.Sprintf("Worker: %s x %d\n", m.recommendI.InstanceTypeWP, m.recommendI.WorkerPlaneCount))
		}

		specsSection := specsStyle.Render(specsStr.String())

		cardStyle := lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(borderColor).
			Width(cardWidth).
			Padding(0, 1)

		cardContent := lipgloss.JoinVertical(lipgloss.Left,
			priceSection,
			separator,
			specsSection,
		)

		cards[i-startIdx] = cardStyle.Render(cardContent)
	}

	var rowContent strings.Builder

	rowContent.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cards...))

	centeredStyle := lipgloss.NewStyle().Width(maxWidth).Align(lipgloss.Center)
	builder.WriteString(centeredStyle.Render(rowContent.String()))
	builder.WriteString("\n\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).
		Width(maxWidth)

	instructions := "← → to navigate • enter to select plan • q to skip changing region"
	instructions += " • Currently it costs " + fmt.Sprintf("`$%.2f`", m.recommendI.CurrentTotalCost) + " in " + color.HiCyanString(m.recommendI.CurrentRegion.Name) + "\n\n"
	builder.WriteString(instructionStyle.Render(instructions))

	return builder.String()
}

type RegionRecommendation struct {
	t *tea.Program
}

func NewRegionRecommendation(
	clusterType consts.KsctlClusterType,
	recommendI *optimizer.RecommendationAcrossRegions) *RegionRecommendation {
	model := newModel(clusterType, recommendI)
	t := tea.NewProgram(model)
	// t := tea.NewProgram(model, tea.WithAltScreen())
	return &RegionRecommendation{t: t}
}

func (t *RegionRecommendation) Run() (string, error) {

	finalModel, err := t.t.Run()
	if err != nil {
		return "", err
	}

	m, ok := finalModel.(Model)
	if !ok {
		return "", fmt.Errorf("failed to cast final model to Model type")
	}

	if m.selectedPlan >= 0 && m.selectedPlan < len(m.recommendI.RegionRecommendations) {
		plan := m.recommendI.RegionRecommendations[m.selectedPlan]
		return plan.Region.Sku, nil
	}

	if m.quitting {
		return "", nil
	}
	return "", fmt.Errorf("internal problem. invalid selected plan index: %d", m.selectedPlan)
}
