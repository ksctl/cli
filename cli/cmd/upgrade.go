package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ksctl/ksctl/poller"
	"github.com/rogpeppe/go-internal/semver"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

func fetchLatestVersion() ([]string, error) {

	logCli.Print(ctx, "Fetching available versions")

	poller.InitSharedGithubReleasePoller()
	return poller.GetSharedPoller().Get("ksctl", "cli")
}

func filterToUpgradeableVersions(versions []string) []string {
	var upgradeableVersions []string
	for _, version := range versions {
		if semver.Compare(version, Version) > 0 {
			upgradeableVersions = append(upgradeableVersions, version)
		}
	}
	return upgradeableVersions
}

func downloadFile(url, localFilename string) error {
	logCli.Print(ctx, "Downloading file", "url", url, "localFilename", localFilename)

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
func verifyChecksum(filePath, checksumfileLoc string) (bool, error) {
	logCli.Print(ctx, "Verifying checksum", "file", filePath, "checksumfile", checksumfileLoc)

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
		return false, logCli.NewError(ctx, "Checksum not found in checksum file")
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
	osName, err := getOs()
	if err != nil {
		return err
	}
	archName, err := getOsArch()
	if err != nil {
		return err
	}

	logCli.Print(ctx, "Delected System", "OS", osName, "Arch", archName)
	downloadURLBase := fmt.Sprintf("https://github.com/ksctl/cli/releases/download/%s", version)
	tarFile := fmt.Sprintf("ksctl-cli_%s_%s_%s.tar.gz", version[1:], osName, archName)
	checksumFile := fmt.Sprintf("ksctl-cli_%s_checksums.txt", version[1:])

	tarUri := fmt.Sprintf("%s/%s", downloadURLBase, tarFile)
	checksumUri := fmt.Sprintf("%s/%s", downloadURLBase, checksumFile)

	defer func() {
		logCli.Print(ctx, "Cleaning up")
		if err := os.Remove(checksumFile); err != nil {
			logCli.Error("Failed to remove checksum file", "error", err)
		}

		if err := os.Remove(tarFile); err != nil {
			logCli.Error("Failed to remove checksum file", "error", err)
		}
	}()

	if err := downloadFile(tarUri, tarFile); err != nil {
		return err
	}

	if err := downloadFile(checksumUri, checksumFile); err != nil {
		return err
	}

	match, err := verifyChecksum(tarFile, checksumFile)
	if err != nil {
		return logCli.NewError(ctx, "Failed to verify checksum", "error", err)
	}
	if !match {
		return logCli.NewError(ctx, "Checksum verification failed")
	}
	logCli.Success(ctx, "Checksum verification successful")

	tempDir, err := os.MkdirTemp("", "ksctl-update")
	if err != nil {
		return logCli.NewError(ctx, "Failed to create temp dir", "error", err)
	}
	file, err := os.Open(tarFile)
	if err != nil {
		return logCli.NewError(ctx, "Failed to open tar file", "error", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return logCli.NewError(ctx, "Failed to read gzip file", "error", err)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return logCli.NewError(ctx, "Failed to read tar file", "error", err)
		}
		if header.Name == "ksctl" {
			outFile, err := os.Create(filepath.Join(tempDir, "ksctl"))
			if err != nil {
				return logCli.NewError(ctx, "Failed to create ksctl binary", "error", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return logCli.NewError(ctx, "Failed to copy ksctl binary", "error", err)
			}
			break
		}
	}

	logCli.Print(ctx, "Making ksctl executable...")
	if err := os.Chmod(filepath.Join(tempDir, "ksctl"), 0550); err != nil {
		return logCli.NewError(ctx, "Failed to make ksctl executable", "error", err)
	}

	logCli.Print(ctx, "Moving ksctl to /usr/local/bin (requires sudo)...")
	cmd := exec.Command("sudo", "mv", "-v", filepath.Join(tempDir, "ksctl"), "/usr/local/bin/ksctl")
	err = cmd.Run()
	if err != nil {
		return logCli.NewError(ctx, "Failed to move ksctl to /usr/local/bin", "error", err)
	}

	_, err = exec.LookPath("ksctl")
	if err != nil {
		return logCli.NewError(ctx, "Failed to find ksctl in PATH", "error", err)
	}

	return nil
}

var selfUpdate = &cobra.Command{
	Use:   "self-update",
	Short: "update the ksctl cli",
	Long:  LongMessage("to self-update ksctl cli"),
	Run: func(cmd *cobra.Command, args []string) {

		if Version == "dev" {
			logCli.Error("Cannot update dev version", "msg", "Please use a stable version to update")
			os.Exit(1)
		}

		logCli.Warn(ctx, "Currently no migrations are supported", "msg", "Please help us by creating a PR to support migrations. Thank you!")

		vers, err := fetchLatestVersion()
		if err != nil {
			logCli.Error("Failed to fetch latest version", "error", err)
			os.Exit(1)
		}
		vers = filterToUpgradeableVersions(vers)

		if len(vers) == 0 {
			logCli.Success(ctx, "You are already on the latest version", "version", Version)
			os.Exit(0)
		}

		logCli.Print(ctx, "Available versions to update")
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(vers).Show()

		newVer := selectedOption

		if err := update(newVer); err != nil {
			logCli.Error("Failed to update ksctl cli", "error", err)
			os.Exit(1)
		}

		logCli.Success(ctx, "Updated Ksctl cli", "previousVer", Version, "newVer", newVer)
		logCli.Note(ctx, "Please restart your terminal to use the updated version")
	},
}

func init() {
	RootCmd.AddCommand(selfUpdate)
	storageFlag(selfUpdate)

	selfUpdate.Flags().BoolP("verbose", "v", true, "for verbose output")
}
