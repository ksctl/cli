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
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
)

// SummaryUI is responsible for rendering cluster summary with enhanced UI
type SummaryUI struct {
	writer io.Writer
}

// NewSummaryUI creates a new instance of SummaryUI
func NewSummaryUI(w io.Writer) *SummaryUI {
	return &SummaryUI{
		writer: w,
	}
}

// RenderClusterSummary renders the cluster summary with enhanced UI
func (ui *SummaryUI) RenderClusterSummary(summary *common.SummaryOutput) {
	parentBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2).
		Width(90).
		Align(lipgloss.Center)

	banner := lipgloss.NewStyle().
		Padding(0, 1).
		Width(80).
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
		Width(80)

	keyValueRow := func(key, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Foreground(lipgloss.Color("14")).PaddingRight(3).Width(28).Align(lipgloss.Left).Render(key),
			lipgloss.NewStyle().Width(50).Render(value),
		)
	}

	var parentBoxContent strings.Builder

	bannerContent := fmt.Sprintf("✨ %s ✨\n\n%s",
		lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#FFFFFF")).Render("Cluster Summary"),
		lipgloss.NewStyle().Italic(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#DDDDDD")).Render("Health and status of your Kubernetes cluster"))

	parentBoxContent.WriteString(banner.Render(bannerContent))
	parentBoxContent.WriteString("\n\n")

	// Cluster basics section
	{
		var content strings.Builder

		if summary.ClusterName != "" {
			content.WriteString(keyValueRow("Name", summary.ClusterName))
			content.WriteString("\n")
		}
		if summary.CloudProvider != "" {
			content.WriteString(keyValueRow("Cloud", summary.CloudProvider))
			content.WriteString("\n")
		}
		if summary.ClusterType != "" {
			content.WriteString(keyValueRow("Type", summary.ClusterType))
			content.WriteString("\n")
		}
		if summary.RoundTripLatency != "" {
			content.WriteString(keyValueRow("Round Trip Latency", summary.RoundTripLatency))
			content.WriteString("\n")
		}
		if summary.KubernetesVersion != "" {
			content.WriteString(keyValueRow("Kubernetes Version", summary.KubernetesVersion))
			content.WriteString("\n")
		}

		if summary.APIServerHealthCheck != nil {
			if !summary.APIServerHealthCheck.Healthy {
				content.WriteString(keyValueRow("API Server Health", color.HiRedString("Unhealthy")))
			} else {
				content.WriteString(keyValueRow("API Server Health", color.HiGreenString("Healthy")))
			}
			content.WriteString("\n")
			if len(summary.APIServerHealthCheck.FailedComponents) > 0 {
				v := color.HiRedString(strings.Join(summary.APIServerHealthCheck.FailedComponents, ", "))
				content.WriteString(keyValueRow("Components Unhealthy", v))
				content.WriteString("\n")
			}
		}
		for k, v := range summary.ControlPlaneComponentVers {
			content.WriteString(keyValueRow(k, v))
			content.WriteString("\n")
		}

		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("🔑 Key Attributes")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		parentBoxContent.WriteString(fullSection)
		parentBoxContent.WriteString("\n")
	}
	yesNoColor := func(v bool, inverse bool) (string, func(format string, a ...interface{}) string) {
		if v {
			if inverse {
				return "Yes", color.HiRedString
			}
			return "Yes", color.HiGreenString
		}
		if inverse {
			return "No", color.HiGreenString
		}
		return "No", color.HiRedString
	}
	applyYesNoColor := func(msg string, formatter func(format string, a ...interface{}) string) string {
		return formatter(msg)
	}

	{
		// Nodes section
		var content strings.Builder
		for _, node := range summary.Nodes {
			content.WriteString("\n")
			content.WriteString(color.HiMagentaString(node.Name))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Ready", applyYesNoColor(yesNoColor(node.Ready, false))))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Kubelet Healthy", applyYesNoColor(yesNoColor(node.KubeletHealthy, false))))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Memory Pressure", applyYesNoColor(yesNoColor(node.MemoryPressure, true))))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Disk Pressure", applyYesNoColor(yesNoColor(node.DiskPressure, true))))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Unreachable", applyYesNoColor(yesNoColor(node.NetworkUnavailable, true))))
			content.WriteString("\n")

			content.WriteString(keyValueRow(" Kubelet Version", node.KubeletVersion))
			content.WriteString("\n")

			content.WriteString(keyValueRow(" CRI", node.ContainerRuntimeVersion))
			content.WriteString("\n")

			cpuUtlization := "😢 Currently Unavailable"
			if node.CPUUtilization > 0 {
				cpuUtlization = fmt.Sprintf("%.2f%% (Total: %s)", node.CPUUtilization, node.CPUUnits)
			}
			memoryUtlization := "😢 Currently Unavailable"
			if node.MemoryUtilization > 0 {
				memoryUtlization = fmt.Sprintf("%.2f%% (Total: %s)", node.MemoryUtilization, node.MemUnits)
			}
			content.WriteString(keyValueRow(" CPU Usage", cpuUtlization))
			content.WriteString("\n")
			content.WriteString(keyValueRow(" Memory Usage", memoryUtlization))
			content.WriteString("\n")
		}

		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("🤖 Nodes")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		parentBoxContent.WriteString(fullSection)
		parentBoxContent.WriteString("\n")
	}
	{
		// Workload section
		var content strings.Builder
		workloads := summary.WorkloadSummary
		content.WriteString(keyValueRow("Deployments", fmt.Sprintf("%d", workloads.Deployments)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("StatefulSets", fmt.Sprintf("%d", workloads.StatefulSets)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("DaemonSets", fmt.Sprintf("%d", workloads.DaemonSets)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("CronJobs", fmt.Sprintf("%d", workloads.CronJobs)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Namespaces", fmt.Sprintf("%d", workloads.Namespaces)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Persistent Volumes", fmt.Sprintf("%d", workloads.PV)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Persistent Volume Claims", fmt.Sprintf("%d", workloads.PVC)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Storage Class", fmt.Sprintf("%d", workloads.SC)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Service (ClusterIP)", fmt.Sprintf("%d", workloads.ClusterIPSVC)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Service (LoadBalancer)", fmt.Sprintf("%d", workloads.LoadbalancerSVC)))
		content.WriteString("\n")
		content.WriteString(keyValueRow("Pods (Running)", fmt.Sprintf("%d", workloads.RunningPods)))
		content.WriteString("\n")

		if len(workloads.UnHealthyPods) > 0 {
			content.WriteString(color.HiMagentaString("Unhealthy Pods"))
			content.WriteString("\n")
		}

		for _, pod := range workloads.UnHealthyPods {
			content.WriteString("\n")
			content.WriteString(" " + color.HiRedString(pod.Name+"@"+pod.Namespace))
			content.WriteString("\n")
			objectRef := []string{}
			for _, ref := range pod.OwnerRef {
				v := fmt.Sprintf("%s/%s@%s", ref.Kind, ref.Name, pod.Namespace)
				objectRef = append(objectRef, v)
			}
			if len(objectRef) == 0 {
				objectRef = append(objectRef, "No Owner Reference")
			}
			content.WriteString(keyValueRow("  "+"Owner Ref", strings.Join(objectRef, ", ")))
			content.WriteString("\n")
			content.WriteString(keyValueRow("  "+"Failed", applyYesNoColor(yesNoColor(pod.IsFailed, true))))
			content.WriteString("\n")
			content.WriteString(keyValueRow("  "+"Pending", applyYesNoColor(yesNoColor(pod.IsPending, true))))
			content.WriteString("\n")
			for _, container := range pod.FailedContainers {
				content.WriteString("\n")
				content.WriteString(keyValueRow("  "+"Container", container.Name))
				content.WriteString("\n")
				content.WriteString(keyValueRow("   "+"Restart Count", fmt.Sprintf("%d", container.RestartCount)))
				content.WriteString("\n")
				content.WriteString(keyValueRow("   "+"Reason", container.WaitingProblem.Reason))
				content.WriteString("\n")
				content.WriteString(keyValueRow("   "+"Message", container.WaitingProblem.Message))
				content.WriteString("\n")
			}
			content.WriteString("\n")
		}

		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("📦 Workloads")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		parentBoxContent.WriteString(fullSection)
		parentBoxContent.WriteString("\n")
	}
	{
		// Recent events section
		var content strings.Builder
		if len(summary.RecentWarningEvents) == 0 {
			content.WriteString(color.HiGreenString("No recent events"))
			content.WriteString("\n")
		} else {
			content.WriteString(color.HiMagentaString("Recent Warning Events (Past 24h)"))
			content.WriteString("\n")
			for _, event := range summary.RecentWarningEvents {
				content.WriteString("\n")
				content.WriteString(color.HiMagentaString(fmt.Sprintf("%s/%s@%s", event.Kind, event.Name, event.Namespace)))
				content.WriteString("\n")
				content.WriteString(keyValueRow("  Reason", event.Reason))
				content.WriteString("\n")
				content.WriteString(keyValueRow("  Message", event.Message))
				content.WriteString("\n")
				content.WriteString(keyValueRow("  Happened On", event.Time.String()))
				content.WriteString("\n")
				content.WriteString(keyValueRow("  Reported By", event.ReportedBy))
				content.WriteString("\n")
			}
		}

		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("📢 Events")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		parentBoxContent.WriteString(fullSection)
		parentBoxContent.WriteString("\n")
	}
	{
		// Issues section
		var content strings.Builder
		if len(summary.DetectedIssues) == 0 {
			content.WriteString(color.HiGreenString("No issues detected"))
			content.WriteString("\n")
		} else {
			content.WriteString(color.HiMagentaString("Detected Issues"))
			content.WriteString("\n")
			severityColor := func(severity string) string {
				switch severity {
				case "Critical":
					return color.New(color.BgRed, color.FgBlack).Add(color.Bold).Sprintf("%s", severity)
				case "Error":
					return color.HiRedString(severity)
				case "Warning":
					return color.HiYellowString(severity)
				default:
					return severity
				}
			}
			for _, issue := range summary.DetectedIssues {
				content.WriteString(color.HiMagentaString(issue.Component))
				content.WriteString("\n")
				content.WriteString(keyValueRow(" Severity", severityColor(issue.Severity)))
				content.WriteString("\n")
				content.WriteString(keyValueRow(" Message", issue.Message))
				content.WriteString("\n")
				content.WriteString(keyValueRow(" Recommendation", issue.Recommendation))
				content.WriteString("\n")
			}
		}

		contentStr := strings.TrimSuffix(content.String(), "\n")
		contentBlock := infoBlock.Render(contentStr)
		titleBlock := sectionTitle.Render("🚩 Problems")
		fullSection := lipgloss.JoinVertical(lipgloss.Left, titleBlock, contentBlock)

		parentBoxContent.WriteString(fullSection)
		parentBoxContent.WriteString("\n")
	}

	fmt.Fprintln(ui.writer, parentBox.Render(parentBoxContent.String()))
}
