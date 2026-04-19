package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ashutoshsinghai/goclip/internal/style"
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

// Upgrade fetches the latest release and replaces the current binary.
func Upgrade(currentVersion string) {
	fmt.Println(style.Dim.Render("Checking for updates..."))
	release, err := fetchRelease("latest")
	if err != nil {
		fmt.Println(style.Red.Render("Error reaching GitHub: ") + err.Error())
		os.Exit(1)
	}
	if currentVersion != "" && currentVersion == release.TagName {
		fmt.Println(style.Green.Render("Already up to date ") + style.Dim.Render("("+currentVersion+")"))
		return
	}
	fmt.Println(style.Bold.Render("Upgrading ") + style.Dim.Render(currentVersion) + style.Bold.Render(" → ") + style.Green.Render(release.TagName))
	applyRelease(release)
	fmt.Println(style.Green.Render("Done! goclip upgraded to ") + style.Bold.Render(release.TagName))
}

// Install downloads and installs a specific version, or the latest if "--latest" is passed.
func Install(version, currentVersion string) {
	// normalise: accept both "v0.5.0" and "0.5.0"
	tag := version
	if tag == "--latest" || tag == "latest" {
		tag = "latest"
	} else if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}

	fmt.Println(style.Dim.Render("Fetching release info..."))
	release, err := fetchRelease(tag)
	if err != nil {
		fmt.Println(style.Red.Render("Error: ") + err.Error())
		os.Exit(1)
	}

	if currentVersion != "" && currentVersion == release.TagName {
		fmt.Println(style.Green.Render("Already on ") + style.Bold.Render(release.TagName))
		return
	}

	if tag == "latest" {
		fmt.Println(style.Bold.Render("Installing ") + style.Dim.Render(currentVersion) + style.Bold.Render(" → ") + style.Green.Render(release.TagName))
	} else {
		fmt.Println(style.Bold.Render("Installing ") + style.Green.Render(release.TagName) + style.Dim.Render(" (current: "+currentVersion+")"))
	}

	applyRelease(release)
	fmt.Println(style.Green.Render("Done! goclip is now at ") + style.Bold.Render(release.TagName))
}

// fetchRelease fetches release metadata from GitHub.
// Pass "latest" for the latest release, or a tag like "v0.5.0" for a specific one.
func fetchRelease(tag string) (githubRelease, error) {
	var url string
	if tag == "latest" {
		url = "https://api.github.com/repos/" + repo + "/releases/latest"
	} else {
		url = "https://api.github.com/repos/" + repo + "/releases/tags/" + tag
	}

	resp, err := http.Get(url)
	if err != nil {
		return githubRelease{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return githubRelease{}, fmt.Errorf("version %q not found — check available releases at https://github.com/%s/releases", tag, repo)
	}
	if resp.StatusCode != 200 {
		return githubRelease{}, fmt.Errorf("GitHub returned HTTP %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return githubRelease{}, fmt.Errorf("could not parse release info: %w", err)
	}
	return release, nil
}

// applyRelease downloads the right asset for the current platform and replaces the binary.
func applyRelease(release githubRelease) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var assetName string
	if goos == "windows" {
		assetName = fmt.Sprintf("goclip_%s_%s.zip", goos, goarch)
	} else {
		assetName = fmt.Sprintf("goclip_%s_%s.tar.gz", goos, goarch)
	}

	downloadURL := ""
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		fmt.Printf("No binary found for %s/%s in release %s\n", goos, goarch, release.TagName)
		os.Exit(1)
	}

	fmt.Println(style.Dim.Render("Downloading " + assetName + "..."))

	tmpDir, _ := os.MkdirTemp("", "goclip-install")
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		fmt.Printf("Download failed: %v\n", err)
		os.Exit(1)
	}

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

	currentBinary, err := os.Executable()
	if err != nil {
		fmt.Printf("Could not find current binary path: %v\n", err)
		os.Exit(1)
	}
	currentBinary, _ = filepath.EvalSymlinks(currentBinary)

	tmpBinary := currentBinary + ".new"
	if err := copyFile(newBinaryPath, tmpBinary); err != nil {
		fmt.Printf("Could not write new binary: %v\n", err)
		os.Exit(1)
	}
	os.Chmod(tmpBinary, 0755)

	if err := replaceBinary(tmpBinary, currentBinary); err != nil {
		fmt.Printf("Could not replace binary: %v\n", err)
		os.Remove(tmpBinary)
		os.Exit(1)
	}
}

// replaceBinary swaps tmpBinary into place as dest.
// On Windows, the running executable can't be overwritten directly, so we
// rename it to a .old file first, then rename the new binary into place.
func replaceBinary(src, dest string) error {
	if runtime.GOOS != "windows" {
		return os.Rename(src, dest)
	}
	old := dest + ".old"
	os.Remove(old) // remove any leftover from a previous install
	if err := os.Rename(dest, old); err != nil {
		return err
	}
	if err := os.Rename(src, dest); err != nil {
		// best-effort rollback
		os.Rename(old, dest)
		return err
	}
	os.Remove(old)
	return nil
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
