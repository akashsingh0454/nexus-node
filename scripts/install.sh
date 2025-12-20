#!/bin/bash
# One-line installer for Nexus Node Agent (Linux/macOS)

REPO="YOUR_USERNAME/nexus-node" # PLACEHOLDER
INSTALL_DIR="/opt/nexus-node"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize Arch
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "🕷️  Installing Nexus Node Agent for $OS/$ARCH..."

# 1. Get Download URL
API_URL="https://api.github.com/repos/$REPO/releases/latest"
RELEASE_JSON=$(curl -s $API_URL)
DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep "browser_download_url" | grep "$OS-$ARCH" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Error: Could not find release asset for nexus-$OS-$ARCH"
    exit 1
fi

echo "   - Downloading from: $DOWNLOAD_URL"

# 2. Download & Install
mkdir -p "$INSTALL_DIR"
curl -L -o /tmp/nexus.tar.gz "$DOWNLOAD_URL"
tar -xzf /tmp/nexus.tar.gz -C "$INSTALL_DIR"
rm /tmp/nexus.tar.gz

chmod +x "$INSTALL_DIR/nexus"

echo "✅ Installation Complete!"
echo "   Location: $INSTALL_DIR/nexus"
echo "   Run with: $INSTALL_DIR/nexus --config $INSTALL_DIR/config.yaml"
