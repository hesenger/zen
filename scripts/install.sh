#!/bin/bash
set -e

REPO="hesenger/zen"
BINARY_NAME="app"
INSTALL_DIR="/opt/zen"
SERVICE_NAME="zen"

echo "Installing Zen from GitHub releases..."

echo "Creating installation directory..."
mkdir -p "$INSTALL_DIR/data"

LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not fetch latest release"
    exit 1
fi

echo "Latest release: $LATEST_RELEASE"

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_NAME"

echo "Downloading $DOWNLOAD_URL..."
curl -L -o "/tmp/$BINARY_NAME" "$DOWNLOAD_URL"

chmod +x "/tmp/$BINARY_NAME"

echo "Installing to $INSTALL_DIR..."
mv "/tmp/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "Setting up PATH..."
if ! grep -q "$INSTALL_DIR" /etc/environment; then
    echo "PATH=\"\$PATH:$INSTALL_DIR\"" >> /etc/environment
fi

echo "Downloading systemd service file..."
curl -L -o "/tmp/$SERVICE_NAME.service" "https://raw.githubusercontent.com/$REPO/main/scripts/app.service"

echo "Installing systemd service..."
mv "/tmp/$SERVICE_NAME.service" "/etc/systemd/system/$SERVICE_NAME.service"

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"

echo "✓ Zen installed successfully!"
echo "Service status:"
systemctl status "$SERVICE_NAME" --no-pager
