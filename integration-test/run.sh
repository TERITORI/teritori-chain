#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=267a3f9604d48c6f1ea5c31d13a5d24c6ed35210

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