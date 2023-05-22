package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgBurnTokens{}

var MsgTypeBurnTokens = "burn_tokens"

func NewMsgBurnTokens(
	sender string,
	amount sdk.Coins,
) *MsgBurnTokens {
	return &MsgBurnTokens{
		Sender: sender,
		Amount: amount,
	}
}

func (m *MsgBurnTokens) Route() string {
	return ModuleName
}

func (m *MsgBurnTokens) Type() string {
	return MsgTypeBurnTokens
}

func (m *MsgBurnTokens) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyAddress
	}

	return nil
}

func (m *MsgBurnTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgBurnTokens) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}
