package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	MsgTypeRegisterNftStaking = "register_nft_staking"
	MsgTypeSetAccessInfo      = "set_access_info"
)

var _ sdk.Msg = &MsgRegisterNftStaking{}

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
	return MsgTypeRegisterNftStaking
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

var _ sdk.Msg = &MsgSetAccessInfo{}

func NewMsgSetAccessInfo(
	sender string,
	accessInfo Access,
) *MsgSetAccessInfo {
	return &MsgSetAccessInfo{
		Sender:     sender,
		AccessInfo: accessInfo,
	}
}

func (m *MsgSetAccessInfo) Route() string {
	return ModuleName
}

func (m *MsgSetAccessInfo) Type() string {
	return MsgTypeSetAccessInfo
}

func (m *MsgSetAccessInfo) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptySender
	}

	return nil
}

func (m *MsgSetAccessInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetAccessInfo) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{
		sender,
	}
}
