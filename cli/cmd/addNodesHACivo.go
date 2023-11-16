package cmd

// authors Dipankar <dipankar@dipankar-das.com>

import (
	"os"

	control_pkg "github.com/kubesimplify/ksctl/pkg/controllers"

	"github.com/kubesimplify/ksctl/pkg/utils/consts"
	"github.com/spf13/cobra"
)

var addMoreWorkerNodesHACivo = &cobra.Command{
	Use:   "add-nodes",
	Short: "Use to add more worker nodes in HA CIVO k3s cluster",
	Long: `It is used to add nodes to worker nodes in cluster with the given name from user. For example:

ksctl create-cluster ha-civo add-nodes <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		if err := control_pkg.InitializeStorageFactory(&cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}
		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudCivo
		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)
		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
		cli.Client.Metadata.K8sVersion = k8sVer

		cli.Client.Metadata.IsHA = true

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := createApproval(cmd.Flags().Lookup("approve").Changed); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		err := controller.AddWorkerPlaneNode(&cli.Client)
		if err != nil {
			log.Error("Failed to scale up", "Reason", err)
			os.Exit(1)
		}
		log.Success("Scale up successful")
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

	addMoreWorkerNodesHACivo.MarkFlagRequired("name")
	addMoreWorkerNodesHACivo.MarkFlagRequired("region")
}
