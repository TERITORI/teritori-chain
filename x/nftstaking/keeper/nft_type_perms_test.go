package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestNftTypePermsGetSet() {
	newNftType := types.NftType(100)

	nftTypePerms := suite.app.NftStakingKeeper.GetNftTypePerms(suite.ctx, newNftType)
	suite.Require().Equal(nftTypePerms, types.NftTypePerms{})

	allNftTypePerms := suite.app.NftStakingKeeper.GetAllNftTypePerms(suite.ctx)
	suite.Require().Len(allNftTypePerms, 0)

	// set
	newNftTypePerms := types.NftTypePerms{
		NftType: newNftType,
		Perms:   []types.Permission{1, 2},
	}
	suite.app.NftStakingKeeper.SetNftTypePerms(suite.ctx, newNftTypePerms)

	// check after set
	nftTypePerms = suite.app.NftStakingKeeper.GetNftTypePerms(suite.ctx, newNftType)
	suite.Require().Equal(nftTypePerms, newNftTypePerms)

	allNftTypePerms = suite.app.NftStakingKeeper.GetAllNftTypePerms(suite.ctx)
	suite.Require().Len(allNftTypePerms, 1)

	// check after delete
	suite.app.NftStakingKeeper.DeleteNftTypePerms(suite.ctx, newNftType)

	nftTypePerms = suite.app.NftStakingKeeper.GetNftTypePerms(suite.ctx, newNftType)
	suite.Require().Equal(nftTypePerms, types.NftTypePerms{})

	allNftTypePerms = suite.app.NftStakingKeeper.GetAllNftTypePerms(suite.ctx)
	suite.Require().Len(allNftTypePerms, 0)
}

func (suite *KeeperTestSuite) TestHasPermission() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	newNftType := types.NftType(100)
	hasPermission := suite.app.NftStakingKeeper.HasPermission(suite.ctx, addr.String(), types.Permission(1))
	suite.Require().Equal(hasPermission, false)

	// set
	newNftTypePerms := types.NftTypePerms{
		NftType: newNftType,
		Perms:   []types.Permission{1, 2},
	}
	suite.app.NftStakingKeeper.SetNftTypePerms(suite.ctx, newNftTypePerms)
	suite.app.NftStakingKeeper.SetNftStaking(suite.ctx, types.NftStaking{
		NftType:       newNftType,
		NftIdentifier: "nft1",
		NftMetadata:   "nftmetadata1",
		RewardAddress: addr.String(),
		RewardWeight:  1,
	})

	hasPermission = suite.app.NftStakingKeeper.HasPermission(suite.ctx, addr.String(), types.Permission(1))
	suite.Require().Equal(hasPermission, true)
}
