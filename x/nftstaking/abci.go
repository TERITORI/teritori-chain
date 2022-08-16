package nftstaking

import (
	"github.com/TERITORI/teritori-chain/x/nftstaking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	stakings := k.GetAllNftStakings(ctx)

	totalInflation := sdk.NewInt(1000_000_000)
	totalStakingPower := uint64(0)

	for _, staking := range stakings {
		totalStakingPower += staking.RewardWeight
	}

	for _, staking := range stakings {
		reward := totalInflation.Mul(sdk.NewInt(int64(staking.RewardWeight))).Quo(sdk.NewInt(int64(totalStakingPower)))

		rewardAddr, err := sdk.AccAddressFromBech32(staking.RewardAddress)
		if err != nil {
			continue
		}
		cachedCtx, write := ctx.CacheContext()
		err = k.AllocateTokensToRewardAddress(cachedCtx, rewardAddr, sdk.NewCoin(k.BondDenom(ctx), reward))
		if err == nil {
			write()
		}
	}
}
