#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=1366a1f06d43e9d1cd53e3dd022df3b4ee47c8d3

if [[ -z "${TERITORI_DAPP_REPO:-}" ]]; then
    rm -fr teritori-dapp
    git clone https://github.com/TERITORI/teritori-dapp.git
    cd teritori-dapp
    git checkout $commit
else
    cd $TERITORI_DAPP_REPO
fi

yarn

while ! curl -s http://localhost:26657/status | jq -e '.result.sync_info.latest_block_height|tonumber > 0'; do sleep 5; done

npx tsx packages/scripts/network-setup/deploy teritori-localnet validator