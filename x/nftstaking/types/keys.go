package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// constants
var (
	ModuleName = "nftstaking"

	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixKeyNftStaking        = []byte{0x0}
	PrefixKeyNftStakingByOwner = []byte{0x1}
	PrefixKeyAccessInfo        = []byte{0x2}
	PrefixKeyNftTypePerms      = []byte{0x3}
)

func NftStakingKey(identifier string) []byte {
	return append(PrefixKeyNftStaking, identifier...)
}

func NftStakingByOwnerKey(owner string, identifier string) []byte {
	return append(append(PrefixKeyNftStakingByOwner, owner...), identifier...)
}

func AccessInfoKey(addr sdk.AccAddress) []byte {
	return append(PrefixKeyAccessInfo, addr...)
}

func NftTypePermsKey(nftType NftType) []byte {
	return append(PrefixKeyNftTypePerms, sdk.Uint64ToBigEndian(uint64(nftType))...)
}
