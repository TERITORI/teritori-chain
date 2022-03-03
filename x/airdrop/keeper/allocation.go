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

func (k Keeper) ClaimAllocation(ctx sdk.Context, address string, rewardAddress string, signature []byte) error {
	// ensure allocation exists for the address
	allocation := k.GetAllocation(ctx, address)
	if allocation == nil {
		return types.ErrAirdropAllocationDoesNotExists
	}

	// ensure allocation is not claimed already
	unclaimed := allocation.Amount.Sub(allocation.ClaimedAmount)
	if unclaimed.IsZero() {
		return types.ErrAirdropAllocationAlreadyClaimed
	}

	// verify native chain account with signature
	sigOk := verifySignature(allocation.Chain, allocation.Address, rewardAddress, signature)
	if !sigOk {
		return types.ErrNativeChainAccountSigVerificationFailure
	}

	// send coins from airdrop module account to beneficiary address
	sdkAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdkAddr, sdk.Coins{sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unclaimed)})
	if err != nil {
		return err
	}

	// update claimed amount and set the record on-chain
	allocation.ClaimedAmount = allocation.Amount
	k.SetAllocation(ctx, *allocation)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimAllocation,
			sdk.NewAttribute(types.AttributeKeyAddress, address),
			sdk.NewAttribute(types.AttributeKeyAmount, unclaimed.String()),
			sdk.NewAttribute(types.AttributeKeyRewardAddress, rewardAddress),
		),
	)

	return nil
}
