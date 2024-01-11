#!/bin/bash

rm -rf $HOME/.teritorid/

cd $HOME

teritorid140 init --chain-id=testing testing --home=$HOME/.teritorid
teritorid140 keys add validator --keyring-backend=test --home=$HOME/.teritorid
teritorid140 add-genesis-account $(teritorid140 keys show validator -a --keyring-backend=test --home=$HOME/.teritorid) 100000000000utori,100000000000stake --home=$HOME/.teritorid
teritorid140 gentx validator 500000000stake --keyring-backend=test --home=$HOME/.teritorid --chain-id=testing
teritorid140 collect-gentxs --home=$HOME/.teritorid

VALIDATOR=$(teritorid140 keys show -a validator --keyring-backend=test --home=$HOME/.teritorid)

sed -i '' -e 's/"owner": ""/"owner": "'$VALIDATOR'"/g' $HOME/.teritorid/config/genesis.json
sed -i '' -e 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/enable = false/enable = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' $HOME/.teritorid/config/config.toml 
sed -i '' 's/"voting_period": "172800s"/"voting_period": "20s"/g' $HOME/.teritorid/config/genesis.json

teritorid140 start --home=$HOME/.teritorid