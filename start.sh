#!/bin/bash

rm -rf $HOME/.teritorid/

cd $HOME

teritorid init --chain-id=testing testing --home=$HOME/.teritorid
teritorid keys add validator --keyring-backend=test --home=$HOME/.teritorid
teritorid add-genesis-account $(teritorid keys show validator -a --keyring-backend=test --home=$HOME/.teritorid) 100000000000utori,100000000000stake --home=$HOME/.teritorid
teritorid gentx validator 500000000stake --keyring-backend=test --home=$HOME/.teritorid --chain-id=testing
teritorid collect-gentxs --home=$HOME/.teritorid

VALIDATOR=$(teritorid keys show -a validator --keyring-backend=test --home=$HOME/.teritorid)

sed -i '' -e 's/"owner": ""/"owner": "'$VALIDATOR'"/g' $HOME/.teritorid/config/genesis.json
sed -i '' -e 's/"voting_period": "172800s"/"voting_period": "20s"/g' $HOME/.teritorid/config/genesis.json
sed -i '' -e 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/enable = false/enable = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' $HOME/.teritorid/config/config.toml 
jq '.app_state.gov.voting_params.voting_period = "20s"'  $HOME/.teritorid/config/genesis.json > temp.json ; mv temp.json $HOME/.teritorid/config/genesis.json;

teritorid start --home=$HOME/.teritorid

# git checkout v1.3.0
# go install ./cmd/teritorid
# sh start.sh
# teritorid tx gov submit-proposal software-upgrade "v1.4.0" --upgrade-height=12 --title="title" --description="description" --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes --broadcast-mode=block --deposit="100000000stake"
# teritorid tx gov vote 1 Yes --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes  --broadcast-mode=block
# teritorid query gov proposals
# git checkout ica_controller
# go install ./cmd/teritorid
# teritorid start
# teritorid query interchain-accounts controller params
