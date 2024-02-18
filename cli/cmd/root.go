/*
Kubesimplify
authors			Dipankar <dipankar@dipankar-das.com>
				Anurag Kumar <contact.anurag7@gmail.com>
				Avinesh Tripathi <avineshtripathi1@gmail.com>
*/

package cmd

import (
	"os"
	"time"

	controlPkg "github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/resources"
	"github.com/ksctl/ksctl/pkg/resources/controllers"

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
	//storage     string  // Currently only local storage is present
	distro string
	k8sVer string
	cloud  map[int]string
)

type CobraCmd struct {
	ClusterName string
	Region      string
	Client      resources.KsctlClient
	Version     string
}

var (
	cli        *CobraCmd
	controller controllers.Controller
	log        resources.LoggerFactory
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ksctl",
	Short: "CLI tool for managing multiple K8s clusters",
	Long: `CLI tool which can manage multiple K8s clusters
from local clusters to cloud provider specific clusters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	cli = new(CobraCmd)
	controller = controlPkg.GenKsctlController()
	log = logger.NewDefaultLogger(0, os.Stdout)
	log.SetPackageName("cli")

	cloud = map[int]string{
		1: string(consts.CloudAws),
		2: string(consts.CloudAzure),
		3: string(consts.CloudCivo),
		4: string(consts.CloudLocal),
	}
	cli.Client.Metadata.StateLocation = consts.StoreLocal

	timer := time.Now()
	err := rootCmd.Execute()
	defer log.Print("Time Took", "⏰", time.Since(timer).String())

	if err != nil {
		log.Error("Initialization of cli failed", "Reason", err)
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
