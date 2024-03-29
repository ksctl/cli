package cmd

// authors Dipankar <dipankar@dipankar-das.com>

import (
	"context"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers"

	"github.com/spf13/cobra"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
)

var addMoreWorkerNodesHAAzure = &cobra.Command{
	Use:   "add-nodes",
	Short: "Use to add more worker nodes in HA azure k3s cluster",
	Long: `It is used to add nodes to worker nodes in cluster with the given name from user. For example:

ksctl create-cluster ha-azure add-nodes <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		cli.Client.Metadata.NoWP = noWP
		cli.Client.Metadata.WorkerPlaneNodeType = nodeSizeWP
		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.IsHA = true
		cli.Client.Metadata.K8sDistro = consts.KsctlKubernetes(distro)
		cli.Client.Metadata.K8sVersion = k8sVer

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

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
	createClusterHAAzure.AddCommand(addMoreWorkerNodesHAAzure)
	clusterNameFlag(addMoreWorkerNodesHAAzure)
	noOfWPFlag(addMoreWorkerNodesHAAzure)
	nodeSizeWPFlag(addMoreWorkerNodesHAAzure)
	regionFlag(addMoreWorkerNodesHAAzure)
	k8sVerFlag(addMoreWorkerNodesHAAzure)
	distroFlag(addMoreWorkerNodesHAAzure)
	storageFlag(addMoreWorkerNodesHAAzure)

	addMoreWorkerNodesHAAzure.MarkFlagRequired("name")
	addMoreWorkerNodesHAAzure.MarkFlagRequired("region")
}
