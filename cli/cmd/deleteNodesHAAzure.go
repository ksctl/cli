package cmd

// authors Dipankar <dipankar@dipankar-das.com>

import (
	"context"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var deleteNodesHAAzure = &cobra.Command{
	Use:   "del-nodes",
	Short: "Use to delete a HA azure k3s cluster",
	Long: `It is used to delete cluster with the given name from user. For example:

ksctl delete-cluster ha-azure delete-nodes <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure
		cli.Client.Metadata.IsHA = true

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

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
	deleteClusterHAAzure.AddCommand(deleteNodesHAAzure)

	clusterNameFlag(deleteNodesHAAzure)
	noOfWPFlag(deleteNodesHAAzure)
	regionFlag(deleteNodesHAAzure)
	distroFlag(deleteNodesHAAzure)
	storageFlag(deleteNodesHAAzure)

	deleteNodesHAAzure.MarkFlagRequired("name")
	deleteNodesHAAzure.MarkFlagRequired("region")
}
