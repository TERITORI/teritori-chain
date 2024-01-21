#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=66134e9580135a07aba64e00b68af9f30f8fdb93

if [[ -z "${TERITORI_DAPP_REPO:-}" ]]; then
    rm -fr teritori-dapp
    git clone https://github.com/TERITORI/teritori-dapp.git
    cd teritori-dapp
    git checkout $commit
else
    cd $TERITORI_DAPP_REPO
fi

yarn

npx tsx packages/scripts/integration-testing/simpleTest ..
npx tsx packages/scripts/integration-testing/upgradeTest142toDir ..
npx tsx packages/scripts/integration-testing/upgradeTest120toDir ..