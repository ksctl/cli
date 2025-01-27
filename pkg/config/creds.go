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

type MongoCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Port     *int   `json:"port,omitempty"`
	SRV      bool   `json:"srv,omitempty"`
}

type AwsCreds struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type AzureCreds struct {
	SubscriptionID string `json:"subscription_id"`
	TenantID       string `json:"tenant_id"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
}

type CloudCreds struct {
	Aws              *AwsCreds         `json:"aws,omitempty"`
	Azure            *AzureCreds       `json:"azure,omitempty"`
	CloudProviderSKU consts.KsctlCloud `json:"cloud_provider_sku"`
}

type ExternalStorageCreds struct {
	Mongo      *MongoCreds       `json:"mongo,omitempty"`
	StorageSKU consts.KsctlStore `json:"storage_sku"`
}

func locateCreds(s, prefix string) (fileLoc string, err error) {
	if len(s) == 0 {
		return "", fmt.Errorf("sku is empty")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "ksctl", "creds")
	configFile := filepath.Join(configDir, prefix+s+".json")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return configFile, fmt.Errorf("failed to create directory %s: %v", configDir, err)
		}
	}
	return configFile, nil
}

func SaveStorageCreds(c *ExternalStorageCreds, s consts.KsctlStore) error {
	credsFile, err := locateCreds(string(s), "s-")
	if err != nil {
		return err
	}

	file, err := os.Create(credsFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", credsFile, err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(c)
}

func LoadStorageCreds(c *ExternalStorageCreds, s consts.KsctlStore) (errC error) {
	credsFile, err := locateCreds(string(s), "s-")
	if err != nil {
		return err
	}

	file, err := os.Open(credsFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", credsFile, err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)
}

func SaveCloudCreds(c *CloudCreds, s consts.KsctlCloud) error {
	credsFile, err := locateCreds(string(s), "c-")
	if err != nil {
		return err
	}

	file, err := os.Create(credsFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", credsFile, err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(c)
}

func LoadCloudCreds(c *CloudCreds, s consts.KsctlCloud) (errC error) {
	credsFile, err := locateCreds(string(s), "c-")
	if err != nil {
		return err
	}

	file, err := os.Open(credsFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", credsFile, err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)
}
