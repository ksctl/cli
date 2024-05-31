package cmd

// authors Dipankar <dipankar.das@ksctl.com>

import (
	"os"

	"github.com/fatih/color"

	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var addMoreWorkerNodesHACivo = &cobra.Command{
	Deprecated: color.HiYellowString("This will be removed in future releases once autoscaling is stable"),
	Use:        "add-nodes",
	Short:      "Use to add more worker nodes in HA CIVO k3s cluster",
	Long: `It is used to add nodes to worker nodes in cluster with the given name from user. For example:

ksctl create-cluster ha-civo add-nodes <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")

		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
		cli.Client.Metadata.K8sVersion = k8sVer
		cli.Client.Metadata.IsHA = true
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		if err := createApproval(ctx, log, cmd.Flags().Lookup("approve").Changed); err != nil {
			log.Error(ctx, "createApproval", "Reason", err)
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

		err = m.AddWorkerPlaneNodes()
		if err != nil {
			log.Error(ctx, "Failed to scale up", "Reason", err)
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
	k8sVerFlag(addMoreWorkerNodesHACivo)
	distroFlag(addMoreWorkerNodesHACivo)
	storageFlag(addMoreWorkerNodesHACivo)

	addMoreWorkerNodesHACivo.MarkFlagRequired("name")
	addMoreWorkerNodesHACivo.MarkFlagRequired("region")
}
