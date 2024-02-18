package cmd

// authors 	Dipankar Das <dipankardas0115@gmail.com>

import (
	"context"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers"
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:     "delete-cluster",
	Short:   "Use to delete a cluster",
	Aliases: []string{"delete"},
	Long: `It is used to delete cluster of given provider. For example:

ksctl delete-cluster ["azure", "ha-<provider>", "civo", "local"]
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
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		deleteManaged(cmd.Flags().Lookup("approve").Changed)
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
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		deleteManaged(cmd.Flags().Lookup("approve").Changed)

	},
}

var deleteClusterHAAzure = &cobra.Command{
	Use:   "ha-azure",
	Short: "Use to delete a HA k3s cluster in Azure",
	Long: `It is used to delete cluster with the given name from user. For example:

	ksctl delete-cluster ha-azure <arguments to civo cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		deleteHA(cmd.Flags().Lookup("approve").Changed)
	},
}

var deleteClusterHACivo = &cobra.Command{
	Use:   "ha-civo",
	Short: "Use to delete a HA CIVO k3s cluster",
	Long: `It is used to delete cluster with the given name from user. For example:

ksctl delete-cluster ha-civo <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		deleteHA(cmd.Flags().Lookup("approve").Changed)
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
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudLocal

		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Inialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		deleteManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

func init() {
	rootCmd.AddCommand(deleteClusterCmd)

	deleteClusterCmd.AddCommand(deleteClusterHACivo)
	deleteClusterCmd.AddCommand(deleteClusterCivo)
	deleteClusterCmd.AddCommand(deleteClusterHAAzure)
	deleteClusterCmd.AddCommand(deleteClusterAzure)
	deleteClusterCmd.AddCommand(deleteClusterLocal)

	deleteClusterAzure.MarkFlagRequired("name")
	deleteClusterAzure.MarkFlagRequired("region")
	deleteClusterCivo.MarkFlagRequired("name")
	deleteClusterCivo.MarkFlagRequired("region")
	deleteClusterHAAzure.MarkFlagRequired("name")
	deleteClusterHAAzure.MarkFlagRequired("region")
	deleteClusterHACivo.MarkFlagRequired("name")
	deleteClusterLocal.MarkFlagRequired("name")
}
