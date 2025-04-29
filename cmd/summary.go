package cmd

import (
	"os"
	"strconv"

	"github.com/gookit/goutil/dump"
	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Summary() *cobra.Command {

	cmd := &cobra.Command{
		Use: "summary",
		Example: `
ksctl cluster summary --help
		`,
		Short: "Use to get summary of the created cluster",
		Long:  "It is used to get summary cluster",

		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Error("No clusters found to connect")
				os.Exit(1)
			}

			selectDisplay := make(map[string]string, len(clusters))
			valueMaping := make(map[string]controller.Metadata, len(clusters))

			for idx, cluster := range clusters {
				selectDisplay[makeHumanReadableList(cluster)] = strconv.Itoa(idx)
				valueMaping[strconv.Itoa(idx)] = controller.Metadata{
					ClusterName:   cluster.Name,
					ClusterType:   cluster.ClusterType,
					Provider:      cluster.CloudProvider,
					Region:        cluster.Region,
					StateLocation: k.KsctlConfig.PreferedStateStore,
					K8sDistro:     cluster.K8sDistro,
					K8sVersion:    cluster.K8sVersion,
				}
			}

			selectedCluster, err := k.menuDriven.DropDown(
				"Select the cluster to for summary",
				selectDisplay,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterConnect, telemetry.TelemetryMeta{
				CloudProvider:     m.Provider,
				StorageDriver:     m.StateLocation,
				Region:            m.Region,
				ClusterType:       m.ClusterType,
				BootstrapProvider: m.K8sDistro,
				K8sVersion:        m.K8sVersion,
				Addons:            telemetry.TranslateMetadata(m.Addons),
			}); err != nil {
				k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
			}

			if k.loadCloudProviderCreds(m.Provider) != nil {
				os.Exit(1)
			}

			c, err := common.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: m,
				},
			)
			if err != nil {
				k.l.Error("Failed to create the controller", "Reason", err)
				os.Exit(1)
			}

			health, err := c.ClusterSummary()
			if err != nil {
				k.l.Error("Failed to connect to the cluster", "Reason", err)
				os.Exit(1)
			}
			d := dump.NewWithOptions(dump.SkipNilField(), dump.SkipPrivate())
			d.MaxDepth = 10
			d.Println(health)
		},
	}

	return cmd
}

// func printClusterSummary(summary *common.SummaryOutput) {
// 	// Format and color-code the output for better readability
// 	fmt.Printf("Cluster Summary: %s (%s)\n", summary.ClusterName, summary.OverallStatus)
// 	fmt.Printf("Kubernetes: %s %s, Provider: %s\n",
// 		summary.K8sDistro, summary.K8sVersion, summary.CloudProvider)

// 	fmt.Println("\n=== Node Status ===")
// 	fmt.Printf("Total Nodes: %d (Control Plane: %d, Workers: %d)\n",
// 		summary.NodeCounts["total"], summary.NodeCounts["master"], summary.NodeCounts["worker"])

// 	// Print nodes with issues first
// 	fmt.Println("\nNodes with issues:")
// 	hasNodeIssues := false
// 	for _, node := range summary.Nodes {
// 		if !node.Ready || node.MemoryPressure || node.DiskPressure || node.NetworkUnavailable {
// 			hasNodeIssues = true
// 			fmt.Printf("  ❌ %s - Ready: %v, MemoryPressure: %v, DiskPressure: %v, Network: %v\n",
// 				node.Name, node.Ready, node.MemoryPressure, node.DiskPressure, node.NetworkUnavailable)
// 		}
// 	}
// 	if !hasNodeIssues {
// 		fmt.Println("  No issues detected 👍")
// 	}

// 	fmt.Println("\n=== Resource Utilization ===")
// 	fmt.Printf("CPU: %.1f%% requested, %.1f%% limit\n",
// 		summary.ResourceUtilization.CPURequestPercentage,
// 		summary.ResourceUtilization.CPULimitPercentage)
// 	fmt.Printf("Memory: %.1f%% requested, %.1f%% limit\n",
// 		summary.ResourceUtilization.MemoryRequestPercentage,
// 		summary.ResourceUtilization.MemoryLimitPercentage)
// 	fmt.Printf("Pods: %d/%d (%.1f%%)\n",
// 		summary.ResourceUtilization.PodCount,
// 		summary.ResourceUtilization.PodCapacity,
// 		float64(summary.ResourceUtilization.PodCount)/float64(summary.ResourceUtilization.PodCapacity)*100)

// 	// Print other sections...

// 	if len(summary.DetectedIssues) > 0 {
// 		fmt.Println("\n=== Detected Issues ===")
// 		for _, issue := range summary.DetectedIssues {
// 			var icon string
// 			switch issue.Severity {
// 			case "Critical":
// 				icon = "🔴"
// 			case "Error":
// 				icon = "❌"
// 			case "Warning":
// 				icon = "⚠️"
// 			default:
// 				icon = "ℹ️"
// 			}
// 			fmt.Printf("%s %s - %s\n", icon, issue.Component, issue.Message)
// 			fmt.Printf("   Recommendation: %s\n", issue.Recommendation)
// 		}
// 	} else {
// 		fmt.Println("\n✅ No issues detected!")
// 	}
// }
