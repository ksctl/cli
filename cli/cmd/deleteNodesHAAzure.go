package cmd

// authors Dipankar <dipankar.das@ksctl.com>

import (
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var deleteNodesHAAzure = &cobra.Command{
	Deprecated: color.HiYellowString("This will be removed in future releases once autoscaling is stable"),
	Use:        "del-nodes",
	Short:      "Use to delete a HA azure k3s cluster",
	Long: `It is used to delete cluster with the given name from user. For example:

ksctl delete-cluster ha-azure delete-nodes <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure
		cli.Client.Metadata.IsHA = true

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)

		if err := deleteApproval(ctx, log, cmd.Flags().Lookup("approve").Changed); err != nil {
			log.Error(ctx, "deleteApproval", "Reason", err)
			os.Exit(1)
		}

		m, err := controllers.NewManagerClusterSelfManaged(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error(ctx, "Failed to init manager", "Reason", err)
			os.Exit(1)
		}

		err = m.DelWorkerPlaneNodes()
		if err != nil {
			log.Error(ctx, "Failed to scale down", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Scale down successful")
	},
}

func init() {
	deleteClusterHAAzure.AddCommand(deleteNodesHAAzure)

	clusterNameFlag(deleteNodesHAAzure)
	noOfWPFlag(deleteNodesHAAzure)
	regionFlag(deleteNodesHAAzure)
	distroFlag(deleteNodesHAAzure)
	storageFlag(deleteNodesHAAzure)

	deleteNodesHAAzure.MarkFlagRequired("name")
	deleteNodesHAAzure.MarkFlagRequired("region")
}
