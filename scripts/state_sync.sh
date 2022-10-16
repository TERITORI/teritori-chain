#!/bin/bash
# microtick and bitcanna contributed significantly here.
# Pebbledb state sync script.
# invoke like: bash scripts/ss.bash
set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Install with pebbledb 
go mod edit -replace github.com/tendermint/tm-db=github.com/baabeetaa/tm-db@pebble
go mod tidy
go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb -X github.com/tendermint/tm-db.ForceSync=1' -tags pebbledb ./...

# go install ./...

# NOTE: ABOVE YOU CAN USE ALTERNATIVE DATABASES, HERE ARE THE EXACT COMMANDS
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags rocksdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags badgerdb ./...
# go install -ldflags '-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb' -tags boltdb ./...

rm  ~/.teritorid/config/genesis.json

# Initialize chain.
teritorid init test
teritorid config chain-id teritori-1

# Get Genesis
wget https://github.com/TERITORI/teritori-chain/blob/main/mainnet/teritori-1/genesis.json?raw=true -O genesis.json
mv genesis.json ~/.teritorid/config/genesis.json


# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT=$(curl -s https://teritori-rpc.polkachu.com/block | jq -r .result.block.header.height)
BLOCK_HEIGHT=$(($LATEST_HEIGHT-$INTERVAL)) 
TRUST_HASH=$(curl -s "https://teritori-rpc.polkachu.com/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export TERITORID_STATESYNC_ENABLE=true
export TERITORID_P2P_MAX_NUM_OUTBOUND_PEERS=200
export TERITORID_STATESYNC_RPC_SERVERS="https://teritori-rpc.polkachu.com:443,https://teritori-rpc.polkachu.com:443"
export TERITORID_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export TERITORID_STATESYNC_TRUST_HASH=$TRUST_HASH

# Fetch and set list of seeds from chain registry.
export TERITORID_P2P_SEEDS=$(curl -s https://raw.githubusercontent.com/cosmos/chain-registry/master/teritori/chain.json | jq -r '[foreach .peers.seeds[] as $item (""; "\($item.id)@\($item.address)")] | join(",")')

# Start chain.
teritorid start --x-crisis-skip-assert-invariants --db_backend pebbledb