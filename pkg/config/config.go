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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ksctl/ksctl/v2/pkg/consts"
)

// ksctl config in ~/.config/ksctl/config.json (handled by ksctl:cli)
// ksctl credentials in ~/.config/ksctl/creds/(aws|azure|mongodb).json (handled by ksctl:cli)
// ksctl state in ~/.ksctl/state/..... (handled by the ksctl:core:storage)

// NOTE
// store all the credentials in the local file system and for state we decide on the external storage
// we will retrieve the credentials and load them in the context.Background() and pass it to the respective clients
// So we need to initialize the HostStorage() to get the credentials from the local file system once thats done
// we can populate the context.Background() with the credentials. and initialze the client which that preferedstateStore unless specified in the command argument

type Config struct {
	PreferedStateStore consts.KsctlStore `json:"preferedStateStore"`
	DefaultProvider    consts.KsctlCloud `json:"defaultCloud"`
}

func LoadConfig(c *Config) (errC error) {

	configFile, err := locateConfig()
	if err != nil {
		return err
	}

	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", configFile, err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)
}

func locateConfig() (fileLoc string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "ksctl")
	configFile := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return configFile, fmt.Errorf("failed to create directory %s: %v", configDir, err)
		}
	}
	return configFile, nil
}

func SaveConfig(c *Config) error {
	configFile, err := locateConfig()
	if err != nil {
		return err
	}

	file, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", configFile, err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(c)
}
