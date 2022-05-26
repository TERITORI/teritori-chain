package keeper

import (
	"github.com/NXTPOP/teritori-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	bk       types.BankKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "upop"
}
