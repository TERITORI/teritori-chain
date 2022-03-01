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

	// TODO: check allocation exists for the address
	// TODO: ensure allocation is not claimed already (Allocation.Amount - Allocation.ClaimedAmount) > 0
	// TODO: ensure signature is correct based on chain and reward address
	// TODO: send tokens from airdrop module account to individual account
	// TODO: update claimed amount and set the record on-chain
	// TODO: emit event for native-chain-address, amount, reward_address

	return &types.MsgClaimAllocationResponse{}, nil
}
