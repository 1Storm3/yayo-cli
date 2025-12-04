#!/usr/bin/env bash

set -e

# ОС
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    Linux)
        OS="linux"
        ;;
    Darwin)
        OS="darwin"
        ;;
    *)
        echo "❌ Unsupported OS: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

ASSET_NAME="yayo-cli_${OS}-${ARCH}"

echo "→ Определена платформа: $ASSET_NAME"

LATEST=$(curl -s https://api.github.com/repos/1Storm3/yayo-cli/releases/latest \
    | grep "browser_download_url" \
    | grep "$ASSET_NAME" \
    | cut -d '"' -f 4)

if [ -z "$LATEST" ]; then
    echo "❌ Не найден бинарь для $ASSET_NAME в последнем релизе"
    exit 1
fi

echo "→ Загружаю: $LATEST"

sudo curl -L "$LATEST" -o /usr/local/bin/yayo-cli
sudo chmod +x /usr/local/bin/yayo-cli

echo "✅ Успешно установлено: yayo-cli"
echo "Используйте: yayo-cli --help"
