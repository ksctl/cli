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
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ksctl/cli/v2/pkg/cli"
	"github.com/ksctl/cli/v2/pkg/config"
	"github.com/ksctl/cli/v2/pkg/telemetry"
	"github.com/ksctl/ksctl/v2/pkg/poller"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func (k *KsctlCommand) CheckForUpdates() (bool, error) {

	cacheFile := &config.UpdateCache{}
	if errC := config.LoadUpdateCache(cacheFile); errC != nil {
		k.l.Error("Failed to load update cache", "error", errC)
		return false, errC
	}
	if config.InDevMode() {
		return false, nil
	}

	if !cacheFile.LastChecked.IsZero() && time.Since(cacheFile.LastChecked) < cacheFile.UpdateCheckInterval {
		return cacheFile.AvailableVersions, nil
	}

	versions, err := k.fetchLatestVersion()
	if err != nil {
		return false, err
	}

	upgradeableVersions := k.filterToUpgradeableVersions(versions)

	cacheFile.LastChecked = time.Now()
	cacheFile.AvailableVersions = len(upgradeableVersions) > 0

	if err := config.SaveUpdateCache(cacheFile); err != nil {
		k.l.Error("Failed to save update cache", "error", err)
		return false, err
	}

	return cacheFile.AvailableVersions, nil
}

func (k *KsctlCommand) NotifyAvailableUpdates() {
	k.l.Box(k.Ctx, "Update Available! ✨", "Run 'ksctl self-update' to upgrade to the latest version!")
}

func (k *KsctlCommand) SelfUpdate() *cobra.Command {

	cmd := &cobra.Command{
		Use: "self-update",
		Example: `
ksctl self-update --help
`,
		Short: "Use to update the ksctl cli",
		Long:  "It is used to update the ksctl cli",
		Run: func(cmd *cobra.Command, args []string) {

			if config.InDevMode() {
				k.l.Error("Cannot update dev version", "msg", "Please use a stable version to update")
				os.Exit(1)
			}

			k.l.Warn(k.Ctx, "Currently no migrations are supported", "msg", "Please help us by creating a PR to support migrations. Thank you!")

			k.l.Print(k.Ctx, "Fetching available versions")
			vers, err := k.fetchLatestVersion()
			if err != nil {
				k.l.Error("Failed to fetch latest version", "error", err)
				os.Exit(1)
			}
			vers = k.filterToUpgradeableVersions(vers)

			if len(vers) == 0 {
				k.l.Note(k.Ctx, "You are already on the latest version", "version", config.Version)
				os.Exit(0)
			}

			selectedOption, err := k.menuDriven.DropDownList("Select a version to update", vers, cli.WithDefaultValue(vers[0]))
			if err != nil {
				return
			}

			newVer := selectedOption

			if err := k.telemetry.Send(k.Ctx, k.l, telemetry.EventClusterUpgrade, telemetry.TelemetryMeta{}); err != nil {
				k.l.Debug(k.Ctx, "Failed to send the telemetry", "Reason", err)
			}

			{
				c := &config.UpdateCache{}
				err := config.LoadUpdateCache(c)
				if err == nil {
					c.LastChecked = time.Now()
					c.AvailableVersions = false
					if err := config.SaveUpdateCache(c); err != nil {
						k.l.Error("Failed to save update cache", "error", err)
					}
				} else {
					k.l.Error("Failed to load update cache", "error", err)
				}
			}

			if err := k.update(newVer); err != nil {
				k.l.Error("Failed to update ksctl cli", "error", err)
				os.Exit(1)
			}

			k.l.Box(k.Ctx, "Updated Successful 🎉", "ksctl has been updated to version "+newVer)
		},
	}

	return cmd
}

func (k *KsctlCommand) fetchLatestVersion() ([]string, error) {

	poller.InitSharedGithubReleasePoller()
	return poller.GetSharedPoller().Get("ksctl", "cli")
}

func (k *KsctlCommand) filterToUpgradeableVersions(versions []string) []string {
	var upgradeableVersions []string
	for _, version := range versions {
		if semver.Compare(version, config.Version) > 0 {
			upgradeableVersions = append(upgradeableVersions, version)
		}
	}
	return upgradeableVersions
}

