#!/bin/bash

teritorid140 tx gov submit-proposal software-upgrade "v1.4.0" \
--upgrade-height=15 \
--title="Upgrade to v1.4.0" --description="Upgrade to v1.4.0" \
--from=validator --keyring-backend=test \
--chain-id=testing --home=$HOME/.teritorid --yes -b block --deposit="100000000stake"

teritorid140 tx gov vote 1 yes --from validator --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend test

teritorid140 query gov proposals

sleep 50

killall teritorid140 &> /dev/null || true

teritorid200 start --home=$HOME/.teritorid

# Check mint module params update
# Check packet forward queries
# Check group module tx
# Check newly added queries
# - inflation
# - stakingAPR
