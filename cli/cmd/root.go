/*
Kubesimplify
authors			Dipankar <dipankar@dipankar-das.com>
				Anurag Kumar <contact.anurag7@gmail.com>
				Avinesh Tripathi <avineshtripathi1@gmail.com>
*/

package cmd

import (
	"context"
	"os"
	"time"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/spf13/cobra"
)

var (
	clusterName string
	region      string
	noCP        int
	noWP        int
	noMP        int
	noDS        int
	nodeSizeMP  string
	nodeSizeCP  string
	nodeSizeWP  string
	nodeSizeLB  string
	nodeSizeDS  string
	apps        string
	cni         string
	provider    string
	storage     string
	distro      string
	k8sVer      string
	cloud       map[int]string
)

type CobraCmd struct {
	ClusterName string
	Region      string
	Client      types.KsctlClient
	Version     string
}

var (
	cli    *CobraCmd
	logCli types.LoggerFactory
	ctx    context.Context
)

var rootCmd = &cobra.Command{
	Use:   "ksctl",
	Short: "CLI tool for managing multiple K8s clusters",
	Long: `
  _             _   _ 
 | |           | | | |
 | | _____  ___| |_| |
 | |/ / __|/ __| __| |
 |   <\__ \ (__| |_| |
 |_|\_\___/\___|\__|_|

	CLI tool which can manage multiple K8s clusters
from local clusters to cloud provider specific clusters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {

	ctx = context.WithValue(
		context.Background(),
		consts.KsctlModuleNameKey,
		"cli",
	)
	ctx = context.WithValue(
		ctx, "USERID", "cli",
	)
	if _, ok := os.LookupEnv("KSCTL_FAKE_FLAG_ENABLED"); ok {
		ctx = context.WithValue(
			ctx,
			consts.KsctlTestFlagKey,
			"true",
		)
	}

	cli = new(CobraCmd)
	logCli = logger.NewGeneralLogger(0, os.Stdout)

	cloud = map[int]string{
		1: string(consts.CloudAws),
		2: string(consts.CloudAzure),
		3: string(consts.CloudCivo),
		4: string(consts.CloudLocal),
	}

	timer := time.Now()
	err := rootCmd.Execute()
	defer logCli.Print(ctx, "Time Took", "time", time.Since(timer).String())

	if err != nil {
		logCli.Error(ctx, "Initialization of cli failed", "Reason", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Kubesimpctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	verboseFlags()

	argsFlags()
}
