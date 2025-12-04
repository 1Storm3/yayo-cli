#!/usr/bin/env bash

set -e

LATEST=$(curl -s https://api.github.com/repos/1Storm3/yayo-cli/releases/latest \
    | grep "browser_download_url" \
    | grep "linux-amd64" \
    | cut -d '"' -f 4)

curl -L $LATEST -o /usr/local/bin/yayo-cli
chmod +x /usr/local/bin/yayo-cli

echo "Успешно установлено: yayo-cli"
