package keeper

import (
	"context"

	"github.com/TERITORI/teritori-chain/x/airdrop/types"
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

	err := k.keeper.ClaimAllocation(ctx, msg.Address, msg.PubKey, msg.RewardAddress, msg.Signature)
	return &types.MsgClaimAllocationResponse{}, err
}

func (k msgServer) SetAllocation(goCtx context.Context, msg *types.MsgSetAllocation) (*types.MsgSetAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: check permission
	k.keeper.SetAllocation(ctx, msg.Allocation)
	return &types.MsgSetAllocationResponse{}, nil
}
