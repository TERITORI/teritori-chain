package types

// MinterKey is the key to use for the keeper store at which
// the Minter and its BlockProvisions are stored.
var MinterKey = []byte{0x00}

// LastReductionBlockKey is the key to use for the keeper store
// for storing the last block at which reduction occurred.
var LastReductionBlockKey = []byte{0x03}

// TeamVestingMonthInfoKey is the key to use to store month information since genesis
// for non-linear team token vesting
var TeamVestingMonthInfoKey = []byte{0x04}

const (
	// ModuleName is the module name.
	ModuleName = "mint"
	// DeveloperVestingModuleAcctName is the module acct name for developer vesting.
	DeveloperVestingModuleAcctName = "developer_vesting_unvested"

	// StoreKey is the default store key for mint.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// QueryParameters is an endpoint path for querying mint parameters.
	QueryParameters = "parameters"

	// QueryBlockProvisions is an endpoint path for querying mint block provisions.
	QueryBlockProvisions = "block_provisions"
)
