package types

import "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEmptySender        = errors.Register(ModuleName, 1, "empty sender")
	ErrEmptyRewardAddress = errors.Register(ModuleName, 2, "empty reward address")
)
