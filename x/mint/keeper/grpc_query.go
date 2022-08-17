package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TERITORI/teritori-chain/x/mint/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/mint keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Params returns params of the mint module.
func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// BlockProvisions returns minter.BlockProvisions of the mint module.
func (q Querier) BlockProvisions(c context.Context, _ *types.QueryBlockProvisionsRequest) (*types.QueryBlockProvisionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := q.Keeper.GetMinter(ctx)

	return &types.QueryBlockProvisionsResponse{BlockProvisions: minter.BlockProvisions}, nil
}
