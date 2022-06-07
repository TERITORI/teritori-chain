package keeper

import (
	"github.com/NXTPOP/teritori-chain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get params
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	return types.Params{}, nil
}

// SetParams set params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return nil
}
