#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=8853dff7dd50ca7f3f4a2c444eede38f5a40d0bc

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