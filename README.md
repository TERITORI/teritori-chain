# Teritori chain

Teritori chain supports cosmwasm, airdrop, and nftstaking.

## How to setup binary

```
make install
```

## How to run test

```
make test
```

## Testnet setup

### Bootstrap genesis and configurations

```
teritorid testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2 --chain-id=teritori-testnet-v2 --keyring-backend=test
```

### Replace configuration files

#### Update `config.toml` files of testnet

```
192.168.10.2 => 176.9.19.162
192.168.10.3 => 176.9.149.15
192.168.10.4 => 5.9.40.222
192.168.10.5 => 78.46.106.69
```

```
prometheus = true
cors_allowed_origins = ["*"]
```

#### Update `app.toml` files of testnet

```
enabled-unsafe-cors = true
```

### Update all "stake" denom on genesis to "utori"

### Update public genesis

Update `<repo>/genesis/genesis.json` to latest one generated

### Create zip file for output and setup 4 nodes

```sh
ssh root@176.9.19.162
ssh root@176.9.149.15
ssh root@5.9.40.222
ssh root@78.46.106.69

# copy zip via scp into servers
scp ./testnet.zip root@176.9.19.162:~/testnet.zip
scp ./testnet.zip root@176.9.149.15:~/testnet.zip
scp ./testnet.zip root@5.9.40.222:~/testnet.zip
scp ./testnet.zip root@78.46.106.69:~/testnet.zip

# kill previous processes if required
ps -e | grep nxtpopd
kill -9 78428
rm -rf nodehome/
rm -rf testnet
rm -rf testnet_bkup2.zip
rm -rf nxtpop-chain/

sudo apt-get install unzip
unzip testnet.zip

rm -rf /usr/local/go
curl -LO https://go.dev/dl/go1.18.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
sudo nano ~/.profile

export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
source ~/.profile
go version

sudo apt install build-essential
sudo apt install git-all

eval `ssh-agent -s`
ssh-keygen
ssh-add
cat .ssh/id_rsa.pub
git clone git@github.com:NXTPOP/teritori-chain.git

cd teritori-chain/
go install ./cmd/teritorid/

mv testnet/node0/teritorid/ nodehome/
mv testnet/node1/teritorid/ nodehome/
mv testnet/node2/teritorid/ nodehome/
mv testnet/node3/teritorid/ nodehome/

teritorid start --home=nodehome &

```

### Setup 5th node following normal validator joining process

```
ssh root@178.63.25.244
```

Update `chain_id` and `persistent_peers` field to updated testnet version on the document.
https://docs.google.com/document/d/1wa9XOnvKSIBe9sPukV2xEgkgfvbsaboyI-NWnX-ERYU/edit?usp=sharing
