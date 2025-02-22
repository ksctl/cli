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
	"context"
	"os"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/config"
	cLogger "github.com/ksctl/cli/v2/pkg/logger"
	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/ksctl/ksctl/v2/pkg/provider"
	"github.com/ksctl/ksctl/v2/pkg/storage"
	"github.com/spf13/cobra"
)

type KsctlCommand struct {
	Ctx                     context.Context
	CliLog                  logger.Logger
	l                       logger.Logger
	ksctlStorage            storage.Storage
	root                    *cobra.Command
	verbose                 int
	debugMode               bool
	menuDriven              cli.MenuDriven
	KsctlConfig             *config.Config
	telemetry               *telemetry.Telemetry
	inMemInstanceTypesInReg map[string]provider.InstanceRegionOutput
}

func New() (*KsctlCommand, error) {
	k := new(KsctlCommand)
	k.KsctlConfig = new(config.Config)

	k.Ctx = context.WithValue(
		context.WithValue(
			context.Background(),
			consts.KsctlModuleNameKey,
			"cli",
		),
		consts.KsctlContextUserID,
		"cli",
	)

	k.root = k.NewRootCmd()

	k.CliLog = cLogger.NewLogger(0, os.Stdout)

	return k, nil
}

func (k *KsctlCommand) ForDocs() (*cobra.Command, error) {
	if err := k.CommandMapping(); err != nil {
		return nil, err
	}

	return k.root, nil
}

func (k *KsctlCommand) Execute() error {

	if err := config.LoadConfig(k.KsctlConfig); err != nil {
		return err
	}

	if err := k.CommandMapping(); err != nil {
		return err
	}

	return k.root.Execute()
}
