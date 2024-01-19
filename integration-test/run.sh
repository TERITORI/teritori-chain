#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

commit=1d53bd16f9041de01089bdc0642868f9d65bbfad

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