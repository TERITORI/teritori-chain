package types

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(minter Minter, params Params, reductionStartedBlock int64, monthInfo TeamVestingMonthInfo) *GenesisState {
	return &GenesisState{
		Minter:                minter,
		Params:                params,
		ReductionStartedBlock: reductionStartedBlock,
		MonthInfo:             monthInfo,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Minter:                DefaultInitialMinter(),
		Params:                DefaultParams(),
		ReductionStartedBlock: 0,
		MonthInfo: TeamVestingMonthInfo{
			OneMonthPeriodInBlocks: 525600, // 1 month - 86400 x 365 / 12 / 5			,
		},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return data.Minter.Validate()
}
