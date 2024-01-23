#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

if [[ -z "${TERITORI_DAPP_REPO:-}" ]]; then
    commit=7e968801a0a03f47f59dd7683f1653935222ea88
    rm -fr teritori-dapp
    git clone https://github.com/TERITORI/teritori-dapp.git
    cd teritori-dapp
    git checkout $commit
else
    cd $TERITORI_DAPP_REPO
fi

yarn

npx tsx packages/scripts/integration-testing/simpleTest ..
npx tsx packages/scripts/integration-testing/upgradeTest ..