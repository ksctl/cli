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
	"os"

	"github.com/ksctl/cli/pkg/cli"
	cLogger "github.com/ksctl/cli/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/storage"
	"github.com/spf13/cobra"
)

type KsctlCommand struct {
	Log          logger.Logger
	ksctlStorage storage.Storage
	root         *cobra.Command
	verbose      int
}

func New() (*KsctlCommand, error) {
	k := new(KsctlCommand)

	k.root = NewRootCmd()

	cli.AddVerboseFlag(k.root, &k.verbose)
	k.Log = cLogger.NewLogger(k.verbose, os.Stdout)

	return k, nil
}

func (k *KsctlCommand) Execute() error {
	return k.root.Execute()
}
