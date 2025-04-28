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
			dump.P(health)
		},
	}

	return cmd
}
