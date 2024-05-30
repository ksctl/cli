package cmd

/*
Kubesimplify
authors Dipankar <dipankar@dipankar-das.com>
				Anurag Kumar <contact.anurag7@gmail.com>
*/

import (
	"os"

	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"

	"github.com/spf13/cobra"
)

type printer struct {
	ClusterName string `json:"cluster_name"`
	Region      string `json:"region"`
	Provider    string `json:"provider"`
}

// viewClusterCmd represents the viewCluster command
var getClusterCmd = &cobra.Command{
	Use:     "get-clusters",
	Aliases: []string{"get"},
	Short:   "Use to get clusters",
	Long: `It is used to view clusters. For example:

ksctl get-clusters `,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)

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
			log.Error(ctx, "failed to init", "Reason", err)
			os.Exit(1)
		}

		err = m.GetCluster()
		if err != nil {
			log.Error(ctx, "Get cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Get cluster successfull")
	},
}

func init() {
	rootCmd.AddCommand(getClusterCmd)
	storageFlag(getClusterCmd)

	getClusterCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
}
