package types

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgClaimAllocation(t *testing.T) {
	msg := MsgClaimAllocation{
		Address:       "0x9d967594Cc61453aFEfD657313e5F05be7c6F88F",
		RewardAddress: "pop18mu5hhgy64390q56msql8pfwps0uesn0gf0elf",
		Signature:     "0xb89733c05568385a861fa20f5c4abe53c23a13962515bf5510638b4e3947b1236963b53de549ae762bbd45427dbd3712ae7d169a935d21e44e7da86b1c552f471b",
	}

	fmt.Println("sign_bytes", string(msg.GetSignBytes()))

	require.NoError(t, msg.ValidateBasic())
}

func TestMsgSignData(t *testing.T) {
	msg := MsgSignData{
		Signer: "secret1kap5kvfahhvwufvsj5x6dvyqxztyjhsv8y3da7",
		Data:   []byte(`{"address":"secret1kap5kvfahhvwufvsj5x6dvyqxztyjhsv8y3da7","chain":"secret","rewardAddr":"tori1hwf62gw7h39xmd69st3p487r8x3sphm2hx79ae"}`),
	}

	fmt.Println("sign_bytes", string(msg.GetSignBytes()))

	require.NoError(t, msg.ValidateBasic())
}
