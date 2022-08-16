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

func (q Querier) QueryNftStaking(goCtx context.Context, request *types.QueryNftStakingRequest) (*types.QueryNftStakingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryNftStakingResponse{
		Nftstaking: q.keeper.GetNftStaking(ctx, request.Identifier),
	}, nil
}

func (q Querier) QueryNftStakingsByOwner(goCtx context.Context, request *types.QueryNftStakingsByOwnerRequest) (*types.QueryNftStakingsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryNftStakingsByOwnerResponse{
		Nftstakings: q.keeper.GetNftStakingsByOwner(ctx, request.Owner),
	}, nil
}

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

func (q Querier) QueryAllNftTypePerms(goCtx context.Context, request *types.QueryAllNftTypePermsRequest) (*types.QueryAllNftTypePermsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllNftTypePermsResponse{
		AllNftTypePerms: q.keeper.GetAllNftTypePerms(ctx),
	}, nil
}

func (q Querier) QueryNftTypePerms(goCtx context.Context, request *types.QueryNftTypePermsRequest) (*types.QueryNftTypePermsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryNftTypePermsResponse{
		NftTypePerms: q.keeper.GetNftTypePerms(ctx, request.NftType),
	}, nil
}

func (q Querier) QueryHasPermission(goCtx context.Context, request *types.QueryHasPermissionRequest) (*types.QueryHasPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryHasPermissionResponse{
		HasPermission: q.keeper.HasPermission(ctx, request.Address, types.Permission(types.Permission_value[request.Permission])),
	}, nil
}
