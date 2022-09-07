package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestParamsGetSet() {
	params := suite.app.MintKeeper.GetParams(suite.ctx)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	params = types.Params{
		MintDenom:               "utori",
		GenesisBlockProvisions:  sdk.NewDec(1000),
		ReductionPeriodInBlocks: 86400,
		ReductionFactor:         sdk.NewDecWithPrec(5, 1),
		DistributionProportions: types.DistributionProportions{
			GrantsProgram:    sdk.NewDecWithPrec(2, 1),
			CommunityPool:    sdk.NewDecWithPrec(2, 1),
			UsageIncentive:   sdk.NewDecWithPrec(2, 1),
			Staking:          sdk.NewDecWithPrec(2, 1),
			DeveloperRewards: sdk.NewDecWithPrec(2, 1),
		},
		WeightedDeveloperRewardsReceivers: []types.MonthlyVestingAddress{
			{
				Address:        "",
				MonthlyAmounts: []sdk.Int{sdk.NewInt(7000), sdk.NewInt(7000), sdk.NewInt(7000)},
			},
		},
		UsageIncentiveAddress:                addr.String(),
		GrantsProgramAddress:                 addr.String(),
		TeamReserveAddress:                   addr.String(),
		MintingRewardsDistributionStartBlock: 1,
	}

	suite.app.MintKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.MintKeeper.GetParams(suite.ctx)
	suite.Require().Equal(params, newParams)
}
