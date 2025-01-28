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
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) ScaleUp() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaleup",
		Example: `
ksctl update scaleup --help
		`,
		Short: "Use to manually scaleup a selfmanaged cluster",
		Long:  "It is used to manually scaleup a selfmanaged cluster",

		Run: func(cmd *cobra.Command, args []string) {
			l := k.l
			ctx := k.Ctx

			l.Box(ctx, "update", "scaleup cluster")
			l.Print(ctx, "info", "args", args)
		},
	}
	return cmd
}

func (k *KsctlCommand) ScaleDown() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaledown",
		Example: `
ksctl update scaledown --help
		`,
		Short: "Use to manually scaledown a selfmanaged cluster",
		Long:  "It is used to manually scaledown a selfmanaged cluster",

		Run: func(cmd *cobra.Command, args []string) {
			l := k.l
			ctx := k.Ctx

			l.Box(ctx, "update", "scaledown cluster")
			l.Print(ctx, "info", "args", args)
		},
	}
	return cmd
}
