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

import "github.com/spf13/cobra"

func (k *KsctlCommand) Addons() *cobra.Command {

	cmd := &cobra.Command{
		Use: "addons",
		Example: `
ksctl addons --help
`,
		Short: "Use to work with addons",
		Long:  "It is used to work with addons",
	}

	return cmd
}

func (k *KsctlCommand) ListAddon() *cobra.Command {

	cmd := &cobra.Command{
		Use: "list",
		Example: `
ksctl addons list --help
`,
		Short: "Use to list the addons",
		Long:  "It is used to list the addons",
	}

	return cmd
}

func (k *KsctlCommand) EnableAddon() *cobra.Command {

	cmd := &cobra.Command{
		Use: "enable",
		Example: `
ksctl addons enable --help
`,
		Short: "Use to enable an addon",
		Long:  "It is used to enable an addon",
	}
	return cmd
}

func (k *KsctlCommand) DisableAddon() *cobra.Command {

	cmd := &cobra.Command{
		Use: "disable",
		Example: `
ksctl addons disable --help
`,
		Short: "Use to disable an addon",
		Long:  "It is used to disable an addon",
	}
	return cmd
}
