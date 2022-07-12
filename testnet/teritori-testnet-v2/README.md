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
git clone https://github.com/TERITORI/teritori-chain && cd teritori-chain && git checkout teritori-testnet-v2 && make install
```  

Verify the installation:  
```shell
teritorid version
#Should return  teritori-testnet-v2-0f4e5cb1d529fa18971664891a9e8e4c114456c6
```  

Init the chain:
```shell
teritorid init <YOUR_MONIKER> --chain-id teritori-testnet-v2
```  

Add peers in the config file:
```shell
`sed -i.bak "s/persistent_peers =.*/persistent_peers = "0b42fd287d3bb0a20230e30d54b4b8facc412c53@176.9.149.15:26656,2371b28f366a61637ac76c2577264f79f0965447@176.9.19.162:26656,2f394edda96be07bf92b0b503d8be13d1b9cc39f@5.9.40.222:26656"/" $HOME/.teritorid/config/config.toml`
```  

Launch the node:
```shell
teritorid start
```  

Wait for the chain to synchronize with the current block...  

Create your validator:  
```shell 
teritorid tx staking create-validator \
 --commission-max-change-rate=0.01 \
 --commission-max-rate=0.2 \
 --commission-rate=0.05 \
 --amount 1000000utori \
 --pubkey=$(teritorid tendermint show-validator) \
 --moniker=<YOUR_MONIKER> \
 --chain-id=teritori-testnet-v2 \
 --details="<DESCRIPTION_OF_YOUR_VALIDATOR>" \
 --security-contact="<YOUR_EMAIL_ADDRESS" \
 --website="<YOUR_WEBSITE>" \
 --identity="<YOUR_KEYBASE_ID>" \
 --min-self-delegation=1000000 \
 --from=<YOUR_KEY>
 ```  
