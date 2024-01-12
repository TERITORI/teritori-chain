package types

import (
	"cosmossdk.io/errors"
)

// errors
var (
	ErrEmptyAddress = errors.Register(ModuleName, 1, "empty address")
)
