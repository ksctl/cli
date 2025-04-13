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
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (p *genericMenuDriven) CardSelection(pack CardPack) (string, error) {
	t := tea.NewProgram(newCardRunner(pack))
	finalModel, err := t.Run()
	if err != nil {
		return "", err
	}

	m, ok := finalModel.(cardRunner)
	if !ok {
		return "", fmt.Errorf("failed to cast final model to cardRunner type")
	}

	if m.selectedPlan >= 0 && m.selectedPlan < pack.LenOfItems() {
		plan := pack.GetResult(m.selectedPlan)
		return plan, nil
	}

	if m.quitting {
		return "", nil
	}
	return "", fmt.Errorf("internal problem. invalid selected plan index: %d", m.selectedPlan)
}

type cardRunner struct {
	b            CardPack
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

func newCardRunner(b CardPack) cardRunner {
	return cardRunner{
		b:            b,
		currentPlan:  0,
		help:         help.New(),
		keys:         newKeyMap(),
		selectedPlan: -1,
	}
}

func (m cardRunner) Init() tea.Cmd {
	return nil
}

func (m cardRunner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.currentPlan < m.b.LenOfItems()-1 {
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

func (m cardRunner) View() string {
	if m.quitting {
		if m.selectedPlan >= 0 && m.selectedPlan < m.b.LenOfItems() {
			return m.b.GetResult(m.selectedPlan)
		}
		return ""
	}

	maxWidth := m.windowWidth
	if maxWidth == 0 {
		maxWidth = 100 // Default width
	}

	cardWidth, visibleCards := m.b.GetCardConfiguration()

	// Add some top padding
	var builder strings.Builder
	builder.WriteString("\n\n")

	// Determine which cards to show
	var startIdx int

	implLen := m.b.LenOfItems()
	if m.currentPlan == 0 {
		startIdx = 0
	} else if m.currentPlan == implLen-1 {
		startIdx = max(0, implLen-visibleCards)
	} else {
		startIdx = max(0, m.currentPlan-1)
	}
	endIdx := min(implLen, startIdx+visibleCards)

	cards := make([]string, endIdx-startIdx)

	for i := startIdx; i < endIdx; i++ {
		plan := m.b.GetItem(i)
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

		priceSection := priceStyle.Render(plan.GetUpper())

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

		specsSection := specsStyle.Render(plan.GetLower())

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

	builder.WriteString(instructionStyle.Render(m.b.GetInstruction() + "\n\n"))

	return builder.String()
}
