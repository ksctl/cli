package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"

	"github.com/spf13/cobra"
)

var infoClusterCmd = &cobra.Command{
	Use:     "info-cluster",
	Aliases: []string{"info"},
	Example: `
ksctl info --provider azure --name demo --region eastus --storage store-local
ksctl info -p ha-azure -n ha-demo-kubeadm -r eastus -s store-local --verbose -1
`,
	Short: "Use to info cluster",
	Long:  `It is used to detailed data for a given cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

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

		case string(consts.CloudAws):
			cli.Client.Metadata.Provider = consts.CloudAws

		case string(consts.CloudAzure):
			cli.Client.Metadata.Provider = consts.CloudAzure
		default:
			log.Error("invalid provider specified", "provider", provider)
			os.Exit(1)
		}

		m, err := controllers.NewManagerClusterKsctl(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error("failed to init", "Reason", err)
			os.Exit(1)
		}

		_, err = m.InfoCluster()
		if err != nil {
			log.Error("info cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "info cluster successfull")
	},
}

func init() {
	RootCmd.AddCommand(infoClusterCmd)
	storageFlag(infoClusterCmd)
	clusterNameFlag(infoClusterCmd)
	regionFlag(infoClusterCmd)

	infoClusterCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
	infoClusterCmd.MarkFlagRequired("name")
	infoClusterCmd.MarkFlagRequired("provider")
}
