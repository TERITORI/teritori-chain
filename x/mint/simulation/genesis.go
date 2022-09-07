package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/TERITORI/teritori-chain/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants.
const (
	blockProvisionsKey         = "genesis_block_provisions"
	reductionFactorKey         = "reduction_factor"
	reductionPeriodInBlocksKey = "reduction_period_in_blocks"

	MintingRewardsDistributionStartBlockKey = "minting_rewards_distribution_start_block"

	maxInt64 = int(^uint(0) >> 1)
)

var (
	distributionProportions = types.DistributionProportions{
		GrantsProgram:    sdk.NewDecWithPrec(10, 2), // 10%
		CommunityPool:    sdk.NewDecWithPrec(10, 2), // 10%
		UsageIncentive:   sdk.NewDecWithPrec(25, 2), // 25%
		Staking:          sdk.NewDecWithPrec(40, 2), // 40%
		DeveloperRewards: sdk.NewDecWithPrec(15, 2), // 15%
	}
	weightedDevRewardReceivers = []types.MonthlyVestingAddress{
		{
			Address:        "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			MonthlyAmounts: []sdk.Int{sdk.NewInt(7000), sdk.NewInt(7000), sdk.NewInt(7000)},
		},
		{
			Address:        "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			MonthlyAmounts: []sdk.Int{sdk.NewInt(2000), sdk.NewInt(2000), sdk.NewInt(2000)},
		},
		{
			Address:        "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			MonthlyAmounts: []sdk.Int{sdk.NewInt(1000), sdk.NewInt(1000), sdk.NewInt(1000)},
		},
	}
)

// RandomizedGenState generates a random GenesisState for mint.
func RandomizedGenState(simState *module.SimulationState) {
	var blockProvisions sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, blockProvisionsKey, &blockProvisions, simState.Rand,
		func(r *rand.Rand) { blockProvisions = genBlockProvisions(r) },
	)

	var reductionFactor sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, reductionFactorKey, &reductionFactor, simState.Rand,
		func(r *rand.Rand) { reductionFactor = genReductionFactor(r) },
	)

	var reductionPeriodInBlocks int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, reductionPeriodInBlocksKey, &reductionPeriodInBlocks, simState.Rand,
		func(r *rand.Rand) { reductionPeriodInBlocks = genReductionPeriodInBlocks(r) },
	)

	var mintintRewardsDistributionStartBlock int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintingRewardsDistributionStartBlockKey, &mintintRewardsDistributionStartBlock, simState.Rand,
		func(r *rand.Rand) { mintintRewardsDistributionStartBlock = genMintintRewardsDistributionStartBlock(r) },
	)

	reductionStartedBlock := genReductionStartedBlock(simState.Rand)

	mintDenom := sdk.DefaultBondDenom
	params := types.NewParams(
		mintDenom,
		blockProvisions,
		reductionFactor,
		reductionPeriodInBlocks,
		distributionProportions,
		weightedDevRewardReceivers,
		mintintRewardsDistributionStartBlock)

	minter := types.NewMinter(blockProvisions)

	mintGenesis := types.NewGenesisState(minter, params, reductionStartedBlock, types.TeamVestingMonthInfo{})

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}

func genBlockProvisions(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(int64(r.Intn(maxInt64)))
}

func genReductionFactor(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(10)), 1)
}

func genReductionPeriodInBlocks(r *rand.Rand) int64 {
	return int64(r.Intn(maxInt64))
}

func genMintintRewardsDistributionStartBlock(r *rand.Rand) int64 {
	return int64(r.Intn(maxInt64))
}

func genReductionStartedBlock(r *rand.Rand) int64 {
	return int64(r.Intn(maxInt64))
}
