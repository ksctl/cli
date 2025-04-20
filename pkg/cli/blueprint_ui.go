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
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
)

// BlueprintUI is responsible for rendering the cluster blueprint with enhanced UI
type BlueprintUI struct {
	writer io.Writer
}

// NewBlueprintUI creates a new instance of BlueprintUI
func NewBlueprintUI(w io.Writer) *BlueprintUI {
	return &BlueprintUI{
		writer: w,
	}
}

// RenderClusterBlueprint renders the cluster metadata with enhanced UI
func (ui *BlueprintUI) RenderClusterBlueprint(meta controller.Metadata) {
	// Border styles
	banner := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")).
		Padding(0, 1).
		Width(50).
		Align(lipgloss.Center)

	sectionTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("13")).
		Bold(true).
		MarginTop(1).
		Padding(0, 1)

	infoBlock := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")).
		Padding(1, 2).
		MarginTop(1).
		Width(50)

	keyValueRow := func(key, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Foreground(lipgloss.Color("14")).PaddingRight(3).Width(22).Align(lipgloss.Left).Render(key),
			lipgloss.NewStyle().Width(40).Render(value),
		)
	}

	// Header banner
	bannerContent := fmt.Sprintf("✨ %s ✨\n\n%s",
		lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#FFFFFF")).Render("Cluster Blueprint"),
		lipgloss.NewStyle().Italic(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#DDDDDD")).Render("Your Kubernetes cluster plan"))

	fmt.Fprintln(ui.writer, banner.Render(bannerContent))
	fmt.Fprintln(ui.writer)

	// Define a container for all sections
	gridContainer := lipgloss.NewStyle().
		Width(120). // Double the width to accommodate two sections side by side
		Align(lipgloss.Center)

	// Create slices to hold our sections and their titles
	var sectionBlocks []string

	// Key Attributes section
	{
		var content strings.Builder
		content.WriteString(keyValueRow("🔖 Cluster Name", meta.ClusterName))
		content.WriteString("\n")
		content.WriteString(keyValueRow("📍 Region", meta.Region))
		content.WriteString("\n")
		content.WriteString(keyValueRow("🌐 Cloud Provider", string(meta.Provider)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("🔧 Cluster Type", string(meta.ClusterType)))

		contentBlock := infoBlock.Render(content.String())
		titleBlock := sectionTitle.Render("🔑 Key Attributes")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		sectionBlocks = append(sectionBlocks, fullSection)
	}

	// Infrastructure section
	if meta.NoCP > 0 || meta.NoWP > 0 || meta.NoDS > 0 || len(meta.ManagedNodeType) > 0 {
		var content strings.Builder
		if meta.NoCP > 0 {
			content.WriteString(keyValueRow("🎮 Control Plane", fmt.Sprintf("%d × %s", meta.NoCP, color.HiMagentaString(meta.ControlPlaneNodeType))))
			content.WriteString("\n")
		}
		if meta.NoWP > 0 {
			content.WriteString(keyValueRow("🔋 Worker Nodes", fmt.Sprintf("%d × %s", meta.NoWP, color.HiMagentaString(meta.WorkerPlaneNodeType))))
			content.WriteString("\n")
		}
		if meta.NoDS > 0 {
			content.WriteString(keyValueRow("💾 Etcd Nodes", fmt.Sprintf("%d × %s", meta.NoDS, color.HiMagentaString(meta.DataStoreNodeType))))
			content.WriteString("\n")
		}
		if meta.LoadBalancerNodeType != "" {
			content.WriteString(keyValueRow("🔀 Load Balancer", color.HiMagentaString(meta.LoadBalancerNodeType)))
			content.WriteString("\n")
		}
		if len(meta.ManagedNodeType) > 0 {
			content.WriteString(keyValueRow("🌐 Managed Nodes", fmt.Sprintf("%d × %s", meta.NoMP, color.HiMagentaString(meta.ManagedNodeType))))
		}

		// Trim trailing newline if present
		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("🖥️  Infrastructure")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		sectionBlocks = append(sectionBlocks, fullSection)
	}

	// Kubernetes configuration section
	if meta.K8sDistro != "" || meta.EtcdVersion != "" || meta.K8sVersion != "" {
		var content strings.Builder
		if meta.K8sDistro != "" {
			content.WriteString(keyValueRow("🚀 Bootstrap Provider", string(meta.K8sDistro)))
			content.WriteString("\n")
		}
		if meta.K8sVersion != "" {
			content.WriteString(keyValueRow("🔄 Kubernetes Version", meta.K8sVersion))
			content.WriteString("\n")
		}
		if meta.EtcdVersion != "" {
			content.WriteString(keyValueRow("📦 Etcd Version", meta.EtcdVersion))
		}

		// Trim trailing newline if present
		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("⚙️  Kubernetes Configuration")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		sectionBlocks = append(sectionBlocks, fullSection)
	}

	// Addons section
	if len(meta.Addons) > 0 {
		var sectionContent strings.Builder
		addonBlock := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("10")).
			Padding(1, 2).
			MarginTop(1).
			Width(50)

		for i, addon := range meta.Addons {
			addonTitle := color.HiMagentaString(addon.Name)

			config := addon.Config
			vConfig := ""
			if config == nil {
				vConfig = "No configuration available"
			} else {
				var v any
				if err := json.Unmarshal([]byte(*config), &v); err != nil {
					vConfig = "Invalid configuration format"
				} else {
					_v, _ := json.MarshalIndent(v, "\t", "  ")
					vConfig = string(_v)
				}
			}
			addonInfo := fmt.Sprintf("%s\n\t%s: %s\n\t%s: %s",
				addonTitle,
				color.HiCyanString("From"),
				color.HiGreenString(addon.Label),
				color.HiCyanString("Config"),
				color.HiGreenString(vConfig),
			)
			if addon.IsCNI {
				addonInfo += "\n\t" + color.HiCyanString("CNI Add-on")
			}

			sectionContent.WriteString(addonBlock.Render(addonInfo))
			if i < len(meta.Addons)-1 {
				sectionContent.WriteString("\n")
			}
		}

		titleBlock := sectionTitle.Render("🧩 Cluster Add-ons")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, sectionContent.String())

		sectionBlocks = append(sectionBlocks, fullSection)
	}

	// Render the grid layout with 2 sections per row
	var gridRows []string

	// Process sections in pairs
	for i := 0; i < len(sectionBlocks); i += 2 {
		row := ""

		if i+1 < len(sectionBlocks) {
			// If we have a pair, join them horizontally with padding
			left := sectionBlocks[i]
			right := sectionBlocks[i+1]

			// Create padding between columns
			spacing := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render("")

			row = lipgloss.JoinHorizontal(
				lipgloss.Top, // Align top edges of sections
				left,
				spacing,
				right,
			)
		} else {
			// If we have an odd number of sections, center the last one
			row = lipgloss.NewStyle().Align(lipgloss.Center).Render(sectionBlocks[i])
		}

		gridRows = append(gridRows, row)
	}

	// Join all rows vertically with padding
	finalGrid := lipgloss.JoinVertical(
		lipgloss.Center,
		gridRows...,
	)

	fmt.Fprintln(ui.writer, gridContainer.Render(finalGrid))

	// Footer note
	noteStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#919191")).
		Padding(1, 0).
		MarginTop(1).
		Align(lipgloss.Center)

	fmt.Fprintln(ui.writer)
	fmt.Fprintln(ui.writer, noteStyle.Render("Your cluster will be provisioned with these specifications"))
}
