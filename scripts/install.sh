#!/usr/bin/env sh
# goclip installer
# Usage: curl -sSL https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/scripts/install.sh | sh

set -e

REPO="ashutoshsinghai/goclip"
BINARY="goclip"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)
    echo "❌ Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  arm64 | aarch64) ARCH="arm64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Decide install directory (SMART LOGIC)
if [ "$OS" = "darwin" ]; then
  if [ -d "/opt/homebrew/bin" ]; then
    INSTALL_DIR="/opt/homebrew/bin"
  else
    INSTALL_DIR="/usr/local/bin"
  fi
else
  INSTALL_DIR="/usr/local/bin"
fi

FALLBACK_DIR="$HOME/.local/bin"

# Get latest release version
echo "🔍 Fetching latest version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "❌ Could not fetch latest version."
  exit 1
fi

TARBALL="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"

echo "⬇️ Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

TMP=$(mktemp -d)
curl -fsSL "$URL" -o "${TMP}/${TARBALL}"
tar -xzf "${TMP}/${TARBALL}" -C "$TMP"

install_binary() {
  DEST="$1"

  if mv "${TMP}/${BINARY}" "${DEST}/${BINARY}" 2>/dev/null; then
    return 0
  fi

  return 1
}

# Try primary install dir (with sudo if needed)
if [ -w "$INSTALL_DIR" ]; then
  install_binary "$INSTALL_DIR"
else
  echo "🔐 Trying sudo install to $INSTALL_DIR..."
  if sudo mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}" 2>/dev/null; then
    :
  else
    echo "⚠️ Falling back to user install..."

    mkdir -p "$FALLBACK_DIR"
    mv "${TMP}/${BINARY}" "${FALLBACK_DIR}/${BINARY}"
    INSTALL_DIR="$FALLBACK_DIR"
  fi
fi

rm -rf "$TMP"

echo "✅ Installed to: ${INSTALL_DIR}/${BINARY}"

# Immediate usability check
if command -v "$BINARY" >/dev/null 2>&1; then
  echo "🎉 ${BINARY} is ready to use!"
else
  echo ""
  echo "⚠️ ${BINARY} is not in your PATH"
  echo ""
  echo "👉 Add this to your shell config (~/.zshrc or ~/.bashrc):"
  echo "export PATH=\"${INSTALL_DIR}:\$PATH\""
  echo ""
  echo "Then run:"
  echo "source ~/.zshrc  # or ~/.bashrc"
fi

echo ""

# Linux: remind about clipboard dependencies
if [ "$OS" = "linux" ]; then
  echo "📋 Linux clipboard note:"
  echo "   goclip needs xclip, xsel, or wl-clipboard to read your clipboard."
  echo "   Install one if you haven't already:"
  echo "     sudo apt install xclip       # X11"
  echo "     sudo apt install wl-clipboard # Wayland"
  echo ""
fi

echo "💡 goclip needs a background daemon to capture clipboard history."
echo "   Without it, goclip pick will show nothing."
echo ""

# Only prompt if stdin is a real terminal (not curl | sh pipe)
if [ -t 0 ] && command -v "$BINARY" >/dev/null 2>&1; then
  printf "Start the daemon now? [Y/n] "
  read answer </dev/tty
  case "$answer" in
    [nN]*)
      echo ""
      echo "   Start it later with: goclip daemon"
      ;;
    *)
      "$BINARY" daemon
      ;;
  esac
else
  echo "👉 Run this to start capturing clipboard history:"
  echo "   ${BINARY} daemon"
fi

echo ""
echo "👉 Then open the picker with:"
echo "   ${BINARY} pick"
echo ""
