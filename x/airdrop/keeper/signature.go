package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	appparams "github.com/POPSmartContract/nxtpop-chain/app/params"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solana "github.com/gagliardetto/solana-go"
)

type SignMessage struct {
	Chain      string `json:"chain"`
	Address    string `json:"address"`
	RewardAddr string `json:"rewardAddr"`
}

func verifySignature(chain string, address string, rewardAddr string, signatureBytes string) bool {
	fmt.Println("chain = ", chain)
	fmt.Println("address = ", address)
	fmt.Println("rewardAddr = ", rewardAddr)
	fmt.Println("signatureBytes = ", signatureBytes)

	signMsg := SignMessage{
		Chain:      chain,
		Address:    address,
		RewardAddr: rewardAddr,
	}
	signBytes, err := json.Marshal(signMsg)
	fmt.Println("signBytes = ", signBytes)
	fmt.Println("signBytesString = ", string(signBytes))

	if err != nil {
		return false
	}

	switch chain {
	case "solana":
		pubkey := solana.MustPublicKeyFromBase58(address)
		signatureData, err := hex.DecodeString(signatureBytes[2:])
		if err != nil {
			return false
		}
		signature := solana.SignatureFromBytes(signatureData)
		return signature.Verify(pubkey, signBytes)
	case "evm":
		// TODO: implement evm signature verification
		// ethsecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
		// if !ethsecp256k1.VerifySignature(pubKey.Bytes(), sigHash, data.Signature[:len(data.Signature)-1]) {
		// }

		// "github.com/ethereum/go-ethereum/common"
		// "github.com/ethereum/go-ethereum/crypto"
		// code is derivated from github.com/ethereum/go-ethereum
		// func assertSignature(addr common.Address, index uint64, hash [32]byte, r, s [32]byte, v uint8, expect common.Address) bool {
		// 	buf := make([]byte, 8)
		// 	binary.BigEndian.PutUint64(buf, index)
		// 	data := append([]byte{0x19, 0x00}, append(addr.Bytes(), append(buf, hash[:]...)...)...)

		// 	pubkey, err := crypto.Ecrecover(crypto.Keccak256(data), append(r[:], append(s[:], v-27)...))
		// 	if err != nil {
		// 		return false
		// 	}
		// 	var signer common.Address
		// 	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

		// 	return bytes.Equal(signer.Bytes(), expect.Bytes())
		// }

		// TODO: should check this code's accuracy
		// TODO: how to get r,s,v from signatureBytes?
		var r, s [32]byte
		var v uint8
		ethpubkey, err := crypto.Ecrecover(crypto.Keccak256(signBytes), append(r[:], append(s[:], v-27)...))
		if err != nil {
			return false
		}
		var signer common.Address
		copy(signer[:], crypto.Keccak256(ethpubkey[1:])[12:])

		return signer.String() == address
	case "terra":
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
