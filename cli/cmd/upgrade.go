package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

func fetchLatestVersion() ([]string, error) {

	type Release struct {
		TagName string `json:"tag_name"`
	}

	resp, err := http.Get("https://api.github.com/repos/ksctl/cli/releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var _releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&_releases); err != nil {
		return nil, err
	}

	var releases []string
	rcRegex := regexp.MustCompile(`.*-rc[0-9]+$`)
	for _, release := range _releases {
		if rcRegex.MatchString(release.TagName) {
			continue
		}
		releases = append(releases, release.TagName)
	}

	return releases, nil
}

func filterToUpgradeableVersions(versions []string) []string {
	var upgradeableVersions []string
	for _, version := range versions {
		if version > Version {
			upgradeableVersions = append(upgradeableVersions, version)
		}
	}
	return upgradeableVersions
}

func downloadFile(url, localFilename string) error {
	fmt.Printf("Downloading %s to %s\n", url, localFilename)
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
func verifyChecksum(filePath, expectedChecksum string) (bool, error) {
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

func getOsArch() (string, error) {
	arch := runtime.GOARCH

	if arch != "amd64" && arch != "arm64" {
		return "", logCli.NewError(ctx, "Unsupported architecture")
	}
	return arch, nil
}

func getOs() (string, error) {
	os := runtime.GOOS

	if os != "linux" && os != "darwin" {
		return "", logCli.NewError(ctx, "Unsupported OS", "message", "will provide support for windows based OS soon")
	}
	return os, nil
}

func update(version string) error {
	os, err := getOs()
	if err != nil {
		return err
	}
	arch, err := getOsArch()
	if err != nil {
		return err
	}
	return nil
}

var selfUpdate = &cobra.Command{
	Use:   "self-update",
	Short: "update the ksctl cli",
	Long:  "setups up update for ksctl cli",
	Run: func(cmd *cobra.Command, args []string) {

		// if Version == "dev" {
		// 	log.Error("Cannot update a dev version of ksctl")
		// 	os.Exit(1)
		// }

		logCli.Warn(ctx, "Currently no migrations are supported", "Message", "Please help us by creating a PR to support migrations. Thank you!")

		vers, err := fetchLatestVersion()
		if err != nil {
			logCli.Error("Failed to fetch latest version", "error", err)
			os.Exit(1)
		}
		vers = filterToUpgradeableVersions(vers)

		logCli.Print(ctx, "Available versions to update")
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(vers).Show()

		newVer := selectedOption

		if err := update(newVer); err != nil {
			logCli.Error("Failed to update ksctl cli", "error", err)
			os.Exit(1)
		}

		logCli.Success(ctx, "Updated Ksctl cli", "previousVer", Version, "newVer", newVer)
	},
}

func init() {
	RootCmd.AddCommand(selfUpdate)
	storageFlag(selfUpdate)

	selfUpdate.Flags().BoolP("verbose", "v", true, "for verbose output")
}
