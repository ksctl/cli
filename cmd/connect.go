// Copyright 2025 Ksctl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/common"
	"github.com/ksctl/ksctl/v2/pkg/handler/cluster/controller"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func (k *KsctlCommand) Connect() *cobra.Command {

	cmd := &cobra.Command{
		Use: "connect",
		Example: `
ksctl connect --help
		`,
		Short: "Connect to existing cluster",
		Long:  "It is used to connect to existing cluster",

		Run: func(cmd *cobra.Command, args []string) {
			clusters, err := k.fetchAllClusters()
			if err != nil {
				k.l.Error("Error in fetching the clusters", "Error", err)
				os.Exit(1)
			}

			if len(clusters) == 0 {
				k.l.Error("No clusters found to connect")
				os.Exit(1)
			}

			selectDisplay := make(map[string]string, len(clusters))
			valueMaping := make(map[string]controller.Metadata, len(clusters))

			for idx, cluster := range clusters {
				selectDisplay[makeHumanReadableList(cluster)] = strconv.Itoa(idx)
				valueMaping[strconv.Itoa(idx)] = controller.Metadata{
					ClusterName:   cluster.Name,
					ClusterType:   cluster.ClusterType,
					Provider:      cluster.CloudProvider,
					Region:        cluster.Region,
					StateLocation: k.KsctlConfig.PreferedStateStore,
					K8sDistro:     cluster.K8sDistro,
				}
			}

			selectedCluster, err := k.menuDriven.DropDown(
				"Select the cluster to delete",
				selectDisplay,
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			m := valueMaping[selectedCluster]

			if k.loadCloudProviderCreds(m.Provider) != nil {
				os.Exit(1)
			}

			c, err := common.NewController(
				k.Ctx,
				k.l,
				&controller.Client{
					Metadata: m,
				},
			)
			if err != nil {
				k.l.Error("Failed to create the controller", "Reason", err)
				os.Exit(1)
			}

			kubeconfig, err := c.Switch()
			if err != nil {
				k.l.Error("Failed to connect to the cluster", "Reason", err)
				os.Exit(1)
			}

			k.l.Note(k.Ctx, "Downloaded the kubeconfig")

			k.writeKubeconfig([]byte(*kubeconfig))

			accessMode, err := k.menuDriven.DropDown(
				"Select the access mode",
				map[string]string{
					"k9s":  "k9s",
					"bash": "shell",
					"none": "none",
				},
				cli.WithDefaultValue("none"),
			)
			if err != nil {
				k.l.Error("Failed to get userinput", "Reason", err)
				os.Exit(1)
			}

			if accessMode == "k9s" {
				K9sAccess(k.l)
			} else if accessMode == "shell" {
				shellAccess(k.l)
			} else {
				k.l.Box(k.Ctx, "Kubeconfig", "You can access the cluster using $ kubectl commands or any other k8s client as its saved to ~/.kube/config")
			}
		},
	}

	return cmd
}

func shellAccess(log logger.Logger) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home dir", "Reason", err)
		os.Exit(1)
	}

	home = filepath.Join(home, ".kube", "config")
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

	fmt.Fprintln(ptmx, "echo Hi from Ksctl team! You are now in the shell session having cluster context.")
	fmt.Fprintln(ptmx, "kubectl get nodes -owide && kubectl cluster-info")

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}

func K9sAccess(log logger.Logger) {
	// home = filepath.Join(home, ".ksctl", "kubeconfig")
	// _cmd := exec.Command("k9s", "--kubeconfig", home)
	_cmd := exec.Command("k9s")

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

func mergeKubeConfigs(configs ...*clientcmdapi.Config) *clientcmdapi.Config {
	merged := clientcmdapi.NewConfig()
	for _, cfg := range configs {
		for name, cluster := range cfg.Clusters {
			merged.Clusters[name] = cluster
		}
		for name, authInfo := range cfg.AuthInfos {
			merged.AuthInfos[name] = authInfo
		}
		for name, context := range cfg.Contexts {
			merged.Contexts[name] = context
		}
		if cfg.CurrentContext != "" {
			merged.CurrentContext = cfg.CurrentContext
		}
	}
	return merged
}

func (k *KsctlCommand) writeKubeconfig(newKubeconfig []byte) {
	home, err := os.UserHomeDir()
	if err != nil {
		k.l.Error("Failed to get the home directory", "Reason", err)
		os.Exit(1)
	}
	orgConfig, err := os.ReadFile(filepath.Join(home, ".kube", "config"))
	if err != nil {
		k.l.Error("Failed to read the kubeconfig", "Reason", err)
		os.Exit(1)
	}

	config1, err := clientcmd.Load(orgConfig)
	if err != nil {
		k.l.Error("Failed to load the kubeconfig in ~/.kube/config", "Reason", err)
		os.Exit(1)
	}
	config2, err := clientcmd.Load(newKubeconfig)
	if err != nil {
		k.l.Error("Failed to load the new kubeconfig", "Reason", err)
		os.Exit(1)
	}

	mergedConfig := mergeKubeConfigs(config1, config2)

	mergedConfig.CurrentContext = config2.CurrentContext

	mergedYAML, err := clientcmd.Write(*mergedConfig)
	if err != nil {
		k.l.Error("Failed to write the merged kubeconfig", "Reason", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(home, ".kube", "config"), mergedYAML, 0640); err != nil {
		k.l.Error("Failed to write the kubeconfig", "Reason", err)
		os.Exit(1)
	}
}
