#!/bin/bash

rm -rf $HOME/.nxtpopd/

cd $HOME

nxtpopd init --chain-id=testing testing --home=$HOME/.nxtpopd
nxtpopd keys add validator --keyring-backend=test --home=$HOME/.nxtpopd
nxtpopd add-genesis-account $(nxtpopd keys show validator -a --keyring-backend=test --home=$HOME/.nxtpopd) 100000000000upop,100000000000stake --home=$HOME/.nxtpopd
nxtpopd gentx validator 500000000stake --keyring-backend=test --home=$HOME/.nxtpopd --chain-id=testing
nxtpopd collect-gentxs --home=$HOME/.nxtpopd

sed -i '' -e 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.nxtpopd/config/app.toml 
sed -i '' -e 's/enable = false/enable = true/g' $HOME/.nxtpopd/config/app.toml 

nxtpopd start --home=$HOME/.nxtpopd

# nxtpopd tx bank send validator pop1pkmvlnstq8q7djns3w882pcu92xh4c9x8hpevr 10000000upop --keyring-backend=test --chain-id=testing --home=$HOME/.nxtpopd/ --yes