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

	params := k.keeper.GetParamSet(ctx)
	if msg.Sender != params.Owner {
		return nil, types.ErrNotEnoughPermission
	}
	k.keeper.SetAllocation(ctx, msg.Allocation)
	return &types.MsgSetAllocationResponse{}, nil
}

func (m msgServer) TransferModuleOwnership(goCtx context.Context, msg *types.MsgTransferModuleOwnership) (*types.MsgTransferModuleOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := m.keeper.GetParamSet(ctx)
	if msg.Sender != params.Owner {
		return nil, types.ErrNotEnoughPermission
	}
	params.Owner = msg.NewOwner
	m.keeper.SetParamSet(ctx, params)

	return &types.MsgTransferModuleOwnershipResponse{}, nil
}

func (m msgServer) DepositTokens(goCtx context.Context, msg *types.MsgDepositTokens) (*types.MsgDepositTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = m.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositTokensResponse{}, nil
}
