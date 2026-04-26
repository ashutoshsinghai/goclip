# goclip installer for Windows (PowerShell)
# Usage: irm https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "ashutoshsinghai/goclip"
$BINARY = "goclip.exe"

# Detect architecture
$ARCH = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Host "Error: 32-bit Windows is not supported." -ForegroundColor Red
    exit 1
}

# Get latest release version from GitHub API
Write-Host "Fetching latest version..."
$release = Invoke-RestMethod "https://api.github.com/repos/$REPO/releases/latest"
$VERSION = $release.tag_name

if (-not $VERSION) {
    Write-Host "Error: Could not fetch latest version." -ForegroundColor Red
    exit 1
}

$ZIPNAME = "goclip_windows_${ARCH}.zip"
$URL = "https://github.com/$REPO/releases/download/$VERSION/$ZIPNAME"

Write-Host "Installing goclip $VERSION (windows/$ARCH)..."

# Download to temp folder
$TMP = Join-Path $env:TEMP "goclip_install"
New-Item -ItemType Directory -Force -Path $TMP | Out-Null

$ZIP_PATH = Join-Path $TMP $ZIPNAME
Invoke-WebRequest -Uri $URL -OutFile $ZIP_PATH

# Extract
Expand-Archive -Path $ZIP_PATH -DestinationPath $TMP -Force

# Install to ~/bin (no admin needed)
$INSTALL_DIR = Join-Path $env:USERPROFILE "bin"
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
Move-Item -Force (Join-Path $TMP $BINARY) (Join-Path $INSTALL_DIR $BINARY)

# Clean up
Remove-Item -Recurse -Force $TMP

# Add ~/bin to PATH for current user if not already there
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$INSTALL_DIR", "User")
    Write-Host "Added $INSTALL_DIR to your PATH." -ForegroundColor Yellow
    Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Done! goclip is installed." -ForegroundColor Green
Write-Host "It will also be set to start automatically at login." -ForegroundColor Green
Write-Host ""

# Prompt the user to start goclip now. Read-Host works under `irm | iex`
# because PowerShell connects it to the host console rather than the input
# stream. If we're running non-interactively (no host UI), Read-Host throws —
# in that case fall back to printing manual instructions.
$BINARY_PATH = Join-Path $INSTALL_DIR $BINARY
try {
    $answer = Read-Host "Start goclip now? [Y/n]"
    if ($answer -match '^[nN]') {
        Write-Host "  Start it later with: goclip start" -ForegroundColor DarkGray
    } else {
        & $BINARY_PATH start
    }
} catch {
    Write-Host "Run this to start capturing clipboard history:" -ForegroundColor Yellow
    Write-Host "  goclip start"
}

Write-Host ""
Write-Host "Then open the picker with: goclip pick"
