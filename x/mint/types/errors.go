package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrEmptyAddress = errors.Register(ModuleName, 1, "empty address")
)
