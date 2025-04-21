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
	"time"
)

type UpdateCache struct {
	LastChecked         time.Time     `json:"lastChecked"`
	AvailableVersions   bool          `json:"availableVersions"`
	UpdateCheckInterval time.Duration `json:"updateCheckInterval"`
}

var DefaultUpdateCache = &UpdateCache{
	UpdateCheckInterval: time.Minute,
}

func locateUpdateCacheFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "ksctl")
	configFile := filepath.Join(configDir, "cache-update.json")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return configFile, fmt.Errorf("failed to create directory %s: %v", configDir, err)
		}
	}
	return configFile, nil
}

func LoadUpdateCache(c *UpdateCache) (errC error) {

	configFile, err := locateUpdateCacheFile()
	if err != nil {
		return err
	}

	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// NOTE: writing default config
			c = DefaultUpdateCache
			return SaveUpdateCache(c)
		}
		return fmt.Errorf("failed to open file %s: %v", configFile, err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)
}

func SaveUpdateCache(c *UpdateCache) error {
	configFile, err := locateUpdateCacheFile()
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
