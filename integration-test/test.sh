#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

rm -fr teritori-dapp
git clone https://github.com/TERITORI/teritori-dapp.git
cd teritori-dapp
git checkout 267a3f9604d48c6f1ea5c31d13a5d24c6ed35210

yarn

while ! curl http://localhost:1317/node_info; do sleep 1; done

npx tsx packages/scripts/network-setup/deploy teritori-localnet validator