package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"

	"github.com/spf13/cobra"
)

var getClusterCmd = &cobra.Command{
	Use:     "get-clusters",
	Aliases: []string{"get", "list"},
	Example: `
ksctl get --provider all --storage store-local
`,
	Short: "Use to get clusters",
	Long: LongMessage(`It is used to view clusters. For example:

ksctl get-clusters `),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}

		if len(provider) == 0 {
			provider = "all"
		}
		SetRequiredFeatureFlags(ctx, log, cmd)
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

		err = m.GetCluster()
		if err != nil {
			log.Error("Get cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Get cluster successfull")
	},
}

func init() {
	RootCmd.AddCommand(getClusterCmd)
	storageFlag(getClusterCmd)

	getClusterCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
}
