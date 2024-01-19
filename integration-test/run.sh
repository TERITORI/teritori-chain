#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=51d92fc37f3326f916834c8677d646808a6d09c8

if [[ -z "${TERITORI_DAPP_REPO:-}" ]]; then
    rm -fr teritori-dapp
    git clone https://github.com/TERITORI/teritori-dapp.git
    cd teritori-dapp
    git checkout $commit
else
    cd $TERITORI_DAPP_REPO
fi

yarn

npx tsx packages/scripts/integration-testing/upgradeTest142toDir ..