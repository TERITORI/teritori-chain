package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestAccessInfoGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	accessInfo := suite.app.NftStakingKeeper.GetAccessInfo(suite.ctx, addr.String())
	suite.Require().Equal(accessInfo, types.Access{})

	accessInfos := suite.app.NftStakingKeeper.GetAllAccessInfos(suite.ctx)
	suite.Require().Len(accessInfos, 0)

	// set
	accessInfoRaw := types.Access{
		Address: addr.String(),
		Servers: []types.ServerAccess{{
			Server:   "teritori",
			Channels: []string{"announcement", "twitter-feed"},
		}},
	}
	suite.app.NftStakingKeeper.SetAccessInfo(suite.ctx, accessInfoRaw)

	// check after set
	accessInfo = suite.app.NftStakingKeeper.GetAccessInfo(suite.ctx, addr.String())
	suite.Require().Equal(accessInfo, accessInfoRaw)

	accessInfos = suite.app.NftStakingKeeper.GetAllAccessInfos(suite.ctx)
	suite.Require().Len(accessInfos, 1)

	// check after delete
	suite.app.NftStakingKeeper.DeleteAccessInfo(suite.ctx, addr.String())

	accessInfo = suite.app.NftStakingKeeper.GetAccessInfo(suite.ctx, addr.String())
	suite.Require().Equal(accessInfo, types.Access{})

	accessInfos = suite.app.NftStakingKeeper.GetAllAccessInfos(suite.ctx)
	suite.Require().Len(accessInfos, 0)
}
