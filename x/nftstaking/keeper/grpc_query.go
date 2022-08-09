package keeper

import (
	"context"

	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) QueryNftStakings(goCtx context.Context, request *types.QueryNftStakingsRequest) (*types.QueryNftStakingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryNftStakingsResponse{
		Nftstakings: q.keeper.GetAllNftStakings(ctx),
	}, nil
}

func (q Querier) QueryAccessInfos(goCtx context.Context, request *types.QueryAccessInfosRequest) (*types.QueryAccessInfosResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAccessInfosResponse{
		AccessInfos: q.keeper.GetAllAccessInfos(ctx),
	}, nil
}

func (q Querier) QueryAccessInfo(goCtx context.Context, request *types.QueryAccessInfoRequest) (*types.QueryAccessInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAccessInfoResponse{
		AccessInfo: q.keeper.GetAccessInfo(ctx, request.Address),
	}, nil
}
