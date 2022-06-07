package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/NXTPOP/teritori-chain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper struct
type Keeper struct {
	cdc           codec.Codec
	storeKey      sdk.StoreKey
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	acountKeeper  types.AccountKeeper
}

// NewKeeper returns keeper
func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, bk types.BankKeeper, sk types.StakingKeeper, ak types.AccountKeeper) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		bankKeeper:    bk,
		stakingKeeper: sk,
		acountKeeper:  ak,
	}
}

// Logger returns logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
