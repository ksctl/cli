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
	"encoding/json"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/config"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/statefile"
	"github.com/ksctl/ksctl/v2/pkg/utilities"
	"github.com/spf13/cobra"
)

func (k *KsctlCommand) Configure() *cobra.Command {

	cmd := &cobra.Command{
		Use: "configure",

		Short: "Configure ksctl cli",
		Long:  "It will display the current ksctl cli configuration",
		Run: func(cmd *cobra.Command, args []string) {
			headers := []string{"Property", "Value"}

			enabled := color.HiCyanString("✔")
			disabled := color.HiRedString("✘")
			telemetry := enabled

			if k.KsctlConfig.Telemetry != nil && !*k.KsctlConfig.Telemetry {
				telemetry = disabled
			}
			rows := [][]string{
				{"PreferedStorage", string(k.KsctlConfig.PreferedStateStore)},
				{"TelemetryStatus", telemetry},
			}

			if k.KsctlConfig.PreferedStateStore == consts.StoreExtMongo {
				if err := k.loadMongoCredentials(); err != nil {
					rows = append(rows, []string{"MongoDB 💾", disabled})
				} else {
					rows = append(rows, []string{"MongoDB 💾", enabled})
				}
			}

			if _, err := k.loadAwsCredentials(); err == nil {
				rows = append(rows, []string{"AWS ☁️", enabled})
			} else {
				rows = append(rows, []string{"AWS ☁️", disabled})
			}

			if _, err := k.loadAzureCredentials(); err == nil {
				rows = append(rows, []string{"Azure ☁️", enabled})
			} else {
				rows = append(rows, []string{"Azure ☁️", disabled})
			}

			k.l.Table(k.Ctx, headers, rows)
		},
	}

	return cmd
}

func (k *KsctlCommand) ConfigureStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use: "storage",

		Short: "Configure storage",
		Long:  "It will help you to configure the storage",
		Run: func(cmd *cobra.Command, args []string) {
			if ok := k.handleStorageConfig(); !ok {
				os.Exit(1)
			}
		},
	}

	return cmd
}

func (k *KsctlCommand) ConfigureCloud() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cloud",

		Short: "Configure cloud",
		Long:  "It will help you to configure the cloud",
		Run: func(cmd *cobra.Command, args []string) {
			if ok := k.handleCloudConfig(); !ok {
				os.Exit(1)
			}
		},
	}

	return cmd
}

func (k *KsctlCommand) ConfigureTelemetry() *cobra.Command {
	cmd := &cobra.Command{
		Use: "telemetry",

		Short: "Configure telemetry",
		Long:  "It will help you to configure the telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := k.menuDriven.Confirmation("Do you want to enable the telemetry?", cli.WithDefaultValue("yes")); err != nil {
				k.l.Error("Failed to get the telemetry status", "Reason", err)
				os.Exit(1)
			} else {
				k.KsctlConfig.Telemetry = utilities.Ptr(v)
				if err := config.SaveConfig(k.KsctlConfig); err != nil {
					k.l.Error("Failed to save the configuration", "Reason", err)
					os.Exit(1)
				}
			}
		},
	}

	return cmd
}

func (k *KsctlCommand) handleStorageConfig() bool {
	if v, err := k.menuDriven.DropDown(
		"What should be your default storageDriver?",
		map[string]string{
			"MongoDb": string(consts.StoreExtMongo),
			"Local":   string(consts.StoreLocal),
		},
		cli.WithDefaultValue("Local"),
	); err != nil {
		k.l.Error("Failed to get the storageDriver", "Reason", err)
		return false
	} else {
		k.KsctlConfig.PreferedStateStore = consts.KsctlStore(v)
		errL := config.SaveConfig(k.KsctlConfig)
		if errL != nil {
			k.l.Error("Failed to save the configuration", "Reason", errL)
			return false
		}

		if consts.KsctlStore(v) == consts.StoreExtMongo {
			k.l.Note(k.Ctx, "You need to provide the credentials for the MongoDB")
			if err := k.storeMongoCredentials(); err != nil {
				k.l.Error("Failed to store the MongoDB credentials", "Reason", err)
				return false
			}
		}
	}
	return true
}

