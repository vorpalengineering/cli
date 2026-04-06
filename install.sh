#!/bin/sh
set -e

REPO="vorpalengineering/cli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Check dependencies
for cmd in curl tar; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Error: $cmd is required but not installed." >&2
    exit 1
  fi
done

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *)
    echo "Error: unsupported operating system: $OS" >&2
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  arm64|aarch64) ARCH="arm64" ;;
  x86_64)        ARCH="amd64" ;;
  *)
    echo "Error: unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

echo "Detected platform: ${OS}-${ARCH}"

# Get latest release version
echo "Fetching latest release..."
VERSION="$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"

if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest release version." >&2
  exit 1
fi

echo "Latest version: ${VERSION}"

# Download
URL="https://github.com/${REPO}/releases/download/${VERSION}/vorpal-${VERSION}-${OS}-${ARCH}.tar.gz"
TMPDIR="$(mktemp -d)"
TARBALL="${TMPDIR}/vorpal.tar.gz"

echo "Downloading ${URL}..."
if ! curl -sSfL -o "$TARBALL" "$URL"; then
  echo "Error: download failed. Check that a release exists for ${OS}-${ARCH}." >&2
  rm -rf "$TMPDIR"
  exit 1
fi

# Extract
tar -xzf "$TARBALL" -C "$TMPDIR"

# Install
if [ ! -w "$INSTALL_DIR" ]; then
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/vorpal" "${INSTALL_DIR}/vorpal"
else
  mv "${TMPDIR}/vorpal" "${INSTALL_DIR}/vorpal"
fi

# Cleanup
rm -rf "$TMPDIR"

# Verify
if command -v vorpal >/dev/null 2>&1; then
  echo ""
  echo "Installed successfully!"
  vorpal version
else
  echo ""
  echo "Installed to ${INSTALL_DIR}/vorpal"
  echo "Make sure ${INSTALL_DIR} is in your PATH."
fi
