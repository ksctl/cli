package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	logoKsctl = `
░  ░░░░  ░░░      ░░░░      ░░░        ░░  ░░░░░░░
▒  ▒▒▒  ▒▒▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒▒▒
▓     ▓▓▓▓▓▓      ▓▓▓  ▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓
▓  ▓▓▓  ▓▓▓▓▓▓▓▓▓  ▓▓  ▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓
█  ████  ███      ████      ██████  █████        █
`
)

// change this using ldflags
var Version string = "dev"

var BuildDate string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ksctl",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(newLogo())

		fmt.Println("Version:", Version)
		fmt.Println("BuildDate:", BuildDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
