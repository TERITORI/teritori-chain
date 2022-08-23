package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestLastReductionBlockNumGetSet() {
	lastBlockNum := suite.app.MintKeeper.GetLastReductionBlockNum(suite.ctx)
	suite.Require().Equal(lastBlockNum, int64(0))

	suite.app.MintKeeper.SetLastReductionBlockNum(suite.ctx, 100)
	lastBlockNum = suite.app.MintKeeper.GetLastReductionBlockNum(suite.ctx)
	suite.Require().Equal(lastBlockNum, int64(100))
}

func (suite *KeeperTestSuite) TestTeamVestingMonthInfoGetSet() {
	monthInfo := suite.app.MintKeeper.GetTeamVestingMonthInfo(suite.ctx)
	suite.Require().Equal(monthInfo, types.DefaultGenesisState().MonthInfo)

	newMonthInfo := types.TeamVestingMonthInfo{
		MonthsSinceGenesis:     1,
		MonthStartedBlock:      1,
		OneMonthPeriodInBlocks: 10000,
	}
	suite.app.MintKeeper.SetTeamVestingMonthInfo(suite.ctx, newMonthInfo)
	monthInfo = suite.app.MintKeeper.GetTeamVestingMonthInfo(suite.ctx)
	suite.Require().Equal(monthInfo, newMonthInfo)
}

func (suite *KeeperTestSuite) TestMinterGetSet() {
	minterInfo := suite.app.MintKeeper.GetMinter(suite.ctx)
	suite.Require().Equal(minterInfo.BlockProvisions, types.DefaultGenesisState().Params.GenesisBlockProvisions)

	newMinterInfo := types.Minter{
		BlockProvisions: sdk.NewDec(1),
	}
	suite.app.MintKeeper.SetMinter(suite.ctx, newMinterInfo)
	minterInfo = suite.app.MintKeeper.GetMinter(suite.ctx)
	suite.Require().Equal(minterInfo, newMinterInfo)
}
