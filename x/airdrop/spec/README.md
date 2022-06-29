# `x/airdrop`

## Abstract

`airdrop` module is to allocate airdrop to different network addresses.

The target networks available to allocate airdrops are

- solana
- evm
- cosmos
- osmosis
- juno
- terra

Airdrop allocation can be set from genesis or admin can set the allocation for different network addresses.
The user with airdrop allocation can send address ownership verification signature to receive airdrop.

## State

Airdrop module keeps the information of `AirdropAllocation` that shows allocation and claimed amounts for an address on different network.

```go
type AirdropAllocation struct {
	Chain         string
	Address       string
	Amount        sdk.Coin
	ClaimedAmount sdk.Coin
}
```

## Messages

### MsgSetAllocation

`MsgSetAllocation` describes the message to set allocation for an account by admin.

```go
type MsgSetAllocation struct {
	Sender     string
	Allocation AirdropAllocation
}
```

### MsgClaimAllocation

`MsgClaimAllocation` describes the message to claim airdrop allocation allocated to different network address.
RewardAddress is the teritori chain address that receives allocation.

```go
type MsgClaimAllocation struct {
	Address       string
	PubKey        string
	RewardAddress string
	Signature     string
}
```
