package keeper_test

import (
	"github.com/TERITORI/teritori-chain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestAllocationGetSet() {
	// get allocation for an address before set
	allocation := suite.app.AirdropKeeper.GetAllocation(suite.ctx, "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9")
	suite.Require().Nil(allocation)

	allocations := suite.app.AirdropKeeper.GetAllAllocations(suite.ctx)
	suite.Require().Len(allocations, 6)

	// set allocation
	evmAllocation := types.AirdropAllocation{
		Chain:         "evm",
		Address:       "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9",
		Amount:        sdk.NewInt64Coin("utori", 1000000),
		ClaimedAmount: sdk.NewInt64Coin("utori", 0),
	}
	suite.app.AirdropKeeper.SetAllocation(suite.ctx, evmAllocation)

	// check allocation after set
	allocation = suite.app.AirdropKeeper.GetAllocation(suite.ctx, "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9")
	suite.Require().Equal(*allocation, evmAllocation)

	allocations = suite.app.AirdropKeeper.GetAllAllocations(suite.ctx)
	suite.Require().Len(allocations, 7)

	// check allocation after delete
	suite.app.AirdropKeeper.DeleteAllocation(suite.ctx, "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9")

	allocation = suite.app.AirdropKeeper.GetAllocation(suite.ctx, "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9")
	suite.Require().Nil(allocation)

	allocations = suite.app.AirdropKeeper.GetAllAllocations(suite.ctx)
	suite.Require().Len(allocations, 6)
}
