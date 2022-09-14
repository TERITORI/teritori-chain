package keeper

import (
	"fmt"

	"github.com/TERITORI/teritori-chain/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParams(ctx)
	blockNumber := ctx.BlockHeight()

	// not distribute rewards if it's not time yet for rewards distribution
	if blockNumber < params.MintingRewardsDistributionStartBlock {
		return
	} else if blockNumber == params.MintingRewardsDistributionStartBlock {
		k.SetLastReductionBlockNum(ctx, blockNumber)
	}
	// fetch stored minter & params
	minter := k.GetMinter(ctx)

	// Check if we have hit an block where we update the inflation parameter.
	// We measure time between reductions in number of blocks.
	// This avoids issues with measuring in block numbers, as blocks have fixed intervals, with very
	// low variance at the relevant sizes. As a result, it is safe to store the block number
	// of the last reduction to be later retrieved for comparison.
	if blockNumber >= params.ReductionPeriodInBlocks+k.GetLastReductionBlockNum(ctx) {
		// Reduce the reward per reduction period
		minter.BlockProvisions = minter.NextBlockProvisions(params)
		k.SetMinter(ctx, minter)
		k.SetLastReductionBlockNum(ctx, blockNumber)
	}

	// implement automatic monthInfo updates
	monthInfo := k.GetTeamVestingMonthInfo(ctx)
	if blockNumber >= monthInfo.OneMonthPeriodInBlocks+monthInfo.MonthStartedBlock {
		monthInfo.MonthsSinceGenesis++
		monthInfo.MonthStartedBlock = ctx.BlockHeight()
		k.SetTeamVestingMonthInfo(ctx, monthInfo)
	}

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	// We over-allocate by the developer vesting portion, and burn this later
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.DistributeMintedCoin(ctx, mintedCoin)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ModuleName,
			sdk.NewAttribute(types.AttributeBlockNumber, fmt.Sprintf("%d", blockNumber)),
			sdk.NewAttribute(types.AttributeKeyBlockProvisions, minter.BlockProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
