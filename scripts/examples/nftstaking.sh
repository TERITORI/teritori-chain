#!/usr/bin/env bash

nxtpopd tx nftstaking register-nft-staking --from validator --nft-identifier "identifier3" --nft-metadata "metadata" --reward-address "pop1snktzg6rrncqtct3acx2vz60aak2a6fke3ny3c" --reward-weight 1000 --chain-id=testing --home=$HOME/.nxtpopd --keyring-backend=test --broadcast-mode=block --yes

nxtpopd query bank balances pop1uef5c6tx7vhjyhfumhzdhvwkepshcmljyv4wh4
