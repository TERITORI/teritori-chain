package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NXTPOP/teritori-chain/x/airdrop/keeper"
	"github.com/NXTPOP/teritori-chain/x/airdrop/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, allocation := range genState.Allocations {
		k.SetAllocation(ctx, *allocation)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
