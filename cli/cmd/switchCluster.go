package cmd

// authors Dipankar <dipankar@dipankar-das.com>

import (
	"os"

	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var switchCluster = &cobra.Command{
	Use:     "switch-cluster",
	Aliases: []string{"switch"},
	Short:   "Use to switch between clusters",
	Long: `It is used to switch cluster with the given ClusterName from user. For example:

ksctl switch-context -p <civo,local,ha-civo,ha-azure,ha-aws,azure>  -n <clustername> -r <region> <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)

		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		switch provider {
		case string(consts.CloudLocal):
			cli.Client.Metadata.Provider = consts.CloudLocal

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudCivo):
			cli.Client.Metadata.Provider = consts.CloudCivo
			cli.Client.Metadata.IsHA = true

		case string(consts.CloudCivo):
			cli.Client.Metadata.Provider = consts.CloudCivo

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudAzure):
			cli.Client.Metadata.Provider = consts.CloudAzure
			cli.Client.Metadata.IsHA = true

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudAws):
			cli.Client.Metadata.Provider = consts.CloudAws
			cli.Client.Metadata.IsHA = true

		case string(consts.CloudAzure):
			cli.Client.Metadata.Provider = consts.CloudAzure
		}

		m, err := controllers.NewManagerClusterKsctl(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error(ctx, "failed to init", "Reason", err)
			os.Exit(1)
		}
		kubeconfig, err := m.SwitchCluster()
		if err != nil {
			log.Error(ctx, "Switch cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Debug(ctx, "kubeconfig output as string", "kubeconfig", kubeconfig)
		log.Success(ctx, "Switch cluster Successful")
	},
}

func init() {
	rootCmd.AddCommand(switchCluster)
	clusterNameFlag(switchCluster)
	regionFlag(switchCluster)
	storageFlag(switchCluster)

	switchCluster.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
	switchCluster.MarkFlagRequired("name")
	switchCluster.MarkFlagRequired("provider")
}
