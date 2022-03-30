package airdrop

import (
	"encoding/hex"
	"fmt"

	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/keeper"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	// 0x028c73bd3ef2c2ce6a520814769faf00f4da1ad11ecd18dcd0c3ea8db12b0fbd41
	pubKeyBytes, _ := hex.DecodeString("028c73bd3ef2c2ce6a520814769faf00f4da1ad11ecd18dcd0c3ea8db12b0fbd41")
	// hrp, bz, err := bech32.DecodeAndConvert("terrapub1addwnpepq2x880f77tpvu6jjpq28d8a0qr6d5xk3rmx33hxsc04gmvftp775zxg87a3")
	// hexValue := hex.EncodeToString(bz)
	// fmt.Println("hrp, bz, err", hrp, hexValue, err)
	pubKey := secp256k1.PubKey{Key: pubKeyBytes}
	fmt.Println(bech32.ConvertAndEncode("terra", pubKey.Address()))
}
