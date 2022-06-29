# `x/nftstaking`

## Abstract

`nftstaking` module is provide rewards to nft stakers on different networks.

The owner registers nft stakers on the module.

Part of inflation rewards are allocated to registered nft stakers and broadcasted on everyblock.

## State

`nftstaking` module keeps the information of `NftStaking` that shows nft information and reward address with weights.

```go
type NftStaking struct {
	NftIdentifier string
	NftMetadata   string
	RewardAddress string
	RewardWeight  uint64
}
```

## Messages

### MsgRegisterNftStaking

`MsgRegisterNftStaking` describes the message to register staked nfts by admin.

```go
type MsgRegisterNftStaking struct {
	Sender     string
	NftStaking NftStaking
}
```
