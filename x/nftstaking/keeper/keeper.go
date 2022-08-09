package keeper

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramSpace paramstypes.Subspace
	bk         types.BankKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	cdc codec.BinaryCodec,
	bk types.BankKeeper) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		paramSpace: paramSpace,
		storeKey:   storeKey,
		bk:         bk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "utori"
}
