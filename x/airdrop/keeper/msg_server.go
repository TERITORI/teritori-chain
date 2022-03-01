package keeper

import (
	"context"

	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/types"
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

func (k msgServer) ClaimAllocation(goCtx context.Context, msg *types.MsgClaimAllocation) (*types.MsgClaimAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgClaimAllocationResponse{}, nil
}
