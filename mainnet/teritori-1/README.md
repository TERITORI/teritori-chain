# [COMING]

## Server Configuration

Server requirements:

- No. of CPUs: 16
- Memory: 32GB
- Disk: 500GB
- OS: Ubuntu 22.04 LTS

Allow all incoming connections from TCP port 26656 and 26657

Notes on the configurations.

1. Multicore is important, regardless the less CPU time usage
2. teritorid uses less than 1GB memory and 2GB should be enough for now.
   Once your new server is running, login to the server and upgrade your packages.

## Setup your machine

If you already have go 1.18+ and packages up to date, you can skip this part and jump to the second section: [Setup the chain](#setup-the-chain)  
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

Setup your environnement (you can skip this part if you already had go installed before):

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

Clone the Teritori repository and install mainnet binary:

```shell
git clone https://github.com/TERITORI/teritori-chain && cd teritori-chain && git checkout v1.1.1 && make install
```

Verify the installation:

```shell
teritorid version
#Should return  mainnet-9bc209964838de9f93be332e07a842bcf1297fbe
```

Init the chain:

```shell
teritorid init <YOUR_MONIKER> --chain-id teritori-1
```

Add peers in the config file:

```shell
sed -i.bak 's/persistent_peers =.*/persistent_peers = ""/' $HOME/.teritorid/config/config.toml
```

Download the genesis file:

```shell
wget -O ~/.teritorid/config/genesis.json https://media.githubusercontent.com/media/TERITORI/teritori-chain/v1.1.1/mainnet/teritori-1/genesis.json
```

## Launch the node

You have multiple choice for launching the chain, choose the one that you prefer in the below list:

- [Manual](https://github.com/TERITORI/teritori-chain/tree/main/mainnet/teritori-1#Manual)
- [Systemctl](https://github.com/TERITORI/teritori-chain/tree/main/mainnet/teritori-1#Systemctl)
- [Cosmovisor](https://github.com/TERITORI/teritori-chain/tree/main/mainnet/teritori-1#Cosmovisor)

### **Manual**

- Create a screen and setup limit of files
- Launch the chain

```shell
screen -S teritori
ulimit -n 4096
teritorid start
```

You can escape the screen pressing `CTRL + AD` and enter it again using:

```shell
screen -r teritori
```

### **Systemctl**

- Create service file
- Enable and launch the service file
- Setup the logs

```shell
tee <<EOF >/dev/null /etc/systemd/system/teritorid.service
[Unit]
Description=Teritori Cosmos daemon
After=network-online.target

[Service]
User=$USER
ExecStart=/home/$USER/go/bin/teritorid start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```

```shell
systemctl enable teritorid
systemctl daemon-reload
systemctl restart teritorid
```

To check the logs:

```shell
journalctl -u teritorid.service -f -n 100
```

### **Cosmovisor**

- Install cosmovisor
- Setup environment variables
- Create folders needed for upgrades
- Copy binaries to cosmovisor bin
- Create service file
- Enable and launch the service file
- Setup the logs

```shell
go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest
```

```shell
export DAEMON_NAME=teritorid
export DAEMON_HOME=$HOME/.teritori
source ~/.profile
```

```shell
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades
```

```shell
cp $HOME/go/bin/teritorid $DAEMON_HOME/cosmovisor/genesis/bin
```

```shell
tee <<EOF >/dev/null /etc/systemd/system/teritorid.service
[Unit]
Description=Teritori Daemon (cosmovisor)

After=network-online.target

[Service]
User=$USER
ExecStart=/home/$USER/go/bin/cosmovisor start
Restart=always
RestartSec=3
LimitNOFILE=4096
Environment="DAEMON_NAME=teritorid"
Environment="DAEMON_HOME=/home/$USER/.teritori"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_LOG_BUFFER_SIZE=512"

[Install]
WantedBy=multi-user.target
EOF
```

```shell
systemctl enable teritorid
systemctl daemon-reload
systemctl restart teritorid
```

To check the logs:

```shell
journalctl -u teritorid.service -f -n 100
```

Wait for the chain to synchronize with the current block... you can do the next step during this time

## Setup validator for genesis validators

Copy `node_key.json` and `priv_validator_key.json` into home directory.

## Setup validator for non-genesis validators

Create an account:

```shell
teritorid keys add <YOUR_KEY>
```

You can also you `--recover` flag to use an already existed key (but we recommend for security reason to use one key per chain to avoid total loss of funds in case one key is missing)

Join our [Discord](https://discord.gg/teritori) and request fund on the `Faucet` channel using this command:

```shell
$request <YOUR_TERITORI_ADDRESS>
```

You can check if you have received fund once your node will be synched using this CLI command:

```shell
teritorid query bank balances <YOUR_TERITORI_ADDRESS> --chain-id teritori-1
```

Once the fund are received and chain is synchronized you can create your validator:

```shell
teritorid tx staking create-validator \
 --commission-max-change-rate=0.01 \
 --commission-max-rate=0.2 \
 --commission-rate=0.05 \
 --amount 1000000utori \
 --pubkey=$(teritorid tendermint show-validator) \
 --moniker=<YOUR_MONIKER> \
 --chain-id=teritori-1 \
 --details="<DESCRIPTION_OF_YOUR_VALIDATOR>" \
 --security-contact="<YOUR_EMAIL_ADDRESS" \
 --website="<YOUR_WEBSITE>" \
 --identity="<YOUR_KEYBASE_ID>" \
 --min-self-delegation=1000000 \
 --from=<YOUR_KEY>
```

Join us on [Discord](https://discord.gg/teritori) for mainnet discussions.
