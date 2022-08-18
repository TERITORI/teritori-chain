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
	weightedDevRewardReceivers = []types.WeightedAddress{
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(2887, 4),
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(229, 3),
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(1625, 4),
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(109, 3),
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(995, 3).Quo(sdk.NewDec(10)), // 0.0995
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(6, 1).Quo(sdk.NewDec(10)), // 0.06
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(15, 2).Quo(sdk.NewDec(10)), // 0.015
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(1, 1).Quo(sdk.NewDec(10)), // 0.01
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(75, 2).Quo(sdk.NewDec(100)), // 0.0075
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(7, 1).Quo(sdk.NewDec(100)), // 0.007
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(5, 1).Quo(sdk.NewDec(100)), // 0.005
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(25, 2).Quo(sdk.NewDec(100)), // 0.0025
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(25, 2).Quo(sdk.NewDec(100)), // 0.0025
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(1, 1).Quo(sdk.NewDec(100)), // 0.001
		},
		{
			Address: "tori1g2escsu26508tgrpv865d80d62pvmw69je2ztn",
			Weight:  sdk.NewDecWithPrec(8, 1).Quo(sdk.NewDec(1000)), // 0.0008
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
