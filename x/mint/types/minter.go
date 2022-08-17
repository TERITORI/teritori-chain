package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	errNilBlockProvisions      = errors.New("block provisions was nil in genesis")
	errNegativeBlockProvisions = errors.New("block provisions should be non-negative")
)

// NewMinter returns a new Minter object with the given block
// provisions values.
func NewMinter(blockProvisions sdk.Dec) Minter {
	return Minter{
		BlockProvisions: blockProvisions,
	}
}

// InitialMinter returns an initial Minter object.
func InitialMinter() Minter {
	return NewMinter(sdk.NewDec(0))
}

// DefaultInitialMinter returns a default initial Minter object for a new chain.
func DefaultInitialMinter() Minter {
	return InitialMinter()
}

// Validate validates minter. Returns nil on success, error otherewise.
func (m Minter) Validate() error {
	if m.BlockProvisions.IsNil() {
		return errNilBlockProvisions
	}

	if m.BlockProvisions.IsNegative() {
		return errNegativeBlockProvisions
	}
	return nil
}

// NextBlockProvisions returns the block provisions.
func (m Minter) NextBlockProvisions(params Params) sdk.Dec {
	return m.BlockProvisions.Mul(params.ReductionFactor)
}

// BlockProvision returns the provisions for a block based on the block
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.BlockProvisions
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
