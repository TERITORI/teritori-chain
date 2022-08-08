package types

// constants
var (
	ModuleName = "nftstaking"

	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixKeyNftStaking   = []byte{0x0}
	PrefixKeyAccessInfo   = []byte{0x1}
	PrefixKeyNftTypePerms = []byte{0x2}
)
