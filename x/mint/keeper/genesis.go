package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TERITORI/teritori-chain/x/mint/types"
)

// InitGenesis new mint genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	if data == nil {
		panic("nil mint genesis state")
	}

	data.Minter.BlockProvisions = data.Params.GenesisBlockProvisions
	k.SetMinter(ctx, data.Minter)
	k.SetParams(ctx, data.Params)

	// The call to GetModuleAccount creates a module account if it does not exist.
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	k.SetLastReductionBlockNum(ctx, data.ReductionStartedBlock)
	k.SetTeamVestingMonthInfo(ctx, data.MonthInfo)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	if params.WeightedDeveloperRewardsReceivers == nil {
		params.WeightedDeveloperRewardsReceivers = make([]types.MonthlyVestingAddress, 0)
	}

	lastReductionBlock := k.GetLastReductionBlockNum(ctx)
	monthInfo := k.GetTeamVestingMonthInfo(ctx)
	return types.NewGenesisState(minter, params, lastReductionBlock, monthInfo)
}
