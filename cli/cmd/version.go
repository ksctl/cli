package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	v0_1Ksctl = `
  _             _   _ 
 | |           | | | |
 | | _____  ___| |_| |
 | |/ / __|/ __| __| |
 |   <\__ \ (__| |_| |
 |_|\_\___/\___|\__|_|

`

	v2_0Ksctl = `
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

		color.HiGreen(v0_1Ksctl)

		x := strings.Split(v2_0Ksctl, "\n")

		y := []string{}

		colorCode := map[int]func(str string) string{
			0: func(str string) string { return color.New(color.BgHiMagenta).Add(color.FgHiBlack).SprintFunc()(str) },
			1: func(str string) string { return color.New(color.BgHiBlue).Add(color.FgHiBlack).SprintFunc()(str) },
			2: func(str string) string { return color.New(color.BgHiCyan).Add(color.FgHiBlack).SprintFunc()(str) },
			3: func(str string) string { return color.New(color.BgHiGreen).Add(color.FgHiBlack).SprintFunc()(str) },
			4: func(str string) string { return color.New(color.BgHiYellow).Add(color.FgHiBlack).SprintFunc()(str) },
			5: func(str string) string { return color.New(color.BgHiRed).Add(color.FgHiBlack).SprintFunc()(str) },
		}

		for i, _x := range x {
			fmt.Println(i, _x)
			if _y, ok := colorCode[i]; ok {
				y = append(y, _y(_x))
			} else {
				fmt.Println("Not found", i)
			}
		}
		fmt.Println(strings.Join(y, "\n"))

		fmt.Println("Version:", Version)
		fmt.Println("BuildDate:", BuildDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
