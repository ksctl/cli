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

// for the newController we should be able to pass some option fields for control more things
// for example whther it is a dry-run for testing

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "ksctl",
		Short: "CLI tool for managing multiple K8s clusters",
		Long:  LongMessage("CLI tool which can manage multiple K8s clusters from local clusters to cloud provider specific clusters."),
	}

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}
