package keeper

import (
	"encoding/json"

	solana "github.com/gagliardetto/solana-go"
)

type SignMessage struct {
	Chain      string `json:"string"`
	Address    string `json:"address"`
	RewardAddr string `json:"rewardAddr"`
}

func verifySignature(chain string, address string, rewardAddr string, signatureBytes []byte) bool {
	signMsg := SignMessage{
		Chain:      chain,
		Address:    address,
		RewardAddr: rewardAddr,
	}
	signByes, err := json.Marshal(signMsg)

	if err != nil {
		return false
	}

	switch chain {
	case "solana":
		pubkey := solana.MustPublicKeyFromBase58(address)
		signature := solana.SignatureFromBytes(signatureBytes)
		return signature.Verify(pubkey, signByes)
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
