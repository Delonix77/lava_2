package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	legacyerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRelayPayment = "relay_payment"

var _ sdk.Msg = &MsgRelayPayment{}

func NewMsgRelayPayment(creator string, relays []*RelaySession, description string) *MsgRelayPayment {
	return &MsgRelayPayment{
		Creator:           creator,
		Relays:            relays,
		DescriptionString: description,
	}
}

func (msg *MsgRelayPayment) Route() string {
	return RouterKey
}

func (msg *MsgRelayPayment) Type() string {
	return TypeMsgRelayPayment
}

func (msg *MsgRelayPayment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRelayPayment) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRelayPayment) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(legacyerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
