#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=8630979356205f87e9c1e43e4dd1fdf6a0280273

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