func (k *KsctlCommand) downloadFile(url, localFilename string) error {
	k.l.Print(k.Ctx, "Downloading file", "url", url, "localFilename", localFilename)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
func (k *KsctlCommand) verifyChecksum(filePath, checksumfileLoc string) (bool, error) {
	k.l.Print(k.Ctx, "Verifying checksum", "file", filePath, "checksumfile", checksumfileLoc)

	rawChecksum, err := os.ReadFile(checksumfileLoc)
	if err != nil {
		return false, err
	}
	checksums := strings.Split(string(rawChecksum), "\n")

	var expectedChecksum string = "LOL"
	for _, line := range checksums {
		if strings.Contains(line, filePath) {
			expectedChecksum = strings.Fields(line)[0]
			break
		}
	}
	if expectedChecksum == "LOL" {
		return false, k.l.NewError(k.Ctx, "Checksum not found in checksum file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	calculatedChecksum := hex.EncodeToString(hash.Sum(nil))
	return calculatedChecksum == expectedChecksum, nil
}

func (k *KsctlCommand) getOsArch() (string, error) {
	arch := runtime.GOARCH

	if arch != "amd64" && arch != "arm64" {
		return "", k.l.NewError(k.Ctx, "Unsupported architecture")
	}
	return arch, nil
}

func (k *KsctlCommand) getOs() (string, error) {
	goos := runtime.GOOS

	if goos != "linux" && goos != "darwin" {
		return "", k.l.NewError(k.Ctx, "Unsupported OS", "message", "will provide support for windows based OS soon")
	}
	return goos, nil
}

func (k *KsctlCommand) update(version string) error {
	osName, err := k.getOs()
	if err != nil {
		return err
	}
	archName, err := k.getOsArch()
	if err != nil {
		return err
	}

	k.l.Print(k.Ctx, "Delected System", "OS", osName, "Arch", archName)
	downloadURLBase := fmt.Sprintf("https://github.com/ksctl/cli/releases/download/%s", version)
	tarFile := fmt.Sprintf("ksctl-cli_%s_%s_%s.tar.gz", version[1:], osName, archName)
	checksumFile := fmt.Sprintf("ksctl-cli_%s_checksums.txt", version[1:])

	tarUri := fmt.Sprintf("%s/%s", downloadURLBase, tarFile)
	checksumUri := fmt.Sprintf("%s/%s", downloadURLBase, checksumFile)

	defer func() {
		k.l.Print(k.Ctx, "Cleaning up")
		if err := os.Remove(checksumFile); err != nil {
			k.l.Error("Failed to remove checksum file", "error", err)
		}

		if err := os.Remove(tarFile); err != nil {
			k.l.Error("Failed to remove checksum file", "error", err)
		}
	}()

	if err := k.downloadFile(tarUri, tarFile); err != nil {
		return err
	}

	if err := k.downloadFile(checksumUri, checksumFile); err != nil {
		return err
	}

	match, err := k.verifyChecksum(tarFile, checksumFile)
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to verify checksum", "error", err)
	}
	if !match {
		return k.l.NewError(k.Ctx, "Checksum verification failed")
	}
	k.l.Success(k.Ctx, "Checksum verification successful")

	tempDir, err := os.MkdirTemp("", "ksctl-update")
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to create temp dir", "error", err)
	}
	file, err := os.Open(tarFile)
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to open tar file", "error", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to read gzip file", "error", err)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return k.l.NewError(k.Ctx, "Failed to read tar file", "error", err)
		}
		if header.Name == "ksctl" {
			outFile, err := os.Create(filepath.Join(tempDir, "ksctl"))
			if err != nil {
				return k.l.NewError(k.Ctx, "Failed to create ksctl binary", "error", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return k.l.NewError(k.Ctx, "Failed to copy ksctl binary", "error", err)
			}
			break
		}
	}

	k.l.Print(k.Ctx, "Making ksctl executable...")
	if err := os.Chmod(filepath.Join(tempDir, "ksctl"), 0550); err != nil {
		return k.l.NewError(k.Ctx, "Failed to make ksctl executable", "error", err)
	}

	k.l.Print(k.Ctx, "Moving ksctl to /usr/local/bin (requires sudo)...")
	cmd := exec.Command("sudo", "mv", "-v", filepath.Join(tempDir, "ksctl"), "/usr/local/bin/ksctl")
	err = cmd.Run()
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to move ksctl to /usr/local/bin", "error", err)
	}

	_, err = exec.LookPath("ksctl")
	if err != nil {
		return k.l.NewError(k.Ctx, "Failed to find ksctl in PATH", "error", err)
	}

	return nil
}
