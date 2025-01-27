package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/v2/pkg/controllers"
	"github.com/ksctl/ksctl/v2/pkg/types"

	"github.com/ksctl/ksctl/v2/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var addMoreWorkerNodesHACivo = &cobra.Command{
	//Deprecated: color.HiYellowString("This will be removed in future releases once autoscaling is stable"),
	Example: `
ksctl create ha-civo add-nodes -n demo -r LON1 -s store-local --noWP 3 --nodeSizeWP g3.medium   # Here the noWP is the desired count of workernodes
	`,
	Use:   "add-nodes",
	Short: "Use to add more worker nodes in self-managed Highly-Available cluster on Civo",
	Long:  "It is used to add nodes to worker nodes in cluster with the given name from user.",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")

		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.IsHA = true
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		if err := createApproval(ctx, log, cmd.Flags().Lookup("yes").Changed); err != nil {
			log.Error("createApproval", "Reason", err)
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

		err = m.AddWorkerPlaneNodes()
		if err != nil {
			log.Error("Failed to scale up", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Scale up successful")
	},
}

func init() {
	createClusterHACivo.AddCommand(addMoreWorkerNodesHACivo)
	clusterNameFlag(addMoreWorkerNodesHACivo)
	noOfWPFlag(addMoreWorkerNodesHACivo)
	nodeSizeWPFlag(addMoreWorkerNodesHACivo)
	regionFlag(addMoreWorkerNodesHACivo)
	storageFlag(addMoreWorkerNodesHACivo)

	addMoreWorkerNodesHACivo.MarkFlagRequired("name")
	addMoreWorkerNodesHACivo.MarkFlagRequired("region")
}
