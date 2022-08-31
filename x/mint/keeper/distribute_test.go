package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestDistributeMintedCoin() {
	grantsAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	usageIncentiveAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	dev1Addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	dev2Addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	teamReserveAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	params := types.Params{
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
				Address:        dev1Addr.String(),
				MonthlyAmounts: []sdk.Int{sdk.NewInt(7000000), sdk.NewInt(7000000), sdk.NewInt(7000000)},
			},
			{
				Address:        dev2Addr.String(),
				MonthlyAmounts: []sdk.Int{sdk.NewInt(3000000), sdk.NewInt(3000000), sdk.NewInt(3000000)},
			},
		},
		UsageIncentiveAddress:                usageIncentiveAddr.String(),
		GrantsProgramAddress:                 grantsAddr.String(),
		TeamReserveAddress:                   teamReserveAddr.String(),
		MintingRewardsDistributionStartBlock: 1,
	}

	tests := []struct {
		testCase    string
		mintedCoins sdk.Coins
		monthIndex  int64
		expectPass  bool
	}{
		{
			"first month distribution test",
			sdk.Coins{sdk.NewInt64Coin(params.MintDenom, 1000000)},
			0,
			true,
		},
		{
			"second month distribution test",
			sdk.Coins{sdk.NewInt64Coin(params.MintDenom, 1000000)},
			1,
			true,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()
		suite.app.MintKeeper.SetParams(suite.ctx, params)

		newMonthInfo := types.TeamVestingMonthInfo{
			MonthsSinceGenesis:     tc.monthIndex,
			MonthStartedBlock:      1,
			OneMonthPeriodInBlocks: 43200,
		}
		suite.app.MintKeeper.SetTeamVestingMonthInfo(suite.ctx, newMonthInfo)

		err := suite.app.MintKeeper.MintCoins(suite.ctx, tc.mintedCoins)
		if err != nil {
			panic(err)
		}

		err = suite.app.MintKeeper.DistributeMintedCoin(suite.ctx, tc.mintedCoins[0])
		if tc.expectPass {
			suite.Require().NoError(err)

			// check grants amount is distributed correctly
			grantsAddrBalance := suite.app.BankKeeper.GetBalance(suite.ctx, grantsAddr, params.MintDenom)
			suite.Require().Equal(grantsAddrBalance, sdk.NewInt64Coin(params.MintDenom, 200000))

			// check usage incentive amount is distributed correctly
			usageIncentiveAddrBalance := suite.app.BankKeeper.GetBalance(suite.ctx, usageIncentiveAddr, params.MintDenom)
			suite.Require().Equal(usageIncentiveAddrBalance, sdk.NewInt64Coin(params.MintDenom, 200000))

			// check staking incentive amount is distributed correctly
			feeCollectorAddr := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
			feeCollectorAddrBalance := suite.app.BankKeeper.GetBalance(suite.ctx, feeCollectorAddr, params.MintDenom)
			suite.Require().Equal(feeCollectorAddrBalance, sdk.NewInt64Coin(params.MintDenom, 200000))

			// check community pool amount is distributed correctly: community pool balance check on distr module
			communityPoolCoins := suite.app.DistrKeeper.GetFeePoolCommunityCoins(suite.ctx)
			suite.Require().Equal(communityPoolCoins, sdk.DecCoins{sdk.NewInt64DecCoin(params.MintDenom, 200000)})

			// check developer reward amount is distributed correctly per month: each address registered on weighted
			dev1AddrCoins := suite.app.BankKeeper.GetBalance(suite.ctx, dev1Addr, params.MintDenom)
			dev1Expected := params.WeightedDeveloperRewardsReceivers[0].MonthlyAmounts[tc.monthIndex].Quo(sdk.NewInt(newMonthInfo.OneMonthPeriodInBlocks))
			suite.Require().Equal(dev1AddrCoins, sdk.NewCoin(params.MintDenom, dev1Expected))
			dev2Expected := params.WeightedDeveloperRewardsReceivers[1].MonthlyAmounts[tc.monthIndex].Quo(sdk.NewInt(newMonthInfo.OneMonthPeriodInBlocks))
			dev2AddrCoins := suite.app.BankKeeper.GetBalance(suite.ctx, dev2Addr, params.MintDenom)
			suite.Require().Equal(dev2AddrCoins, sdk.NewCoin(params.MintDenom, dev2Expected))

			// check team reserve balance
			teamReserveAddrCoins := suite.app.BankKeeper.GetBalance(suite.ctx, teamReserveAddr, params.MintDenom)
			suite.Require().Equal(teamReserveAddrCoins, sdk.NewCoin(params.MintDenom, sdk.NewInt(200000).Sub(dev1Expected).Sub(dev2Expected)))
		} else {
			suite.Require().Error(err)
		}
	}

}
