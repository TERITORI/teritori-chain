#!/bin/bash

teritorid131 tx gov submit-proposal software-upgrade "v1.4.0" \
--upgrade-height=15 \
--title="Upgrade to v1.4.0" --description="Upgrade to v1.4.0" \
--from=validator --keyring-backend=test \
--chain-id=testing --home=$HOME/.teritorid --yes -b block --deposit="100000000stake"

teritorid131 tx gov vote 1 yes --from validator --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend=test

teritorid131 query gov proposals

sleep 50

killall teritorid131 &> /dev/null || true

teritorid140 start --home=$HOME/.teritorid

# teritorid140 query ica controller params
# teritorid140 query bank balances $(teritorid140 keys show -a validator --keyring-backend=test)
# teritorid140 tx mint burn-tokens 500000000stake --from validator --chain-id testing --home $HOME/.teritorid -b block -y --keyring-backend=test

# Restore keys to hermes relayer
hermes --config ./network/hermes/config.toml keys delete --chain test-1 --all
hermes --config ./network/hermes/config.toml keys delete --chain testing --all
echo "alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart" > ./test-1.txt
hermes --config ./network/hermes/config.toml keys add --chain test-1 --mnemonic-file ./test-1.txt &
echo "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset" > ./teritori.txt
hermes --config ./network/hermes/config.toml keys add --chain testing --mnemonic-file ./teritori.txt &

# cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am
# tori17dtl0mjt3t77kpuhg2edqzjpszulwhgz7xjkfq

# teritorid140 query bank balances tori17dtl0mjt3t77kpuhg2edqzjpszulwhgz7xjkfq
# icad query bank balances cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am --node http://localhost:16657


teritorid140 tx bank send validator tori17dtl0mjt3t77kpuhg2edqzjpszulwhgz7xjkfq 10000000stake --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend=test

hermes --config ./network/hermes/config.toml create client --host-chain test-1 --reference-chain testing
hermes --config ./network/hermes/config.toml create client --host-chain testing --reference-chain test-1

hermes --config ./network/hermes/config.toml create channel --a-chain test-1 --b-chain testing --a-port transfer --b-port transfer --new-client-connection --yes
hermes --config ./network/hermes/config.toml start

teritorid140 tx intertx register --connection-id=connection-0 --from validator --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend=test

export VALIDATOR=$(teritorid keys show -a validator --home $HOME/.teritorid --keyring-backend=test)
teritorid query intertx interchainaccounts connection-0 $VALIDATOR
INTERCHAIN_ACCOUNT=cosmos1w8glpvlszm9ets53facn4hl69nnmvuyx4w7lykk4nd44354e6c2qasdsq3

icad tx bank send demowallet1 $INTERCHAIN_ACCOUNT 33333stake --keyring-backend=test --home=./data/test-1 --chain-id=test-1 --node=tcp://localhost:16657 --broadcast-mode=block -y
icad query bank balances $INTERCHAIN_ACCOUNT --node=tcp://localhost:16657
icad query bank balances cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw --node=tcp://localhost:16657

teritorid tx intertx submit \
'{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "from_address":"cosmos1w8glpvlszm9ets53facn4hl69nnmvuyx4w7lykk4nd44354e6c2qasdsq3",
    "to_address":"cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
    "amount": [
        {
            "denom": "stake",
            "amount": "1000"
        }
    ]
}' --connection-id connection-0 --from validator --chain-id testing \
--home $HOME/.teritorid -b block -y --keyring-backend=test