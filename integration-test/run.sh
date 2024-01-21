#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=bd009d1f49f0d82b3fcf81f741892af23acb5b1f

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