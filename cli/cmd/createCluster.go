package cmd

// authors 	Dipankar Das <dipankardas0115@gmail.com>

import (
	"context"
	"os"

	"github.com/ksctl/ksctl/pkg/helpers"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

// createClusterCmd represents the createCluster command
var createClusterCmd = &cobra.Command{
	Use:     "create-cluster",
	Short:   "Use to create a cluster",
	Aliases: []string{"create"},
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster ["azure", "gcp", "aws", "local"]
`,
}

var createClusterHAAws = &cobra.Command{
	Use:   "ha-aws",
	Short: "Use to create a EKS cluster in AWS",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster ha-aws <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeHa)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Initialize Storage Driver", "Reason", err)
			os.Exit(1)
		}
		createHA(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterAzure = &cobra.Command{
	Use:   "azure",
	Short: "Use to create a AKS cluster in Azure",
	Long: `It is used to create cluster with the given name from user. For example:

	ksctl create-cluster azure <arguments to cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Initialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterCivo = &cobra.Command{
	Use:   "civo",
	Short: "Use to create a CIVO k3s cluster",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster civo <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Initialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterLocal = &cobra.Command{
	Use:   "local",
	Short: "Use to create a LOCAL cluster in Docker",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster local <arguments to cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudLocal

		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Initialize Storage Driver", "Reason", err)
			os.Exit(1)
		}

		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterHACivo = &cobra.Command{
	Use:   "ha-civo",
	Short: "Use to create a HA CIVO k3s cluster",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster ha-civo <arguments to cloud provider>
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
		createHA(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterHAAzure = &cobra.Command{
	Use:   "ha-azure",
	Short: "Use to create a HA k3s cluster in Azure",
	Long: `It is used to create cluster with the given name from user. For example:

	ksctl create-cluster ha-azure <arguments to cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout
		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		if err := safeInitializeStorageLoggerFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName())); err != nil {
			log.Error("Failed Initialize Storage Driver", "Reason", err)
			os.Exit(1)
		}
		createHA(cmd.Flags().Lookup("approve").Changed)
	},
}

func init() {
	rootCmd.AddCommand(createClusterCmd)

	createClusterCmd.AddCommand(createClusterAzure)
	createClusterCmd.AddCommand(createClusterCivo)
	createClusterCmd.AddCommand(createClusterLocal)
	createClusterCmd.AddCommand(createClusterHACivo)
	createClusterCmd.AddCommand(createClusterHAAzure)
	createClusterCmd.AddCommand(createClusterHAAws)

	createClusterAzure.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("region")
	createClusterLocal.MarkFlagRequired("name")
	createClusterHAAzure.MarkFlagRequired("name")
	createClusterHACivo.MarkFlagRequired("name")
	createClusterHAAws.MarkFlagRequired("name")
}
