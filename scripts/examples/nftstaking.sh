#!/usr/bin/env bash

nxtpopd tx nftstaking register-nft-staking --from validator --nft-identifier "identifier" --nft-metadata "metadata" --reward-address "pop1hw...tfh" --reward-weight 1000 --chain-id=testing --home=$HOME/.nxtpopd --keyring-backend=test --broadcast-mode=block --yes

nxtpopd query bank balances pop1hw...tfh
