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

package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/ksctl/cli/v2/pkg/config"
	"github.com/ksctl/ksctl/v2/pkg/addons"
	"github.com/ksctl/ksctl/v2/pkg/consts"
	"github.com/ksctl/ksctl/v2/pkg/logger"
)

var (
	clientIdentity string
)

type TelemetryEvent string

const (
	EventClusterCreate       TelemetryEvent = "cluster_create"
	EventClusterDelete       TelemetryEvent = "cluster_delete"
	EventClusterConnect      TelemetryEvent = "cluster_connect"
	EventClusterList         TelemetryEvent = "cluster_list"
	EventClusterGet          TelemetryEvent = "cluster_get"
	EventClusterScaleDown    TelemetryEvent = "cluster_scaledown"
	EventClusterScaleUp      TelemetryEvent = "cluster_scaleup"
	EventClusterUpgrade      TelemetryEvent = "cli_upgrade"
	EventClusterAddonEnable  TelemetryEvent = "cluster_addon_enable"
	EventClusterAddonDisable TelemetryEvent = "cluster_addon_disable"
)

type TelemetryAddon struct {
	Sku     string `json:"sku"`
	Version string `json:"version"`
	Label   string `json:"label"`
}

func TranslateMetadata(addons addons.ClusterAddons) []TelemetryAddon {
	var telemetryAddons []TelemetryAddon

	for _, addon := range addons {
		telemetryAddons = append(telemetryAddons, TelemetryAddon{
			Sku:   addon.Name,
			Label: addon.Label,
		})
	}

	return telemetryAddons
}

type TelemetryMeta struct {
	CloudProvider     consts.KsctlCloud       `json:"cloud_provider"`
	StorageDriver     consts.KsctlStore       `json:"storage_driver"`
	Region            string                  `json:"cloud_provider_region"`
	ClusterType       consts.KsctlClusterType `json:"cluster_type"`
	BootstrapProvider consts.KsctlKubernetes  `json:"bootstrap_provider"`
	K8sVersion        string                  `json:"k8s_version"`
	Addons            []TelemetryAddon        `json:"addons"`
}

type TelemetryData struct {
	ClientIdentity string         `json:"client_identity"`
	UserId         string         `json:"client_id"`
	Event          TelemetryEvent `json:"event"`

	KsctlVer string        `json:"ksctl_ver"`
	OS       string        `json:"os"`
	Arch     string        `json:"arch"`
	Data     TelemetryMeta `json:"meta"`
}

type Telemetry struct {
	userId   string
	ksctlVer string
	endpoint string
	active   bool
	os       string
	arch     string
}

func NewTelemetry(active *bool) *Telemetry {
	return &Telemetry{
		userId:   "ksctl:cli",
		endpoint: "https://telemetry.ksctl.com",
		ksctlVer: config.Version,
		active:   active == nil || *active,
		os:       runtime.GOOS,
		arch:     runtime.GOARCH,
	}
}

func IntegrityCheck() {
	if config.InDevMode() {
		return
	}

	if len(config.Version) == 0 {
		color.New(color.BgHiRed, color.FgHiBlack).Println("corrupted version")
		os.Exit(2)
	}

	if len(clientIdentity) == 0 {
		color.New(color.BgHiRed, color.FgHiBlack).Println("corrupted cli identity")
		os.Exit(2)
	}

	if len(config.KsctlCoreVer) == 0 {
		color.New(color.BgHiRed, color.FgHiBlack).Println("corrupted core version")
		os.Exit(2)
	}

	if len(config.BuildDate) == 0 {
		color.New(color.BgHiRed, color.FgHiBlack).Println("corrupted build date")
		os.Exit(2)
	}
}

func (t *Telemetry) Send(ctx context.Context, l logger.Logger, event TelemetryEvent, data TelemetryMeta) error {
	if !t.active {
		return nil
	}

	telemetryData := TelemetryData{
		ClientIdentity: clientIdentity,
		UserId:         t.userId,
		KsctlVer:       t.ksctlVer,
		Event:          event,
		Data:           data,
		OS:             t.os,
		Arch:           t.arch,
	}

	if config.InDevMode() {
		return nil
	}

	payloadBuf := new(bytes.Buffer)

	if err := json.NewEncoder(payloadBuf).Encode(telemetryData); err != nil {
		return err
	}

	if res, err := http.Post(t.endpoint, "application/json", payloadBuf); err != nil {
		return err
	} else {
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send telemetry, status code: %d", res.StatusCode)
		}

		l.Debug(ctx, "Telemetry sent successfully")

		return nil
	}
}
