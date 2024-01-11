#!/bin/bash

teritorid140 tx gov submit-proposal software-upgrade "v2.0.0" \
--upgrade-height=10 \
--title="Upgrade to v2.0.0" --description="Upgrade to v2.0.0" \
--from=validator --keyring-backend=test \
--chain-id=testing --home=$HOME/.teritorid --yes -b block --deposit="100000000stake"

teritorid140 tx gov vote 1 yes --from validator --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend test

teritorid140 query gov proposals

sleep 50

killall teritorid140 &> /dev/null || true

teritorid start --home=$HOME/.teritorid

# Check mint module params update
teritorid query mint params

# Check packet forward queries
teritorid query packetforward params

# Check group module tx
teritorid query group groups
VALIDATOR=$(teritorid keys show -a validator --home $HOME/.teritorid --keyring-backend=test)

teritorid tx group create-group tori1uxqcel9xcdmx7rwfjgfrmdmzgmn7q3jql3cvhz "" group-members.json \
 --from validator --chain-id testing --keyring-backend=test \
 --chain-id=testing --home=$HOME/.teritorid --yes -b sync

# Check newly added queries
# - inflation
teritorid query mint inflation
# - stakingAPR
teritorid query mint staking-apr
# IBC transfer test
