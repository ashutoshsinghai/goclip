package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const repo = "ashutoshsinghai/goclip"

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Upgrade fetches the latest release from GitHub and replaces the current binary.
func Upgrade(currentVersion string) {
	fmt.Println("Checking for updates...")

	// Fetch latest release info from GitHub API
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		fmt.Printf("Error reaching GitHub: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		fmt.Printf("Error reading release info: %v\n", err)
		os.Exit(1)
	}

	latest := release.TagName

	// Compare versions
	if currentVersion != "" && currentVersion == latest {
		fmt.Printf("Already up to date (%s)\n", currentVersion)
		return
	}
	fmt.Printf("Upgrading %s → %s\n", currentVersion, latest)

	// Detect OS and arch
	goos := runtime.GOOS     // "darwin", "linux", "windows"
	goarch := runtime.GOARCH // "amd64", "arm64"

	// Build the expected asset name (matches goreleaser output)
	var assetName string
	if goos == "windows" {
		assetName = fmt.Sprintf("goclip_%s_%s.zip", goos, goarch)
	} else {
		assetName = fmt.Sprintf("goclip_%s_%s.tar.gz", goos, goarch)
	}

	// Find the matching asset URL
	downloadURL := ""
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		fmt.Printf("No binary found for %s/%s in release %s\n", goos, goarch, latest)
		os.Exit(1)
	}

	fmt.Printf("Downloading %s...\n", assetName)

	// Download the archive to a temp file
	tmpDir, _ := os.MkdirTemp("", "goclip-upgrade")
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		fmt.Printf("Download failed: %v\n", err)
		os.Exit(1)
	}

	// Extract the binary from the archive
	binaryName := "goclip"
	if goos == "windows" {
		binaryName = "goclip.exe"
	}

	newBinaryPath := filepath.Join(tmpDir, binaryName)
	if strings.HasSuffix(archivePath, ".tar.gz") {
		if err := extractTarGz(archivePath, binaryName, newBinaryPath); err != nil {
			fmt.Printf("Extraction failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := extractZip(archivePath, binaryName, newBinaryPath); err != nil {
			fmt.Printf("Extraction failed: %v\n", err)
			os.Exit(1)
		}
	}

	// Find where the current binary lives
	currentBinary, err := os.Executable()
	if err != nil {
		fmt.Printf("Could not find current binary path: %v\n", err)
		os.Exit(1)
	}
	currentBinary, _ = filepath.EvalSymlinks(currentBinary)

	// Replace: write new binary next to the old one, then rename over it
	tmpBinary := currentBinary + ".new"
	if err := copyFile(newBinaryPath, tmpBinary); err != nil {
		fmt.Printf("Could not write new binary: %v\n", err)
		os.Exit(1)
	}
	os.Chmod(tmpBinary, 0755)

	if err := os.Rename(tmpBinary, currentBinary); err != nil {
		fmt.Printf("Could not replace binary (try with sudo): %v\n", err)
		os.Remove(tmpBinary)
		os.Exit(1)
	}

	fmt.Printf("Done! goclip upgraded to %s\n", latest)
}

// downloadFile downloads a URL to a local file.
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
