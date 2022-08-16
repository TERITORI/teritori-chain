package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		NftStakings:  []NftStaking{},
		NftTypePerms: []NftTypePerms{},
		AccessInfos:  []Access{},
	}
}
