# [DEPRECIATED]

## Setup your machine

If you already have go 1.18+ and packages up to date, you can skip this part and jump to the second section: [Setup the chain](https://github.com/TERITORI/teritori-chain/edit/main/testnet/testnet-teritori-v1/README.md#setup-the-chain)  
Make sure your machine is up to date:  
```shell
apt update && apt upgrade -y 
```  

Install few packages:  
```shell
apt install build-essential git curl gcc make jq -y
```

Install Go 1.18+:  
```shell
wget -c https://go.dev/dl/go1.18.3.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz && rm -rf go1.18.3.linux-amd64.tar.gz
``` 

Setup your environnement:  
```shell
echo 'export GOROOT=/usr/local/go' >> $HOME/.bash_profile
echo 'export GOPATH=$HOME/go' >> $HOME/.bash_profile
echo 'export GO111MODULE=on' >> $HOME/.bash_profile
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> $HOME/.bash_profile && . $HOME/.bash_profile
```  

Verify the installation:  
```shell
go version
#Should return go version go1.18.3 linux/amd64
``` 

## Setup the chain  

Clone the Teritori repository and install the v1 of testnet:  
```shell
git clone https://github.com/TERITORI/teritori-chain && cd teritori-chain && git checkout teritori-testnet-v1 && make install
```  

Verify the installation:  
```shell
teritorid version
#Should return  teritori-testnet-v1-9f1943df9cff63ca5e4ec10d0e31521e13c2fc13
```  

Init the chain:
```shell
teritorid init <YOUR_MONIKER> --chain-id teritori-testnet-v1
```  

Add peers in the config file:
```shell
`sed -i.bak "s/persistent_peers =.*/persistent_peers = "0dde2ae55624d822eeea57d1b5e1223b6019a531@176.9.149.15:26656,4d2ea61e6195ee4e449c1e6132cabce98f7d94e1@5.9.40.222:26656,bceb776975aab62bcfd501969c0e1a2734ed7c2e@176.9.19.162:26656"/" $HOME/.teritorid/config/config.toml`
```  

Launch the node:
```shell
teritorid start
```  

Wait for the chain to synchronize with the current block...  

## Setup your account  

Create an account:  
```shell 
teritorid keys add <YOUR_KEY>
 ```  
 
 You can also you `--recover` flag to use an already existed key (but we recommend for security reason to use one key per chain to avoid total loss of funds in case one key is missing)  
 
Create your validator:  
```shell 
teritorid tx staking create-validator \
 --commission-max-change-rate=0.01 \
 --commission-max-rate=0.2 \
 --commission-rate=0.05 \
 --amount 10000stake \
 --pubkey=$(teritorid tendermint show-validator) \
 --moniker=<YOUR_MONIKER> \
 --chain-id=teritori-testnet-v1 \
 --details="<DESCRIPTION_OF_YOUR_VALIDATOR>" \
 --security-contact="<YOUR_EMAIL_ADDRESS" \
 --website="<YOUR_WEBSITE>" \
 --identity="<YOUR_KEYBASE_ID>" \
 --min-self-delegation=10000 \
 --from=<YOUR_KEY>
 ```  
