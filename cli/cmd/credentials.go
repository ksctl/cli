package cmd

import (
	"fmt"
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var credCmd = &cobra.Command{
	Use:   "cred",
	Short: "Login to your Cloud-provider Credentials",
	Long:  "login to your cloud provider credentials",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		log.Print(ctx, `
1> AWS (EKS)
2> Azure (AKS)
3> Civo (K3s)
`)

		choice := 0

		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			panic(err.Error())
		}
		if provider, ok := cloud[choice]; ok {
			cli.Client.Metadata.Provider = consts.KsctlCloud(provider)
		} else {
			log.Error(ctx, "invalid provider")
		}
		m, err := controllers.NewManagerClusterKsctl(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error(ctx, "Failed to initialize", "Reason", err)
			os.Exit(1)
		}

		if err := m.Credentials(); err != nil {
			log.Error(ctx, "Failed to added the credential", "Reason", err)
			os.Exit(1)
		}
		log.Success(ctx, "Credentials added successfully")
	},
}

func init() {
	RootCmd.AddCommand(credCmd)
	storageFlag(credCmd)

	credCmd.Flags().BoolP("verbose", "v", true, "for verbose output")

}
