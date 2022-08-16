package keeper

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllNftStakings(ctx sdk.Context) []types.NftStaking {
	stakings := []types.NftStaking{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixKeyNftStaking)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		staking := types.NftStaking{}
		k.cdc.MustUnmarshal(iterator.Value(), &staking)
		stakings = append(stakings, staking)
	}

	return stakings
}

func (k Keeper) GetNftStakingsByOwner(ctx sdk.Context, addr string) []types.NftStaking {
	stakings := []types.NftStaking{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PrefixKeyNftStakingByOwner, []byte(addr)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		identifier := string(iterator.Value())
		stakings = append(stakings, k.GetNftStaking(ctx, identifier))
	}

	return stakings
}

func (k Keeper) GetNftStaking(ctx sdk.Context, identifier string) types.NftStaking {
	staking := types.NftStaking{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NftStakingKey(identifier))
	if bz == nil {
		return staking
	}
	k.cdc.MustUnmarshal(bz, &staking)
	return staking
}

func (k Keeper) SetNftStaking(ctx sdk.Context, staking types.NftStaking) {
	nft := k.GetNftStaking(ctx, staking.NftIdentifier)
	if nft.NftIdentifier == "" {
		k.DeleteNftStaking(ctx, nft)
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&staking)
	store.Set(types.NftStakingKey(staking.NftIdentifier), bz)
	store.Set(types.NftStakingByOwnerKey(staking.RewardAddress, staking.NftIdentifier), []byte(staking.NftIdentifier))
}

func (k Keeper) DeleteNftStaking(ctx sdk.Context, staking types.NftStaking) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.NftStakingKey(staking.NftIdentifier))
	store.Delete(types.NftStakingByOwnerKey(staking.RewardAddress, staking.NftIdentifier))
}

func (k Keeper) AllocateTokensToRewardAddress(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	err := k.bk.MintCoins(ctx, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{amount})
}
