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
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// change this using ldflags
var Version string = "dev"

var BuildDate string

func (k *KsctlCommand) Version() *cobra.Command {

	logoKsctl := `
░  ░░░░  ░░░      ░░░░      ░░░        ░░  ░░░░░░░
▒  ▒▒▒  ▒▒▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒▒▒
▓     ▓▓▓▓▓▓      ▓▓▓  ▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓
▓  ▓▓▓  ▓▓▓▓▓▓▓▓▓  ▓▓  ▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓
█  ████  ███      ████      ██████  █████        █
`

	cmd := &cobra.Command{
		Use: "version",
		Example: `
ksctl version --help
		`,
		Short: "ksctl version",
		Long:  "To get version for ksctl components",
		Run: func(cmd *cobra.Command, args []string) {
			// color.New(color.BgHiGreen).Add(color.FgHiBlack).Println(logoKsctl)
			for _, line := range strings.Split(logoKsctl, "\n") {
				color.New(color.FgHiGreen).Add(color.BgBlack).Println(line)
			}
			k.l.Note(k.Ctx, "Components", color.HiGreenString("ksctl:cli"), color.HiBlueString(Version), color.HiGreenString("ksctl:core"), color.HiBlueString("v2"))
			k.l.Note(k.Ctx, "Build Information", "date", BuildDate)
		},
	}

	return cmd
}
