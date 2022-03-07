package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrEmptyRewardAddress                       = errors.Register(ModuleName, 1, "empty reward address")
	ErrEmptyOnChainAllocationAddress            = errors.Register(ModuleName, 2, "empty on-chain allocation address")
	ErrAirdropAllocationDoesNotExists           = errors.Register(ModuleName, 3, "airdrop allocation does not exists for the address")
	ErrAirdropAllocationAlreadyClaimed          = errors.Register(ModuleName, 4, "airdrop allocation is already claimed for the address")
	ErrNativeChainAccountSigVerificationFailure = errors.Register(ModuleName, 5, "native chain account signature verification failure")
)
