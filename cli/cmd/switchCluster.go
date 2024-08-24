package cmd

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/controllers"
	"github.com/ksctl/ksctl/pkg/types"

	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/spf13/cobra"
)

var switchCluster = &cobra.Command{
	Use: "switch-cluster",
	Example: `
ksctl switch-context --provider civo --name <clustername> --region <region>
ksctl switch-context --provider local --name <clustername>
ksctl switch-context --provider azure --name <clustername> --region <region>
ksctl switch-context --provider ha-civo --name <clustername> --region <region>
ksctl switch-context --provider ha-azure --name <clustername> --region <region>
ksctl switch-context --provider ha-aws --name <clustername> --region <region>
ksctl switch-context --provider aws --name <clustername> --region <region>

	For Storage specific

ksctl switch-context -s store-local -p civo -n <clustername> -r <region>
ksctl switch-context -s external-store-mongodb -p civo -n <clustername> -r <region>
`,
	Aliases: []string{"switch", "access"},
	Short:   "Use to switch between clusters",
	Long:    LongMessage("It is used to switch cluster with the given ClusterName from user."),
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		if len(storage) == 0 {
			storage = string(consts.StoreLocal)
		}
		SetRequiredFeatureFlags(ctx, log, cmd)

		cli.Client.Metadata.ClusterName = clusterName
		cli.Client.Metadata.Region = region
		cli.Client.Metadata.StateLocation = consts.KsctlStore(storage)

		switch provider {
		case string(consts.CloudLocal):
			cli.Client.Metadata.Provider = consts.CloudLocal

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudCivo):
			cli.Client.Metadata.Provider = consts.CloudCivo
			cli.Client.Metadata.IsHA = true

		case string(consts.CloudCivo):
			cli.Client.Metadata.Provider = consts.CloudCivo

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudAzure):
			cli.Client.Metadata.Provider = consts.CloudAzure
			cli.Client.Metadata.IsHA = true

		case string(consts.ClusterTypeHa) + "-" + string(consts.CloudAws):
			cli.Client.Metadata.Provider = consts.CloudAws
			cli.Client.Metadata.IsHA = true

		case string(consts.CloudAws):
			cli.Client.Metadata.Provider = consts.CloudAws

		case string(consts.CloudAzure):
			cli.Client.Metadata.Provider = consts.CloudAzure
		}

		m, err := controllers.NewManagerClusterKsctl(
			ctx,
			log,
			&cli.Client,
		)
		if err != nil {
			log.Error("failed to init", "Reason", err)
			os.Exit(1)
		}
		kubeconfig, err := m.SwitchCluster()
		if err != nil {
			log.Error("Switch cluster failed", "Reason", err)
			os.Exit(1)
		}
		log.Debug(ctx, "kubeconfig output as string", "kubeconfig", kubeconfig)
		log.Success(ctx, "Switch cluster Successful")

		if accessMode == "k9s" {
			K9sAccess(log)
		} else if accessMode == "shell" {
			shellAccess(log)
		} else {
			log.Print(ctx, "No mode selected")
		}
	},
}

func shellAccess(log types.LoggerFactory) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home dir", "Reason", err)
		os.Exit(1)
	}

	home = filepath.Join(home, ".ksctl", "kubeconfig")
	cmd := exec.Command("/bin/bash")

	cmd.Env = append(os.Environ(), "KUBECONFIG="+home)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Error creating pseudo-terminal:", err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Println("Error resizing pty:", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	// Print welcome message
	fmt.Fprintln(ptmx, "echo Hi from Ksctl team! You are now in the shell session having cluster context.")
	fmt.Fprintln(ptmx, "kubectl get nodes -owide")

	// Copy stdin to ptmx, and ptmx to stdout
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}

func K9sAccess(log types.LoggerFactory) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home dir", "Reason", err)
		os.Exit(1)
	}
	home = filepath.Join(home, ".ksctl", "kubeconfig")
	_cmd := exec.Command("k9s", "--kubeconfig", home)

	_bout := new(strings.Builder)
	_berr := new(strings.Builder)
	_cmd.Stdout = _bout
	_cmd.Stderr = _berr

	if err := _cmd.Run(); err != nil {
		log.Error("Failed to run k9s", "Reason", err)
	}
	_stdout, _stderr := _bout.String(), _berr.String()
	fmt.Println(color.HiBlueString(_stdout))
	fmt.Println(color.HiRedString(_stderr))
}

func init() {
	RootCmd.AddCommand(switchCluster)
	clusterNameFlag(switchCluster)
	regionFlag(switchCluster)
	storageFlag(switchCluster)

	switchCluster.Flags().StringVarP(&provider, "provider", "p", "", "Provider")
	switchCluster.Flags().StringVarP(&accessMode, "mode", "m", "", "Mode of access can be shell or k9s or none")

	switchCluster.MarkFlagRequired("name")
	switchCluster.MarkFlagRequired("provider")
}
