package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/spf13/cobra"
)

var selfUpdate = &cobra.Command{
	Use:   "self-update",
	Short: "update the ksctl cli",
	Long:  "setups up update for ksctl cli",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		newVer := "TODO"

		log.Success(ctx, "Updated Ksctl cli", "previousVer", Version, "newVer", newVer)
	},
}

func init() {
	RootCmd.AddCommand(selfUpdate)
	storageFlag(selfUpdate)

	selfUpdate.Flags().BoolP("verbose", "v", true, "for verbose output")
}
