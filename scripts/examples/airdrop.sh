#!/usr/bin/env bash

# query params
teritorid query airdrop params

# query allocation for XXX address
teritorid query airdrop allocation XXX --home=/Users/admin/.teritorid

# set allocation for XXX address
teritorid tx airdrop set-allocation evm 0x583e8DD54b7C3F5Ea23862E0E852f0e6914475D5 10000000utori 0utori --from=validator --keyring-backend=test --home=$HOME/.teritorid --chain-id=testing --broadcast-mode=block --yes

# claim allocation for XXX address
teritorid keys add acc1 --keyring-backend=test
teritorid tx airdrop claim-allocation 0x583e8DD54b7C3F5Ea23862E0E852f0e6914475D5 "" --from=acc1 --keyring-backend=test --home=$HOME/.teritorid --chain-id=testing --broadcast-mode=block

teritorid tx airdrop set-allocation evm 0x441D470F996D049B698A442e6DDb9dC3cb78AB68 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx airdrop set-allocation solana BW5D1Dv7pydTYZ8rSByEqNXYVRnGpm4qcEhkfHEGqBM 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx airdrop set-allocation terra BW5D1Dv7pydTYZ8rSByEqNXYVRnGpm4qcEhkfHEGqBM 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx airdrop set-allocation cosmos cosmos19ftk3lkfupgtnh38d7enc8c6jp7aljj3s0p6gt 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx airdrop set-allocation osmosis osmo19ftk3lkfupgtnh38d7enc8c6jp7aljj3c5j27e 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx airdrop set-allocation juno juno19ftk3lkfupgtnh38d7enc8c6jp7aljj3xazp0h 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes

teritorid tx bank send node0 pop19ftk3lkfupgtnh38d7enc8c6jp7aljj3qspa84 1stake --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
teritorid tx bank send node0 pop1hwf62gw7h39xmd69st3p487r8x3sphm29dftfh 1stake --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes

teritorid keys add node0 --keyring-backend=test --home=nodehome
teritorid tx airdrop set-allocation evm 0x0bEE910D7CFD039DD24178E2CE8C781f749A4791 100000000utori 0utori --from=node0 --keyring-backend=test --home=nodehome --chain-id=chain-JskdwJ --broadcast-mode=block --yes
