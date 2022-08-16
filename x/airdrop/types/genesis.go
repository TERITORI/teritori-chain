package types

import "github.com/cosmos/cosmos-sdk/types"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Allocations: []AirdropAllocation{
			{
				Chain:         "evm",
				Address:       "0x--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 6)),
			},
			{
				Chain:         "solana",
				Address:       "--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 10)),
			},
			{
				Chain:         "terra",
				Address:       "terra--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 10)),
			},
			{
				Chain:         "cosmos",
				Address:       "cosmos--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 10)),
			},
			{
				Chain:         "juno",
				Address:       "juno--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 10)),
			},
			{
				Chain:         "osmosis",
				Address:       "osmo--",
				Amount:        types.NewCoin("utori", types.NewIntWithDecimal(100, 6)),
				ClaimedAmount: types.NewCoin("utori", types.NewIntWithDecimal(0, 10)),
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
