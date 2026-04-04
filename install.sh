#!/bin/sh
# goclip installer
# Usage: curl -sSL https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/install.sh | sh

set -e

REPO="ashutoshsinghai/goclip"
BINARY="goclip"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  arm64 | aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get latest release version from GitHub API
VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Could not fetch latest version. Check your internet connection."
  exit 1
fi

TARBALL="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"

echo "Installing goclip ${VERSION} (${OS}/${ARCH})..."

TMP=$(mktemp -d)
curl -sSL "$URL" -o "${TMP}/${TARBALL}"
tar -xzf "${TMP}/${TARBALL}" -C "$TMP"

# Try /usr/local/bin, fall back to ~/bin
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  INSTALL_DIR="$HOME/bin"
  mkdir -p "$INSTALL_DIR"
  mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  echo "Installed to ${INSTALL_DIR} (no write access to /usr/local/bin)"
  echo "Make sure ${INSTALL_DIR} is in your PATH."
fi

rm -rf "$TMP"
echo "Done! Run: goclip help"
