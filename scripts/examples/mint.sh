#!/usr/bin/env bash

teritorid query bank balances $(teritorid keys show -a validator --keyring-backend=test)
teritorid tx mint burn-tokens 500000000stake --keyring-backend=test --from=validator --chain-id=testing --home=$HOME/.teritorid/ --yes  --broadcast-mode=block
teritorid query bank balances $(teritorid keys show -a validator --keyring-backend=test)
