#!/usr/bin/env bash

# query allocation for XXX address
nxtpopd query airdrop allocation XXX --home=/Users/admin/.nxtpopd

# set allocation for XXX address
nxtpopd tx airdrop set-allocation evm 0x583e8DD54b7C3F5Ea23862E0E852f0e6914475D5 10000000upop 0upop --from=validator --keyring-backend=test --home=$HOME/.nxtpopd --chain-id=testing --broadcast-mode=block --yes

# claim allocation for XXX address
nxtpopd tx airdrop claim-allocation XXX "" --from=validator --keyring-backend=test --home=$HOME/.nxtpopd --chain-id=testing --broadcast-mode=block