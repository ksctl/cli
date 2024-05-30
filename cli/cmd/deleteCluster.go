package cmd

// authors 	Dipankar Das <dipankardas0115@gmail.com>

import (
	"os"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"
	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:     "delete-cluster",
	Short:   "Use to delete a cluster",
	Aliases: []string{"delete"},
	Long: `It is used to delete cluster of given provider. For example:

ksctl delete-cluster ["aws","azure", "ha-<provider>", "civo", "local"]
`,
}

var deleteClusterAzure = &cobra.Command{
	Use:   "azure",
	Short: "Use to create a azure managed cluster",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster azure <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("approve").Changed)
	},
}

var deleteClusterCivo = &cobra.Command{
	Use:   "civo",
	Short: "Use to delete a CIVO cluster",
	Long: `It is used to delete cluster of given provider. For example:

ksctl delete-cluster civo
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("approve").Changed)

	},
}

var deleteClusterHAAws = &cobra.Command{
	Use:   "ha-aws",
	Short: "Use to delete a HA k3s cluster in Azure",
	Long: `It is used to delete cluster with the given name from user. For example:

	ksctl delete-cluster ha-aws <arguments to cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("approve").Changed)
	},
}

var deleteClusterHAAzure = &cobra.Command{
	Use:   "ha-azure",
	Short: "Use to delete a HA k3s cluster in Azure",
	Long: `It is used to delete cluster with the given name from user. For example:

	ksctl delete-cluster ha-azure <arguments to cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("approve").Changed)
	},
}

var deleteClusterHACivo = &cobra.Command{
	Use:   "ha-civo",
	Short: "Use to delete a HA CIVO k3s cluster",
	Long: `It is used to delete cluster with the given name from user. For example:

ksctl delete-cluster ha-civo <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		deleteHA(ctx, log, cmd.Flags().Lookup("approve").Changed)
	},
}

var deleteClusterLocal = &cobra.Command{
	Use:   "local",
	Short: "Use to delete a LOCAL cluster",
	Long: `It is used to delete cluster of given provider. For example:

ksctl delete-cluster local <arguments to local/Docker provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewGeneralLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudLocal

		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)

		deleteManaged(ctx, log, cmd.Flags().Lookup("approve").Changed)
	},
}

func init() {
	rootCmd.AddCommand(deleteClusterCmd)

	deleteClusterCmd.AddCommand(deleteClusterHACivo)
	deleteClusterCmd.AddCommand(deleteClusterCivo)
	deleteClusterCmd.AddCommand(deleteClusterHAAzure)
	deleteClusterCmd.AddCommand(deleteClusterAzure)
	deleteClusterCmd.AddCommand(deleteClusterLocal)
	deleteClusterCmd.AddCommand(deleteClusterHAAws)

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
