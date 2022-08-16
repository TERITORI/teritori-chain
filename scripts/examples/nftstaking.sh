#!/usr/bin/env bash

teritorid tx nftstaking register-nft-staking --from validator --nft-identifier "identifier3" --nft-metadata "metadata" --reward-address "pop1snktzg6rrncqtct3acx2vz60aak2a6fke3ny3c" --reward-weight 1000 --chain-id=testing --home=$HOME/.teritorid --keyring-backend=test --broadcast-mode=block --yes
teritorid tx nftstaking set-nft-type-perms NFT_TYPE_DEFAULT SET_SERVER_ACCESS --from=validator --chain-id=testing --home=$HOME/.teritorid --keyring-backend=test --broadcast-mode=block --yes
teritorid tx nftstaking set-access-info $(teritorid keys show -a validator --keyring-backend=test) server1#chan1#chan2,server2#chan3 --from=validator --chain-id=testing --home=$HOME/.teritorid --keyring-backend=test --broadcast-mode=block --yes

teritorid query bank balances pop1uef5c6tx7vhjyhfumhzdhvwkepshcmljyv4wh4
teritorid query nftstaking access-infos
teritorid query nftstaking access-info $(teritorid keys show -a validator --keyring-backend=test)
teritorid query nftstaking all-nfttype-perms
teritorid query nftstaking has-permission $(teritorid keys show -a validator --keyring-backend=test) aaa
teritorid query nftstaking nfttype-perms aaa
teritorid query nftstaking staking aaa
teritorid query nftstaking stakings
teritorid query nftstaking stakings_by_owner $(teritorid keys show -a validator --keyring-backend=test)

