#!/usr/bin/env bash

# query allocation for XXX address
nxtpopd query airdrop allocation XXX --home=/Users/admin/.nxtpopd

# claim allocation for XXX address
nxtpopd tx airdrop claim-allocation XXX "" --from=validator --keyring-backend=test --home=$HOME/.nxtpopd --chain-id=testing --broadcast-mode=block