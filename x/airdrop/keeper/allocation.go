package keeper

import (
	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllocation(ctx sdk.Context, address string) *types.AirdropAllocation {
	allocation := types.AirdropAllocation{}

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAirdropAllocation)
	bz := prefixStore.Get([]byte(address))
	if bz == nil {
		return nil
	}
	k.cdc.MustUnmarshal(bz, &allocation)
	return &allocation
}

func (k Keeper) GetAllAllocations(ctx sdk.Context) []types.AirdropAllocation {
	allocations := []types.AirdropAllocation{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixAirdropAllocation)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		allocation := types.AirdropAllocation{}
		k.cdc.MustUnmarshal(iterator.Value(), &allocation)
		allocations = append(allocations, allocation)
	}

	return allocations
}

func (k Keeper) SetAllocation(ctx sdk.Context, allocation types.AirdropAllocation) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAirdropAllocation)
	bz := k.cdc.MustMarshal(&allocation)
	prefixStore.Set([]byte(allocation.Address), bz)
}

func (k Keeper) DeleteAllocation(ctx sdk.Context, address string) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAirdropAllocation)
	prefixStore.Delete([]byte(address))
}
