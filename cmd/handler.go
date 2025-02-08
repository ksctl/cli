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

import "github.com/ksctl/cli/pkg/cli"

func (k *KsctlCommand) CommandMapping() error {
	c := k.Cluster()
	cr := k.Configure()
	cl := k.List()

	cli.RegisterCommand(
		k.root,
		c,
		k.Version(),
		cr,
	)
	cli.RegisterCommand(
		c,
		k.Create(),
		k.Delete(),
		cl,
		k.Connect(),
		k.ScaleUp(),
		k.ScaleDown(),
	)

	cli.RegisterCommand(
		cl,
		k.ListAll(),
	)

	cli.RegisterCommand(
		cr,
		k.ConfigureStorage(),
		k.ConfigureCloud(),
	)

	return nil
}
