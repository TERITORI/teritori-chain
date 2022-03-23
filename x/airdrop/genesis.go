package airdrop

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/keeper"
	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetAllocation(ctx, *types.DefaultGenesis().Allocations[0])
	k.SetAllocation(ctx, *types.DefaultGenesis().Allocations[1])
	k.SetAllocation(ctx, *types.DefaultGenesis().Allocations[2])
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
