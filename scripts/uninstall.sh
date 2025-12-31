#!/bin/bash
set -e

INSTALL_DIR="/opt/zen"
SERVICE_NAME="zen"

echo "Uninstalling Zen..."

if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo "Stopping $SERVICE_NAME service..."
    systemctl stop "$SERVICE_NAME"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
    echo "Disabling $SERVICE_NAME service..."
    systemctl disable "$SERVICE_NAME"
fi

if [ -f "/etc/systemd/system/$SERVICE_NAME.service" ]; then
    echo "Removing systemd service file..."
    rm -f "/etc/systemd/system/$SERVICE_NAME.service"
    systemctl daemon-reload
fi

if [ -d "$INSTALL_DIR" ]; then
    echo "Removing installation directory..."
    rm -rf "$INSTALL_DIR"
fi

echo "✓ Zen uninstalled successfully!"
