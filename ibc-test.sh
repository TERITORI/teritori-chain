make init-hermes
sh ./network/hermes/restore-keys.sh
sh ./network/hermes/create-conn.sh
make start-hermes

# create interchain account from icad
icad tx intertx register --from demowallet1 --connection-id connection-0 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:16657 --keyring-backend test -y --broadcast-mode=block

teritorid tx intertx register --from demowallet2 --connection-id connection-0 --chain-id test-2 --home ./data/test-2 --keyring-backend test -y --broadcast-mode=block

# Query the address of the interchain account
export DEMOWALLET_1=$(icad keys show demowallet1 -a --keyring-backend test --home ./data/test-1) && echo $DEMOWALLET_1;
icad query intertx interchainaccounts connection-0 $DEMOWALLET_1 --home ./data/test-1 --node tcp://localhost:16657

export DEMOWALLET_2=$(teritorid keys show -a demowallet2 --home ./data/test-2 --keyring-backend test)
teritorid query intertx interchainaccounts connection-0 $DEMOWALLET_2
INTERCHAIN_ACCOUNT=cosmos1zz7xehvr6355sr2mkjeap3gpeaq84uul6r4040wfn3w2h2kcw4jsqsye6k

icad tx bank send demowallet1 $INTERCHAIN_ACCOUNT 33333stake --keyring-backend=test --home=./data/test-1 --chain-id=test-1 --node=tcp://localhost:16657 --broadcast-mode=block -y
icad query bank balances $INTERCHAIN_ACCOUNT --node=tcp://localhost:16657
icad query bank balances cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw --node=tcp://localhost:16657

teritorid tx intertx submit \
'{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "from_address":"cosmos1zz7xehvr6355sr2mkjeap3gpeaq84uul6r4040wfn3w2h2kcw4jsqsye6k",
    "to_address":"cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
    "amount": [
        {
            "denom": "stake",
            "amount": "1000"
        }
    ]
}' --connection-id connection-0 --from $DEMOWALLET_2 --chain-id test-2 --home ./data/test-2 --keyring-backend test -y --broadcast-mode=block

# hermes -c ./network/hermes/config.toml create channel test-1 test-2 --port-a transfer --port-b transfer ;

# teritorid q ibc channel channels;

# icad keys list --keyring-backend=test --home=./data/test-1

# icad tx ibc-transfer transfer transfer channel-1 tori1t7cyvydpp4lklprksnrjy2y3xzv3q2l0n4qqvn 1000000stake --chain-id=test-1 --from=demowallet1 --keyring-backend=test -y --broadcast-mode=block --home=./data/test-1 --node=http://localhost:16657

# teritorid query bank balances tori1t7cyvydpp4lklprksnrjy2y3xzv3q2l0n4qqvn 

