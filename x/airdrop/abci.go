package airdrop

import (
	"github.com/TERITORI/teritori-chain/x/airdrop/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
}
