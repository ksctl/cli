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
		cli.Client.Metadata.Provider = consts.KsctlCloud(provider)
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

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
