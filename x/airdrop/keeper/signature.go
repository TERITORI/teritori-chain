package keeper

func verifySignature(chain string, address string, rewardAddr string, signature []byte) bool {
	switch chain {
	case "solana":
		// TODO: implement solana signature verification
		return true
	case "evm":
		// TODO: implement evm signature verification
		return true
	case "terra":
		// TODO: implement converter from terra address to nxtpop address - signature check is not needed
		return true
	case "osmosis":
		// TODO: implement converter from terra address to nxtpop address - signature check is not needed
		return true
	case "juno":
		// TODO: implement converter from terra address to nxtpop address - signature check is not needed
		return true
	case "cosmos":
		return true
	default: // unsupported chain
		return false
	}
}
