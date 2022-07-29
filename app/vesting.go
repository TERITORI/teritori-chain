package teritori

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// VestingTransactionDecorator prevents staking from vesting accounts
type VestingTransactionDecorator struct {
	ak ante.AccountKeeper
}

func NewVestingTransactionDecorator(ak ante.AccountKeeper) VestingTransactionDecorator {
	return VestingTransactionDecorator{
		ak: ak,
	}
}

// AnteHandle prevents staking from vesting accounts
func (vtd VestingTransactionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
		if !ok {
			continue
		}

		sender := sdk.MustAccAddressFromBech32(delegateMsg.DelegatorAddress)
		acc := vtd.ak.GetAccount(ctx, sender)
		if acc == nil {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
				"account %s does not exist", acc)
		}

		// Check if vesting account
		_, isVesting := acc.(*vestingtypes.BaseVestingAccount)
		if !isVesting {
			return next(ctx, tx, simulate)
		}

		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress,
			"cannot perform delegation on vesting account: %s", delegateMsg.DelegatorAddress,
		)
	}

	return next(ctx, tx, simulate)
}
