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
