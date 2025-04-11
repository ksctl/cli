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
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/ksctl/ksctl/v2/pkg/provider/optimizer"
	"os"
	"strings"

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

func newModel(recommendI *optimizer.RecommendationAcrossRegions) Model {
	return Model{
		recommendI:   recommendI,
		currentPlan:  0,
		help:         help.New(),
		keys:         newKeyMap(),
		selectedPlan: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
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

func (m Model) View() string {
	if m.quitting {
		if m.selectedPlan >= 0 && m.selectedPlan < len(m.recommendI.RegionRecommendations) {
			plan := m.recommendI.RegionRecommendations[m.selectedPlan]
			return plan.Region
		}
		return "NULL"
	}

	// Calculate dimensions based on window size
	maxWidth := m.windowWidth
	if maxWidth == 0 {
		maxWidth = 100 // Default width
	}

	// Card styling
	cardWidth := 24
	visibleCards := 3 // Force exactly 3 visible cards

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
		startIdx = max(0, m.currentPlan-1) // Center the selected card
	}
	endIdx := min(len(m.recommendI.RegionRecommendations), startIdx+visibleCards)

	cards := make([]string, endIdx-startIdx)

	for i := startIdx; i < endIdx; i++ {
		plan := m.recommendI.RegionRecommendations[i]
		isActive := i == m.currentPlan

		var borderColor, textColor, separatorColor lipgloss.Color
		var borderStyle lipgloss.Border

		if isActive {
			borderColor = lipgloss.Color("#0066FF") // Blue for active
			textColor = lipgloss.Color("#EEEEC7")   // Light color for text
			separatorColor = lipgloss.Color("#0066FF")
			borderStyle = lipgloss.RoundedBorder()
		} else {
			borderColor = lipgloss.Color("#444444") // Dark gray for inactive
			textColor = lipgloss.Color("#888888")   // Muted text for inactive
			separatorColor = lipgloss.Color("#444444")
			borderStyle = lipgloss.RoundedBorder()
		}

		priceStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Center).
			PaddingTop(1)

		priceStr := fmt.Sprintf("%s/mo\n%s", plan.price, plan.hourlyRate)
		priceSection := priceStyle.Render(priceStr)

		// Separator
		separator := lipgloss.NewStyle().
			Foreground(separatorColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			Render(strings.Repeat("─", cardWidth-6)) // Adjusted width

		// Specs section
		specsStyle := lipgloss.NewStyle().
			Foreground(textColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Left).
			PaddingLeft(1)

		specsStr := fmt.Sprintf("%s / %s\n%s\n%s",
			plan.ram, plan.cpu, plan.storage, plan.transfer)
		specsSection := specsStyle.Render(specsStr)

		// Card container style
		cardStyle := lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(borderColor).
			Width(cardWidth).
			Padding(0, 1) // Add horizontal padding inside the border

		// Put the card together
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

	instructions := "← → to navigate • enter to select plan"
	builder.WriteString(instructionStyle.Render(instructions))

	return builder.String()
}

type RegionRecommendation struct {
	t *tea.Program
}

func NewRegionRecommendation(recommendI *optimizer.RecommendationAcrossRegions) *RegionRecommendation {
	model := newModel(recommendI)
	t := tea.NewProgram(model, tea.WithAltScreen())
	return &RegionRecommendation{t: t}
}

func (t *RegionRecommendation) Run() {

	// Run the program and capture the final model state
	finalModel, err := t.t.Run()
	if err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Cast the final model to our Model type
	m, ok := finalModel.(Model)
	if !ok {
		fmt.Println("Could not get final model state")
		return
	}

	// Print the selection message after the program exits
	if m.selectedPlan >= 0 && m.selectedPlan < len(m.plans) {
		plan := m.plans[m.selectedPlan]
		fmt.Printf("\nSelected plan: %s/mo (%s)\nThank you for your selection!\n",
			plan.price, plan.hourlyRate)
	} else if m.quitting {
		fmt.Println("\nExited without selecting a plan.")
	}
}
