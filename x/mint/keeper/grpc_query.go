package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// Inflation returns minter.Inflation of the mint module.
func (q Querier) Inflation(c context.Context, _ *types.QueryInflationRequest) (*types.QueryInflationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := q.Keeper.GetMinter(ctx)
	params := q.Keeper.GetParams(ctx)
	mintDenomSupply := q.bankKeeper.GetSupply(ctx, params.MintDenom).Amount

	inflation := minter.BlockProvisions.
		Mul(math.LegacyNewDec(int64(params.BlocksPerYear))).
		Quo(math.LegacyNewDecFromInt(mintDenomSupply))

	return &types.QueryInflationResponse{Inflation: inflation}, nil
}

// StakingAPR returns the current staking APR value.
func (q Querier) StakingAPR(c context.Context, _ *types.QueryStakingAPRRequest) (*types.QueryStakingAPRResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := q.Keeper.GetMinter(ctx)
	params := q.Keeper.GetParams(ctx)
	totalStaked := q.stakingKeeper.TotalBondedTokens(ctx)

	stakingApr := minter.BlockProvisions.
		Mul(math.LegacyNewDec(int64(params.BlocksPerYear))).
		Mul(params.DistributionProportions.Staking).
		Quo(math.LegacyNewDecFromInt(totalStaked))

	return &types.QueryStakingAPRResponse{Apr: stakingApr}, nil
}
