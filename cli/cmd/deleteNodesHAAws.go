package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var deleteNodesHAAws = &cobra.Command{
	//Deprecated: color.HiYellowString("This will be removed in future releases once autoscaling is stable"),
	Use: "del-nodes",
	Example: `
ksctl delete ha-aws del-nodes -n demo -r us-east-1 -s store-local --noWP 1      # Here the noWP is the desired count of workernodes
	`,
	Short: "Use to remove worker nodes in self-managed Highly-Available cluster on Aws",
	Long:  "It is used to delete cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws
		cli.Client.Metadata.IsHA = true

		SetDefaults(consts.CloudAws, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		if err := deleteApproval(ctx, log, cmd.Flags().Lookup("yes").Changed); err != nil {
			log.Error("deleteApproval", "Reason", err)
			os.Exit(1)
		}

		m, err := controllers.NewManagerClusterSelfManaged(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error("Failed to init manager", "Reason", err)
			os.Exit(1)
		}

		err = m.DelWorkerPlaneNodes()
		if err != nil {
			log.Error("Failed to scale down", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Scale down successful")
	},
}

func init() {
	deleteClusterHAAws.AddCommand(deleteNodesHAAws)

	clusterNameFlag(deleteNodesHAAws)
	noOfWPFlag(deleteNodesHAAws)
	regionFlag(deleteNodesHAAws)
	storageFlag(deleteNodesHAAws)

	deleteNodesHAAws.MarkFlagRequired("name")
	deleteNodesHAAws.MarkFlagRequired("region")
}
