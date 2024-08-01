package cmd

import (
	"os"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/types"
	"github.com/spf13/cobra"
)

var createClusterCmd = &cobra.Command{
	Use: "create-cluster",
	Example: `
ksctl create --help
	`,
	Short:   "Use to create a cluster",
	Aliases: []string{"create"},
	Long:    "It is used to create cluster with the given name from user",
}

var createClusterAzure = &cobra.Command{
	Use: "azure",
	Example: `
ksctl create-cluster azure -n demo -r eastus -s store-local --nodeSizeMP Standard_DS2_v2 --noMP 3
`,
	Short: "Use to create a AKS cluster in Azure",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeMang)

		createManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterAws = &cobra.Command{
	Use: "aws",
	Example: `
ksctl create-cluster aws -n demo -r ap-south-1 -s store-local --nodeSizeMP t2.micro --noMP 3
`,
	Short: "Use to create a EKS cluster in Aws",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeMang)

		createManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterCivo = &cobra.Command{
	Use: "civo",
	Example: `
ksctl create-cluster civo --name demo --region LON1 --storage store-local --nodeSizeMP g4s.kube.small --noMP 3
`,
	Short: "Use to create a Civo managed k3s cluster",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeMang)

		createManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterLocal = &cobra.Command{
	Use: "local",
	Example: `
ksctl create-cluster local --name demo --storage store-local --noMP 3
`,
	Short: "Use to create a kind cluster",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudLocal

		SetDefaults(consts.CloudLocal, consts.ClusterTypeMang)

		createManaged(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterHAAws = &cobra.Command{
	Use: "ha-aws",
	Example: `
ksctl create-cluster ha-aws -n demo -r us-east-1 --bootstrap k3s -s store-local --nodeSizeCP t2.medium --nodeSizeWP t2.medium --nodeSizeLB t2.micro --nodeSizeDS t2.small --noWP 1 --noCP 3 --noDS 3
`,
	Short: "Use to create a self-managed Highly Available cluster on AWS",
	Long:  "It is used to create cluster with the given name from user.",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAws

		SetDefaults(consts.CloudAws, consts.ClusterTypeHa)

		createHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterHACivo = &cobra.Command{
	Use: "ha-civo",
	Example: `
ksctl create-cluster ha-civo --name demo --region LON1 --bootstrap k3s --storage store-local --nodeSizeCP g3.small --nodeSizeWP g3.medium --nodeSizeLB g3.small --nodeSizeDS g3.small --noWP 1 --noCP 3 --noDS 3
ksctl create-cluster ha-civo --name demo --region LON1 --bootstrap kubeadm --storage store-local --nodeSizeCP g3.medium --nodeSizeWP g3.large --nodeSizeLB g3.small --nodeSizeDS g3.small --noWP 1 --noCP 3 --noDS 3
`,
	Short: "Use to create a self-managed Highly Available cluster on Civo",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudCivo

		SetDefaults(consts.CloudCivo, consts.ClusterTypeHa)

		createHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

var createClusterHAAzure = &cobra.Command{
	Use: "ha-azure",
	Example: `
ksctl create-cluster ha-azure --name demo --region eastus --bootstrap k3s --storage store-local --nodeSizeCP Standard_F2s --nodeSizeWP Standard_F2s --nodeSizeLB Standard_F2s --nodeSizeDS Standard_F2s --noWP 1 --noCP 3 --noDS 3
ksctl create-cluster ha-azure --name demo --region eastus --bootstrap kubeadm --storage store-local --nodeSizeCP Standard_F2s --nodeSizeWP Standard_F4s --nodeSizeLB Standard_F2s --nodeSizeDS Standard_F2s --noWP 1 --noCP 3 --noDS 3
`,
	Short: "Use to create a self-managed Highly-Available cluster on Azure",
	Long:  "It is used to create cluster with the given name from user",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.Provider = consts.CloudAzure

		SetDefaults(consts.CloudAzure, consts.ClusterTypeHa)

		createHA(ctx, log, cmd.Flags().Lookup("yes").Changed)
	},
}

func init() {
	RootCmd.AddCommand(createClusterCmd)

	createClusterCmd.AddCommand(createClusterAzure)
	createClusterCmd.AddCommand(createClusterCivo)
	createClusterCmd.AddCommand(createClusterLocal)
	createClusterCmd.AddCommand(createClusterHACivo)
	createClusterCmd.AddCommand(createClusterHAAzure)
	createClusterCmd.AddCommand(createClusterHAAws)
	createClusterCmd.AddCommand(createClusterAws)

	createClusterAzure.MarkFlagRequired("name")
	createClusterAws.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("region")
	createClusterLocal.MarkFlagRequired("name")
	createClusterHAAzure.MarkFlagRequired("name")
	createClusterHACivo.MarkFlagRequired("name")
	createClusterHAAws.MarkFlagRequired("name")
}
