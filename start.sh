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

teritorid start --home=$HOME/.teritorid

# teritorid tx bank send validator tori17wwqmzwqrksz2jrzy6zk2x2qhd0u77d8ryq6pv 10000000utori --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes --broadcast-mode=block
# teritorid tx airdrop claim-allocation 0x9d967594Cc61453aFEfD657313e5F05be7c6F88F 0xb89733c05568385a861fa20f5c4abe53c23a13962515bf5510638b4e3947b1236963b53de549ae762bbd45427dbd3712ae7d169a935d21e44e7da86b1c552f471b --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes  --broadcast-mode=block
# teritorid tx gov submit-proposal software-upgrade "v1.3.0" --upgrade-height=12 --title="title" --description="description" --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes --broadcast-mode=block --deposit="100000000stake"
# teritorid tx gov vote 1 Yes --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes  --broadcast-mode=block
# teritorid query gov proposals
