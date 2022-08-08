package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyOwner = []byte("Owner")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyOwner, &p.Owner, validateOwner),
	}
}

// NewParams constructs a new Params instance
func NewParams(owner string) Params {
	return Params{
		Owner: owner,
	}
}

// ParamKeyTable returns the TypeTable for the module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{
		Owner: "",
	}
}

// ValidateParams validates the given params
func ValidateParams(p Params) error {
	if err := validateOwner(p.Owner); err != nil {
		return err
	}

	return nil
}

func validateOwner(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
