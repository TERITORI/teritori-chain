package keeper_test

import (
	"github.com/NXTPOP/teritori-chain/x/airdrop/keeper"
)

func (suite *KeeperTestSuite) TestVerifySignature() {
	tests := []struct {
		testCase   string
		chain      string
		address    string
		pubKey     string
		rewardAddr string
		signature  string
		expectPass bool
	}{
		{
			"evm invalid signature verification",
			"evm",
			"0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9",
			"",
			"tori1665x2fj8xyez0vqxs5pjhc6e7ktmmrx9dz864d",
			"0xed41b43ee89d28652b34f0e5c769753ba3b4f0cb6b85bc63f29ba6680cedbf89688d66145b5326c1afb811117a196f0f14b73989b062dde980ce13cd06f1ad801b",
			false,
		},
		{
			"evm successful verification",
			"evm",
			"0x583e8DD54b7C3F5Ea23862E0E852f0e6914475D5",
			"",
			"tori1pkmvlnstq8q7djns3w882pcu92xh4c9x4ukhcd",
			"0xf2cde652dbe26e73e508782d673850ac10880fafef4f7cd2599fd434736ef0ca2d8a2bd65c7f8b67abc1a837f95a23e3c34789dd0ac230cb9b04000641d62b521c",
			true,
		},
	}

	for _, tc := range tests {
		passed := keeper.VerifySignature(tc.chain, tc.address, tc.pubKey, tc.rewardAddr, tc.signature)
		if tc.expectPass {
			suite.Require().True(passed)
		} else {
			suite.Require().False(passed)
		}
	}
}
