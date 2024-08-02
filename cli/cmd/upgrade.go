package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	"github.com/ksctl/cli/logger"
	"github.com/ksctl/ksctl/pkg/types"
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

func update(version string) error {
	return nil
}

var selfUpdate = &cobra.Command{
	Use:   "self-update",
	Short: "update the ksctl cli",
	Long:  "setups up update for ksctl cli",
	Run: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetInt("verbose")
		var log types.LoggerFactory = logger.NewLogger(verbosity, os.Stdout)

		// if Version == "dev" {
		// 	log.Error("Cannot update a dev version of ksctl")
		// 	os.Exit(1)
		// }

		vers, err := fetchLatestVersion()
		if err != nil {
			log.Error("Failed to fetch latest version", "error", err)
			os.Exit(1)
		}
		vers = filterToUpgradeableVersions(vers)

		log.Print(ctx, "Available versions to update")
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(vers).Show()

		newVer := selectedOption

		if err := update(newVer); err != nil {
		}

		log.Success(ctx, "Updated Ksctl cli", "previousVer", Version, "newVer", newVer)
	},
}

func init() {
	RootCmd.AddCommand(selfUpdate)
	storageFlag(selfUpdate)

	selfUpdate.Flags().BoolP("verbose", "v", true, "for verbose output")
}
