package keeper

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllAccessInfos(ctx sdk.Context) []types.Access {
	accessInfos := []types.Access{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixKeyAccessInfo)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		accessInfo := types.Access{}
		k.cdc.MustUnmarshal(iterator.Value(), &accessInfo)
		accessInfos = append(accessInfos, accessInfo)
	}

	return accessInfos
}

func (k Keeper) GetAccessInfo(ctx sdk.Context, address string) types.Access {
	accessInfo := types.Access{}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyAccessInfo)
	bz := prefixStore.Get([]byte(address))
	if bz == nil {
		return accessInfo
	}
	k.cdc.MustUnmarshal(bz, &accessInfo)

	return accessInfo
}

func (k Keeper) SetAccessInfo(ctx sdk.Context, accessInfo types.Access) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyAccessInfo)
	bz := k.cdc.MustMarshal(&accessInfo)
	prefixStore.Set([]byte(accessInfo.Address), bz)
}
