#!/bin/bash

rm -rf $HOME/.teritorid/

cd $HOME

teritorid131 init --chain-id=testing testing --home=$HOME/.teritorid
teritorid131 keys add validator --keyring-backend=test --home=$HOME/.teritorid
teritorid131 add-genesis-account $(teritorid131 keys show validator -a --keyring-backend=test --home=$HOME/.teritorid) 100000000000utori,100000000000stake --home=$HOME/.teritorid
teritorid131 gentx validator 500000000stake --keyring-backend=test --home=$HOME/.teritorid --chain-id=testing
teritorid131 collect-gentxs --home=$HOME/.teritorid

VALIDATOR=$(teritorid131 keys show -a validator --keyring-backend=test --home=$HOME/.teritorid)

sed -i '' -e 's/"owner": ""/"owner": "'$VALIDATOR'"/g' $HOME/.teritorid/config/genesis.json
sed -i '' -e 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/enable = false/enable = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' $HOME/.teritorid/config/config.toml 
sed -i '' 's/"voting_period": "172800s"/"voting_period": "20s"/g' $HOME/.teritorid/config/genesis.json

teritorid131 start --home=$HOME/.teritorid