package keeper

import (
	"github.com/NXTPOP/teritori-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllNftTypePerms(ctx sdk.Context) []types.NftTypePerms {
	allPerms := []types.NftTypePerms{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixKeyNftTypePerms)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		perms := types.NftTypePerms{}
		k.cdc.MustUnmarshal(iterator.Value(), &perms)
		allPerms = append(allPerms, perms)
	}

	return allPerms
}

func (k Keeper) GetNftTypePerms(ctx sdk.Context, nftType types.NftType) types.NftTypePerms {
	perms := types.NftTypePerms{}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyNftTypePerms)
	bz := prefixStore.Get(sdk.Uint64ToBigEndian(uint64(perms.NftType)))
	if bz == nil {
		return perms
	}
	k.cdc.MustUnmarshal(bz, &perms)

	return perms
}

func (k Keeper) SetNftTypePerms(ctx sdk.Context, perms types.NftTypePerms) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyNftTypePerms)
	bz := k.cdc.MustMarshal(&perms)
	prefixStore.Set(sdk.Uint64ToBigEndian(uint64(perms.NftType)), bz)
}
