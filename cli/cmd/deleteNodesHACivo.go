package cmd

// authors Dipankar <dipankar@dipankar-das.com>

import (
	"os"

	control_pkg "github.com/kubesimplify/ksctl/pkg/controllers"
	"github.com/kubesimplify/ksctl/pkg/utils/consts"
	"github.com/spf13/cobra"
)

var deleteNodesHACivo = &cobra.Command{
	Use:   "delete-nodes",
	Short: "Use to delete a HA CIVO k3s cluster",
	Long: `It is used to delete cluster with the given name from user. For example:

ksctl delete-cluster ha-civo delete-nodes <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		if err := control_pkg.InitializeStorageFactory(&cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}
		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudCivo
		cli.Client.Metadata.IsHA = true

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)
		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := deleteApproval(cmd.Flags().Lookup("approve").Changed); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		err := controller.DelWorkerPlaneNode(&cli.Client)
		if err != nil {
			log.Error("Failed to scale down", "Reason", err)
			os.Exit(1)
		}
		log.Success("Scale down successful")
	},
}

func init() {
	deleteClusterHACivo.AddCommand(deleteNodesHACivo)

	clusterNameFlag(deleteNodesHACivo)
	noOfWPFlag(deleteNodesHACivo)
	regionFlag(deleteNodesHACivo)
	//k8sVerFlag(deleteNodesHACivo)
	distroFlag(deleteNodesHACivo)

	deleteNodesHACivo.MarkFlagRequired("name")
	deleteNodesHACivo.MarkFlagRequired("region")
}
