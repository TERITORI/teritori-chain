package keeper

import (
	"context"

	"github.com/TERITORI/teritori-chain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Allocation(c context.Context, req *types.QueryAllocationRequest) (*types.QueryAllocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAllocationResponse{
		Allocation: k.GetAllocation(ctx, req.Address),
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{
		Params: k.GetParamSet(ctx),
	}, nil
}
