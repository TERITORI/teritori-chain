package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}

	_ codectypes.UnpackInterfacesMessage = MsgSubmitTx{}
)

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(owner, connectionID string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionID,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgSubmitTx creates and returns a new MsgSubmitTx instance
func NewMsgSubmitTx(sdkMsg sdk.Msg, connectionID, owner string) (*MsgSubmitTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitTx{
		ConnectionId: connectionID,
		Owner:        owner,
		Msgs:         []*codectypes.Any{any},
	}, nil
}

// PackTxMsgAny marshals the sdk.Msg payload to a protobuf Any type
func PackTxMsgAny(sdkMsg sdk.Msg) (*codectypes.Any, error) {
	msg, ok := sdkMsg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", sdkMsg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return any, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, any := range msg.Msgs {
		var msg sdk.Msg
		err := unpacker.UnpackAny(any, &msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitTx) GetTxMsg() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, any := range msg.Msgs {
		sdkMsg, ok := any.GetCachedValue().(sdk.Msg)
		if sdkMsg != nil {
			sdkMsgs = append(sdkMsgs, sdkMsg)
		}
		if !ok {
			return nil
		}
	}

	return sdkMsgs
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}
