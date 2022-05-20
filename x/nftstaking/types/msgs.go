package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgRegisterNftStaking{}

var MsgTypeClaimAllocation = "claim_allocation"

func NewMsgRegisterNftStaking(
	sender string,
	nftStaking NftStaking,
) *MsgRegisterNftStaking {
	return &MsgRegisterNftStaking{
		Sender:     sender,
		NftStaking: nftStaking,
	}
}

func (m *MsgRegisterNftStaking) Route() string {
	return ModuleName
}

func (m *MsgRegisterNftStaking) Type() string {
	return MsgTypeClaimAllocation
}

func (m *MsgRegisterNftStaking) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptySender
	}

	if m.NftStaking.RewardAddress == "" {
		return ErrEmptyRewardAddress
	}

	return nil
}

func (m *MsgRegisterNftStaking) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterNftStaking) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{
		sender,
	}
}
