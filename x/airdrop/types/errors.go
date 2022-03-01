package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrEmptyRewardAddress            = errors.Register(ModuleName, 1, "empty reward address")
	ErrEmptyOnChainAllocationAddress = errors.Register(ModuleName, 2, "empty on-chain allocation address")
)
