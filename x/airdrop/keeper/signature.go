package keeper

import (
	"encoding/hex"
	"encoding/json"

	appparams "github.com/TERITORI/teritori-chain/app/params"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	solana "github.com/gagliardetto/solana-go"
)

type SignMessage struct {
	Chain      string `json:"chain"`
	Address    string `json:"address"`
	RewardAddr string `json:"rewardAddr"`
}

func VerifySignature(chain string, address string, pubKey string, rewardAddr string, signatureBytes string) bool {
	signMsg := SignMessage{
		Chain:      chain,
		Address:    address,
		RewardAddr: rewardAddr,
	}
	signBytes, err := json.Marshal(signMsg)
	if err != nil {
		return false
	}
	// keplrSignBytes := legacytx.StdSignBytes(
	// 	"", 0, 0, 0,
	// 	legacytx.StdFee{Amount: sdk.Coins{}, Gas: 0},
	// 	[]sdk.Msg{types.NewMsgSignData(address, signBytes)}, "",
	// )

	switch chain {
	case "solana":
		pubkey, err := solana.PublicKeyFromBase58(address)
		if err != nil {
			return false
		}
		signatureData, err := hex.DecodeString(signatureBytes[2:])
		if err != nil {
			return false
		}
		signature := solana.SignatureFromBytes(signatureData)
		return signature.Verify(pubkey, signBytes)
	case "evm":
		signatureData, err := hexutil.Decode(signatureBytes)
		if err != nil {
			return false
		}
		if len(signatureData) != crypto.SignatureLength {
			return false
		}
		signatureData[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
		recovered, err := crypto.SigToPub(accounts.TextHash(signBytes), signatureData)
		if err != nil {
			return false
		}
		recoveredAddr := crypto.PubkeyToAddress(*recovered)
		return recoveredAddr.String() == address
	case "terra":
		pubKeyBytes, err := hexutil.Decode(pubKey)
		if err != nil {
			return false
		}
		secp256k1PubKey := secp256k1.PubKey{Key: pubKeyBytes}
		terraAddr, err := bech32.ConvertAndEncode("terra", secp256k1PubKey.Address())
		if err != nil {
			return false
		}
		if terraAddr != address {
			return false
		}

		signatureData, err := hexutil.Decode(signatureBytes)
		if err != nil {
			return false
		}
		return secp256k1PubKey.VerifySignature(signBytes, signatureData)
	// case "secret":
	// 	pubKeyBytes, err := hexutil.Decode(pubKey)
	// 	if err != nil {
	// 		return false
	// 	}
	// 	secp256k1PubKey := secp256k1.PubKey{Key: pubKeyBytes}
	// 	secretAddr, err := bech32.ConvertAndEncode("secret", secp256k1PubKey.Address())
	// 	if err != nil {
	// 		return false
	// 	}
	// 	if secretAddr != address {
	// 		return false
	// 	}

	// 	signatureData, err := hexutil.Decode(signatureBytes)
	// 	if err != nil {
	// 		return false
	// 	}
	// 	return secp256k1PubKey.VerifySignature(keplrSignBytes, signatureData)
	case "stargaze":
		_, bz, err := bech32.DecodeAndConvert(address)
		if err != nil {
			return false
		}

		bech32Addr, err := bech32.ConvertAndEncode(appparams.Bech32PrefixAccAddr, bz)
		if err != nil {
			return false
		}

		return bech32Addr == rewardAddr
	case "osmosis":
		_, bz, err := bech32.DecodeAndConvert(address)
		if err != nil {
			return false
		}

		bech32Addr, err := bech32.ConvertAndEncode(appparams.Bech32PrefixAccAddr, bz)
		if err != nil {
			return false
		}

		return bech32Addr == rewardAddr
	case "juno":
		_, bz, err := bech32.DecodeAndConvert(address)
		if err != nil {
			return false
		}

		bech32Addr, err := bech32.ConvertAndEncode(appparams.Bech32PrefixAccAddr, bz)
		if err != nil {
			return false
		}

		return bech32Addr == rewardAddr
	case "cosmos":
		_, bz, err := bech32.DecodeAndConvert(address)
		if err != nil {
			return false
		}

		bech32Addr, err := bech32.ConvertAndEncode(appparams.Bech32PrefixAccAddr, bz)
		if err != nil {
			return false
		}

		return bech32Addr == rewardAddr
	default: // unsupported chain
		return false
	}
}
