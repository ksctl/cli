package cmd

// authors 	Dipankar Das <dipankardas0115@gmail.com>

import (
	"context"
	control_pkg "github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/helpers"
	"os"

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

var createClusterAws = &cobra.Command{
	Use:   "aws",
	Short: "Use to create a EKS cluster in AWS",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster aws <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var createClusterAzure = &cobra.Command{
	Use:   "azure",
	Short: "Use to create a AKS cluster in Azure",
	Long: `It is used to create cluster with the given name from user. For example:

	ksctl create-cluster azure <arguments to civo cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}

		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudAzure
		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)
		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterCivo = &cobra.Command{
	Use:   "civo",
	Short: "Use to create a CIVO k3s cluster",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster civo <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}

		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudCivo
		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)
		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterLocal = &cobra.Command{
	Use:   "local",
	Short: "Use to create a LOCAL cluster in Docker",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster local <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}

		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudLocal
		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)
		createManaged(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterHACivo = &cobra.Command{
	Use:   "ha-civo",
	Short: "Use to create a HA CIVO k3s cluster",
	Long: `It is used to create cluster with the given name from user. For example:

ksctl create-cluster ha-civo <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}
		SetRequiredFeatureFlags(cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo
		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)
		createHA(cmd.Flags().Lookup("approve").Changed)
	},
}

var createClusterHAAzure = &cobra.Command{
	Use:   "ha-azure",
	Short: "Use to create a HA k3s cluster in Azure",
	Long: `It is used to create cluster with the given name from user. For example:

	ksctl create-cluster ha-azure <arguments to civo cloud provider>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		cli.Client.Metadata.LogVerbosity = verbosity
		cli.Client.Metadata.LogWritter = os.Stdout

		if err := control_pkg.InitializeStorageFactory(context.WithValue(context.Background(), "USERID", helpers.GetUserName()), &cli.Client); err != nil {
			log.Error("Inialize Storage Driver", "Reason", err)
		}
		SetRequiredFeatureFlags(cmd)
		cli.Client.Metadata.Provider = consts.CloudAzure
		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)
		createHA(cmd.Flags().Lookup("approve").Changed)
	},
}

func init() {
	rootCmd.AddCommand(createClusterCmd)

	createClusterCmd.AddCommand(createClusterAws)
	createClusterCmd.AddCommand(createClusterAzure)
	createClusterCmd.AddCommand(createClusterCivo)
	createClusterCmd.AddCommand(createClusterLocal)
	createClusterCmd.AddCommand(createClusterHACivo)
	createClusterCmd.AddCommand(createClusterHAAzure)

	createClusterAzure.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("region")
	createClusterLocal.MarkFlagRequired("name")
	createClusterHAAzure.MarkFlagRequired("name")
	createClusterHACivo.MarkFlagRequired("name")
}
