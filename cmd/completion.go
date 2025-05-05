package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (k *KsctlCommand) ShellCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion scripts",
		Long: `To load completions:

Bash:

  $ source <(ksctl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ ksctl completion bash > /etc/bash_completion.d/ksctl
  # macOS:
  $ ksctl completion bash > /usr/local/etc/bash_completion.d/ksctl

Zsh:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  $ ksctl completion zsh > "${fpath[1]}/_ksctl"

Fish:

  $ ksctl completion fish | source

  # To load completions for each session, execute once:
  $ ksctl completion fish > ~/.config/fish/completions/ksctl.fish
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}

	return cmd
}
