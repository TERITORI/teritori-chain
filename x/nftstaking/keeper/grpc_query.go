package keeper

import (
	"context"

	"github.com/NXTPOP/teritori-chain/x/nftstaking/types"
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
