package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"
	"github.com/pterm/pterm"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var credCmd = &cobra.Command{
	Use:   "cred",
	Short: "Login to your Cloud-provider Credentials",
	Long:  LongMessage("login to your cloud provider credentials"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		cloudMap := map[string]int{
			"Amazon Web Services": 1,
			"Azure":               2,
			"Civo":                3,
		}
		var options []string
		for k := range cloudMap {
			options = append(options, k)
		}
		log.Print(ctx, "Select the cloud provider")

		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		if provider, ok := cloud[cloudMap[selectedOption]]; ok {
			cli.Client.Metadata.Provider = consts.KsctlCloud(provider)
		} else {
			log.Error("invalid provider")
		}
		m, err := controllers.NewManagerClusterKsctl(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error("Failed to initialize", "Reason", err)
			os.Exit(1)
		}

		if err := m.Credentials(); err != nil {
			log.Error("Failed to added the credential", "Reason", err)
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
