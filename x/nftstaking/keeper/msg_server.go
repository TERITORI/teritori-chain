package keeper

import (
	"context"

	"github.com/NXTPOP/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterNftStaking(goCtx context.Context, msg *types.MsgRegisterNftStaking) (*types.MsgRegisterNftStakingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.keeper.SetNftStaking(ctx, msg.NftStaking)
	return &types.MsgRegisterNftStakingResponse{}, nil
}

func (k msgServer) SetAccessInfo(goCtx context.Context, msg *types.MsgSetAccessInfo) (*types.MsgSetAccessInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.keeper.SetAccessInfo(ctx, msg.AccessInfo)
	return &types.MsgSetAccessInfoResponse{}, nil
}
