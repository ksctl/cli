package cmd

/*
Kubesimplify
authors Dipankar <dipankar@dipankar-das.com>
				Anurag Kumar <contact.anurag7@gmail.com>
*/

import (
	"context"
	"github.com/ksctl/ksctl/pkg/helpers"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers/consts"

	control_pkg "github.com/ksctl/ksctl/pkg/controllers"
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
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}

		if len(provider) == 0 {
			provider = "all"
		}
		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.KsctlCloud(provider)

		err := controller.GetCluster(&cli.Client)
		if err != nil {
			log.Error("Get cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Success("Get cluster successfull")
	},
}

func init() {
	rootCmd.AddCommand(getClusterCmd)
	getClusterCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
}
