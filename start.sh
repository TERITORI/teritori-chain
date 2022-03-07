#!/bin/bash

rm -rf $HOME/.nxtpopd/

cd $HOME

nxtpopd init --chain-id=testing testing --home=$HOME/.nxtpopd
nxtpopd keys add validator --keyring-backend=test --home=$HOME/.nxtpopd
nxtpopd add-genesis-account $(nxtpopd keys show validator -a --keyring-backend=test --home=$HOME/.nxtpopd) 100000000000pop,100000000000stake --home=$HOME/.nxtpopd
nxtpopd gentx validator 500000000stake --keyring-backend=test --home=$HOME/.nxtpopd --chain-id=testing
nxtpopd collect-gentxs --home=$HOME/.nxtpopd

nxtpopd start --home=$HOME/.nxtpopd
