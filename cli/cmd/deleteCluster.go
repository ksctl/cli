package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/types"
	"github.com/spf13/cobra"
)

var deleteClusterCmd = &cobra.Command{
	Use: "delete-cluster",
	Example: `
ksctl delete --help
	`,
	Short:   "Use to delete a cluster",
	Aliases: []string{"delete"},
	Long:    LongMessage("It is used to delete cluster of given provider"),
}

var deleteClusterLocal = &cobra.Command{
	Use: "local",
	Example: `
ksctl delete local --name demo --storage store-local
`,
	Short: "Use to delete a kind cluster",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudLocal

		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterAzure = &cobra.Command{
	Use: "azure",
	Example: `
ksctl delete azure --name demo --region eastus --storage store-local
`,
	Short: "Use to deletes a AKS cluster",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterAws = &cobra.Command{
	Use: "aws",
	Example: `
ksctl delete aws --name demo --region ap-south-1 --storage store-local
`,
	Short: "Use to deletes a EKS cluster",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterCivo = &cobra.Command{
	Use: "civo",
	Example: `
ksctl delete civo --name demo --region LON1 --storage store-local
`,
	Short: "Use to delete a Civo managed k3s cluster",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterHAAws = &cobra.Command{
	Use: "ha-aws",
	Example: `
ksctl delete ha-aws --name demo --region us-east-1 --storage store-local
`,
	Short: "Use to delete a self-managed Highly Available cluster on AWS",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterHAAzure = &cobra.Command{
	Use: "ha-azure",
	Example: `
ksctl delete ha-azure --name demo --region eastus --storage store-local
`,
	Short: "Use to delete a self-managed Highly Available cluster on Azure",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var deleteClusterHACivo = &cobra.Command{
	Use: "ha-civo",
	Example: `
ksctl delete ha-civo --name demo --region LON1 --storage store-local
`,
	Short: "Use to delete a self-managed Highly Available cluster on Civo",
	Long:  LongMessage("It is used to delete cluster of given provider"),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

func init() {
	RootCmd.AddCommand(deleteClusterCmd)

	deleteClusterCmd.AddCommand(deleteClusterHACivo)
	deleteClusterCmd.AddCommand(deleteClusterCivo)
	deleteClusterCmd.AddCommand(deleteClusterHAAzure)
	deleteClusterCmd.AddCommand(deleteClusterAzure)
	deleteClusterCmd.AddCommand(deleteClusterLocal)
	deleteClusterCmd.AddCommand(deleteClusterHAAws)
	deleteClusterCmd.AddCommand(deleteClusterAws)

	deleteClusterAws.MarkFlagRequired("name")
	deleteClusterAws.MarkFlagRequired("region")
	deleteClusterAzure.MarkFlagRequired("name")
	deleteClusterAzure.MarkFlagRequired("region")
	deleteClusterCivo.MarkFlagRequired("name")
	deleteClusterCivo.MarkFlagRequired("region")
	deleteClusterHAAzure.MarkFlagRequired("name")
	deleteClusterHAAzure.MarkFlagRequired("region")
	deleteClusterHACivo.MarkFlagRequired("name")
	deleteClusterLocal.MarkFlagRequired("name")
	deleteClusterHAAws.MarkFlagRequired("name")
}
