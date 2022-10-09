# Using teritorid command line client

## Download go

Download go from [here](https://go.dev/dl/) and install it.

## Build teritorid binary

Run `go build` in this directory. A binary file named `teritorid` will be created in this directory.

## Usage of teritorid

The following are some examples of how you can use teritorid to interact with Teritori network. Every subcommand can be explored by: `teritorid <subcommand> --help`

### Creating wallet

This will create a new Teritori wallet locally on your computer `teritorid keys add <WALLET_NAME>`

### Querying an address on mainnet

`teritorid query account <ADDRESS> --node "https://rpc.mainnet.teritori.com:443" --chain-id "teritori-1"` 

### Staking to a specific validator

`teritorid tx staking  delegate <VALIDATOR_ADDRESS> "<amount>utori" --from <WALLET_NAME> --node "https://rpc.mainnet.teritori.com:443" --chain-id "teritori-1"`