func (k *KsctlCommand) handleCloudConfig() bool {
	if v, err := k.menuDriven.DropDown(
		"Credentials",
		map[string]string{
			"Amazon Web Services": string(consts.CloudAws),
			"Azure":               string(consts.CloudAzure),
		},
	); err != nil {
		k.l.Error("Failed to get the credentials", "Reason", err)
		return false
	} else {
		switch consts.KsctlCloud(v) {
		case consts.CloudAws:
			if err := k.storeAwsCredentials(); err != nil {
				k.l.Error("Failed to store the AWS credentials", "Reason", err)
				return false
			}
		case consts.CloudAzure:
			if err := k.storeAzureCredentials(); err != nil {
				k.l.Error("Failed to store the Azure credentials", "Reason", err)
				return false
			}
		}
	}

	return true
}

func (k *KsctlCommand) storeAwsCredentials() (err error) {
	c := new(statefile.CredentialsAws)
	c.AccessKeyId, err = k.menuDriven.TextInputPassword("Enter your AWS Access Key ID")
	if err != nil {
		return err
	}
	c.SecretAccessKey, err = k.menuDriven.TextInputPassword("Enter your AWS Secret Access Key")
	if err != nil {
		return err
	}

	return config.SaveCloudCreds(c, consts.CloudAws)
}

func (k *KsctlCommand) loadAwsCredentials() ([]byte, error) {
	c := new(statefile.CredentialsAws)
	if err := config.LoadCloudCreds(c, consts.CloudAws); err != nil {
		return nil, err
	}
	v, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (k *KsctlCommand) storeAzureCredentials() (err error) {
	c := new(statefile.CredentialsAzure)
	c.SubscriptionID, err = k.menuDriven.TextInputPassword("Enter your Azure Subscription ID")
	if err != nil {
		return err
	}

	c.TenantID, err = k.menuDriven.TextInputPassword("Enter your Azure Tenant ID")
	if err != nil {
		return err
	}

	c.ClientID, err = k.menuDriven.TextInputPassword("Enter your Azure Client ID")
	if err != nil {
		return err
	}
	c.ClientSecret, err = k.menuDriven.TextInputPassword("Enter your Azure Client Secret")
	if err != nil {
		return err
	}

	return config.SaveCloudCreds(c, consts.CloudAzure)
}

func (k *KsctlCommand) loadAzureCredentials() ([]byte, error) {
	c := new(statefile.CredentialsAzure)
	if err := config.LoadCloudCreds(c, consts.CloudAzure); err != nil {
		return nil, err
	}
	v, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (k *KsctlCommand) storeMongoCredentials() (err error) {
	c := new(statefile.CredentialsMongodb)
	srv, err := k.menuDriven.Confirmation("Enter whether MongoDB has SRV record or not", cli.WithDefaultValue("no"))
	if err != nil {
		return err
	}
	c.SRV = srv

	c.Domain, err = k.menuDriven.TextInput("Enter your MongoDB URI")
	if err != nil {
		return err
	}
	c.Username, err = k.menuDriven.TextInputPassword("Enter your MongoDB Username")
	if err != nil {
		return err
	}
	c.Password, err = k.menuDriven.TextInputPassword("Enter your MongoDB Password")
	if err != nil {
		return err
	}
	port := ""
	if port, err = k.menuDriven.TextInput("Enter your MongoDB Port"); err != nil {
		return err
	}
	if len(port) != 0 {
		v, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		c.Port = utilities.Ptr(v)
	}

	return config.SaveStorageCreds(c, consts.StoreExtMongo)
}

func (k *KsctlCommand) loadMongoCredentials() error {
	c := new(statefile.CredentialsMongodb)
	if err := config.LoadStorageCreds(c, consts.StoreExtMongo); err != nil {
		return err
	}
	v, err := json.Marshal(c)
	if err != nil {
		return err
	}

	k.Ctx = context.WithValue(k.Ctx, consts.KsctlMongodbCredentials, v)
	return nil
}
