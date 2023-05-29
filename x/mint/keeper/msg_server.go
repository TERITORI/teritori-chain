package keeper

import (
	"context"

	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates and returns a new types.MsgServer, fulfilling the intertx Msg service interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterAccount implements the Msg/RegisterAccount interface
func (k msgServer) BurnTokens(goCtx context.Context, msg *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, govtypes.ModuleName, sdk.Coins(msg.Amount))
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.BurnCoins(ctx, govtypes.ModuleName, sdk.Coins(msg.Amount))
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnTokensResponse{}, nil
}
