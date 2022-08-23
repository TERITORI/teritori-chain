package keeper

import (
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetLastReductionBlockNum returns last reduction block number.
func (k Keeper) GetLastReductionBlockNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastReductionBlockKey)
	if b == nil {
		return 0
	}

	return int64(sdk.BigEndianToUint64(b))
}

// SetLastReductionBlockNum set last reduction block number.
func (k Keeper) SetLastReductionBlockNum(ctx sdk.Context, blockNum int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastReductionBlockKey, sdk.Uint64ToBigEndian(uint64(blockNum)))
}

// GetTeamVestingMonthInfo returns month information for team vesting
func (k Keeper) GetTeamVestingMonthInfo(ctx sdk.Context) types.TeamVestingMonthInfo {
	store := ctx.KVStore(k.storeKey)

	monthInfo := types.TeamVestingMonthInfo{}
	bz := store.Get(types.TeamVestingMonthInfoKey)
	if bz == nil {
		return monthInfo
	}

	k.cdc.MustUnmarshal(bz, &monthInfo)
	return monthInfo
}

// SetLastReductionBlockNum set last reduction block number.
func (k Keeper) SetTeamVestingMonthInfo(ctx sdk.Context, monthInfo types.TeamVestingMonthInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&monthInfo)
	store.Set(types.TeamVestingMonthInfoKey, bz)
}

// get the minter.
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter.
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}
