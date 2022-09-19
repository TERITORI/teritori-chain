package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgClaimAllocation{}

var MsgTypeClaimAllocation = "claim_allocation"

func NewMsgClaimAllocation(
	address string,
	rewardAddress sdk.AccAddress,
	signature string,
) *MsgClaimAllocation {
	return &MsgClaimAllocation{
		Address:       address,
		RewardAddress: rewardAddress.String(),
		Signature:     signature,
	}
}

func (m *MsgClaimAllocation) Route() string {
	return ModuleName
}

func (m *MsgClaimAllocation) Type() string {
	return MsgTypeClaimAllocation
}

func (m *MsgClaimAllocation) ValidateBasic() error {
	if m.RewardAddress == "" {
		return ErrEmptyRewardAddress
	}

	if m.Address == "" {
		return ErrEmptyOnChainAllocationAddress
	}

	return nil
}

func (m *MsgClaimAllocation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimAllocation) GetSigners() []sdk.AccAddress {
	rewardAddr, err := sdk.AccAddressFromBech32(m.RewardAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		rewardAddr,
	}
}

var _ sdk.Msg = &MsgSetAllocation{}

var MsgTypeSetAllocation = "set_allocation"

func NewMsgSetAllocation(
	sender string,
	allocation AirdropAllocation,
) *MsgSetAllocation {
	return &MsgSetAllocation{
		Sender:     sender,
		Allocation: allocation,
	}
}

func (m *MsgSetAllocation) Route() string {
	return ModuleName
}

func (m *MsgSetAllocation) Type() string {
	return MsgTypeSetAllocation
}

func (m *MsgSetAllocation) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyAddress
	}

	return nil
}

func (m *MsgSetAllocation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetAllocation) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}

var _ sdk.Msg = &MsgTransferModuleOwnership{}

var MsgTypeTransferModuleOwnership = "transfer_module_ownership"

func NewMsgTransferModuleOwnership(
	sender sdk.AccAddress,
	newOwner string,
) *MsgTransferModuleOwnership {
	return &MsgTransferModuleOwnership{
		Sender:   sender.String(),
		NewOwner: newOwner,
	}
}

func (m *MsgTransferModuleOwnership) Route() string {
	return ModuleName
}

func (m *MsgTransferModuleOwnership) Type() string {
	return MsgTypeSetAllocation
}

func (m *MsgTransferModuleOwnership) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyAddress
	}

	return nil
}

func (m *MsgTransferModuleOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgTransferModuleOwnership) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}

var _ sdk.Msg = &MsgSignData{}

var MsgTypeSignData = "sign_data"

func NewMsgSignData(
	signer string,
	data []byte,
) *MsgSignData {
	return &MsgSignData{
		Signer: signer,
		Data:   data,
	}
}

func (m *MsgSignData) Route() string {
	return ModuleName
}

func (m *MsgSignData) Type() string {
	return MsgTypeSignData
}

func (m *MsgSignData) ValidateBasic() error {
	return nil
}

func (m *MsgSignData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSignData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
