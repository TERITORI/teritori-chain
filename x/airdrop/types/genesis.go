package types

import "github.com/cosmos/cosmos-sdk/types"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Allocations: []*AirdropAllocation{
			{
				Chain:         "evm",
				Address:       "--",
				Amount:        types.NewCoin("upop", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("upop", types.NewIntWithDecimal(0, 6)),
			},
			{
				Chain:         "solana",
				Address:       "--",
				Amount:        types.NewCoin("upop", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("upop", types.NewIntWithDecimal(0, 10)),
			},
			{
				Chain:         "terra",
				Address:       "--",
				Amount:        types.NewCoin("upop", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("upop", types.NewIntWithDecimal(0, 10)),
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
