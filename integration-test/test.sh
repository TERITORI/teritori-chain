#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

if [[ -z "${TERITORI_DAPP_REPO:-}" ]]; then
    rm -fr teritori-dapp
    git clone https://github.com/TERITORI/teritori-dapp.git
    cd teritori-dapp
    git checkout ac44e4e6ee0965d38da5c9736945c9ce08370ea6
else
    cd $TERITORI_DAPP_REPO
fi

yarn

while ! curl -s http://localhost:26657/status | jq -e '.result.sync_info.latest_block_height|tonumber > 0'; do sleep 5; done

npx tsx packages/scripts/network-setup/deploy teritori-localnet validator