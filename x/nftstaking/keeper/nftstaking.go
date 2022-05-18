package keeper

import (
	"github.com/POPSmartContract/nxtpop-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k Keeper) SetNftStaking(ctx sdk.Context, staking types.NftStaking) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyNftStaking)
	bz := k.cdc.MustMarshal(&staking)
	prefixStore.Set([]byte(staking.NftIdentifier), bz)
}

func (k Keeper) AllocateTokensToRewardAddress(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	err := k.bk.MintCoins(ctx, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{amount})
}
