#!/bin/bash

rm -rf $HOME/.teritorid/

teritorid init --chain-id=teritori-1 testing --home=$HOME/.teritorid
teritorid prepare-genesis teritori-1 $PWD/cosmos_airdrop.csv $PWD/crew3_airdrop.csv $PWD/evmos_orbital_ape.csv
teritorid keys add validator --keyring-backend=test --home=$HOME/.teritorid
teritorid add-genesis-account $(teritorid keys show validator -a --keyring-backend=test --home=$HOME/.teritorid) 100000000000utori --home=$HOME/.teritorid
teritorid gentx validator 500000000utori --keyring-backend=test --home=$HOME/.teritorid --chain-id=teritori-1
teritorid collect-gentxs --home=$HOME/.teritorid

VALIDATOR=$(teritorid keys show -a validator --keyring-backend=test --home=$HOME/.teritorid)

sed -i '' -e 's/"owner": ""/"owner": "'$VALIDATOR'"/g' $HOME/.teritorid/config/genesis.json
sed -i '' -e 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/enable = false/enable = true/g' $HOME/.teritorid/config/app.toml 
sed -i '' -e 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' $HOME/.teritorid/config/config.toml 

teritorid start --home=$HOME/.teritorid

# teritorid tx bank send validator pop18mu5hhgy64390q56msql8pfwps0uesn0gf0elf 10000000utori --keyring-backend=test --chain-id=teritori-1 --home=$HOME/.teritorid/ --yes --broadcast-mode=block
# teritorid tx airdrop claim-allocation 0x9d967594Cc61453aFEfD657313e5F05be7c6F88F 0xb89733c05568385a861fa20f5c4abe53c23a13962515bf5510638b4e3947b1236963b53de549ae762bbd45427dbd3712ae7d169a935d21e44e7da86b1c552f471b --from=validator --keyring-backend=test --chain-id=testing --home=$HOME/.teritorid/ --yes  --broadcast-mode=block

# teritorid tx airdrop set-allocation evm 0x583e8DD54b7C3F5Ea23862E0E852f0e6914475D5 900000000utori 0utori --from=validator --keyring-backend=test --home=$HOME/.teritorid --chain-id=teritori-1 --broadcast-mode=block --yes
# teritorid tx airdrop set-allocation cosmos cosmos1hwf62gw7h39xmd69st3p487r8x3sphm24jfvxf 1100000000utori 0utori --from=validator --keyring-backend=test --home=$HOME/.teritorid --chain-id=teritori-1 --broadcast-mode=block --yes

# teritorid tx airdrop set-allocation evm 0xc1B929A2C4A7312f0A2c3841A9Ed13A96A81d922 9000000000utori 0utori --from=validator --keyring-backend=test --home=$HOME/.teritorid --chain-id=teritori-1 --broadcast-mode=block --yes
# teritorid tx airdrop set-allocation cosmos cosmos1ch8e2j5vdhtg4af7pr02g6ux6vnftr5nsqyhxa 1100000000utori 0utori --from=validator --keyring-backend=test --home=$HOME/.teritorid --chain-id=teritori-1 --broadcast-mode=block --yes
