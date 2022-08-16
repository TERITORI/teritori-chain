package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestNftStakingGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	nftStaking := suite.app.NftStakingKeeper.GetNftStaking(suite.ctx, addr.String())
	suite.Require().Equal(nftStaking, types.NftStaking{})

	nftStakings := suite.app.NftStakingKeeper.GetAllNftStakings(suite.ctx)
	suite.Require().Len(nftStakings, 0)

	// set
	nftStakingRaw := types.NftStaking{
		NftType:       types.NftType_NftTypeDefault,
		NftIdentifier: "nft1",
		NftMetadata:   "metadata1",
		RewardAddress: addr.String(),
		RewardWeight:  1,
	}
	suite.app.NftStakingKeeper.SetNftStaking(suite.ctx, nftStakingRaw)

	// check after set
	nftStaking = suite.app.NftStakingKeeper.GetNftStaking(suite.ctx, nftStakingRaw.NftIdentifier)
	suite.Require().Equal(nftStaking, nftStakingRaw)

	nftStakings = suite.app.NftStakingKeeper.GetAllNftStakings(suite.ctx)
	suite.Require().Len(nftStakings, 1)

	// check after delete
	suite.app.NftStakingKeeper.DeleteNftStaking(suite.ctx, nftStaking)

	nftStaking = suite.app.NftStakingKeeper.GetNftStaking(suite.ctx, addr.String())
	suite.Require().Equal(nftStaking, types.NftStaking{})

	nftStakings = suite.app.NftStakingKeeper.GetAllNftStakings(suite.ctx)
	suite.Require().Len(nftStakings, 0)
}

func (suite *KeeperTestSuite) TestAllocateTokensToRewardAddress() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	err := suite.app.NftStakingKeeper.AllocateTokensToRewardAddress(suite.ctx, addr, sdk.NewInt64Coin("stake", 1000000))
	suite.Require().NoError(err)

	balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr, "stake")
	suite.Require().True(balance.Equal(sdk.NewInt64Coin("stake", 1000000)))
}
