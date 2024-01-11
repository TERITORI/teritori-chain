package keeper

import (
	"cosmossdk.io/math"
	"github.com/TERITORI/teritori-chain/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// DistributeMintedCoins implements distribution of minted coins from mint to external modules.
func (k Keeper) DistributeMintedCoin(ctx sdk.Context, mintedCoin sdk.Coin) error {
	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	grantsAmount, err := k.distributeToAddress(ctx, params.GrantsProgramAddress, mintedCoin, proportions.GrantsProgram)
	if err != nil {
		return err
	}

	usageIncentiveAmount, err := k.distributeToAddress(ctx, params.UsageIncentiveAddress, mintedCoin, proportions.UsageIncentive)
	if err != nil {
		return err
	}

	// allocate staking incentives into fee collector account to be moved to on next begin blocker by staking module account.
	stakingIncentivesAmount, err := k.distributeToModule(ctx, k.feeCollectorName, mintedCoin, proportions.Staking)
	if err != nil {
		return err
	}

	// allocate dev rewards to respective accounts from developer vesting module account.
	devRewardAmount, err := k.distributeDeveloperRewards(ctx, mintedCoin, proportions.DeveloperRewards, params.WeightedDeveloperRewardsReceivers)
	if err != nil {
		return err
	}

	// subtract from original provision to ensure no coins left over after the allocations
	communityPoolAmount := mintedCoin.Amount.Sub(grantsAmount).Sub(usageIncentiveAmount).Sub(stakingIncentivesAmount).Sub(devRewardAmount)
	err = k.communityPoolKeeper.FundCommunityPool(ctx, sdk.NewCoins(sdk.NewCoin(params.MintDenom, communityPoolAmount)), k.accountKeeper.GetModuleAddress(types.ModuleName))
	if err != nil {
		return err
	}

	// call an hook after the minting and distribution of new coins
	if k.hooks != nil {
		k.hooks.AfterDistributeMintedCoin(ctx)
	}

	return err
}

// distributeToModule distributes mintedCoin multiplied by proportion to the recepient account.
func (k Keeper) distributeToAddress(ctx sdk.Context, recipientAddr string, mintedCoin sdk.Coin, proportion sdk.Dec) (math.Int, error) {
	distributionCoin, err := getProportions(mintedCoin, proportion)
	if err != nil {
		return math.Int{}, err
	}

	recipient, err := sdk.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return math.Int{}, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.NewCoins(distributionCoin)); err != nil {
		return math.Int{}, err
	}
	return distributionCoin.Amount, nil
}

// distributeToModule distributes mintedCoin multiplied by proportion to the recepientModule account.
func (k Keeper) distributeToModule(ctx sdk.Context, recipientModule string, mintedCoin sdk.Coin, proportion sdk.Dec) (math.Int, error) {
	distributionCoin, err := getProportions(mintedCoin, proportion)
	if err != nil {
		return math.Int{}, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, sdk.NewCoins(distributionCoin)); err != nil {
		return math.Int{}, err
	}
	return distributionCoin.Amount, nil
}

func (k Keeper) distributeDeveloperRewards(ctx sdk.Context, totalMintedCoin sdk.Coin, developerRewardsProportion sdk.Dec, developerRewardsReceivers []types.MonthlyVestingAddress) (math.Int, error) {

	params := k.GetParams(ctx)
	totalDevRewards, err := getProportions(totalMintedCoin, developerRewardsProportion)
	if err != nil {
		return math.Int{}, err
	}

	vestedAmount := sdk.ZeroInt()
	// allocate developer rewards to addresses by weight
	for _, w := range developerRewardsReceivers {
		monthInfo := k.GetTeamVestingMonthInfo(ctx)
		if len(w.MonthlyAmounts) <= int(monthInfo.MonthsSinceGenesis) {
			continue
		}
		devPortionAmount := w.MonthlyAmounts[monthInfo.MonthsSinceGenesis].Quo(sdk.NewInt(monthInfo.OneMonthPeriodInBlocks))
		if devPortionAmount.IsZero() {
			continue
		}
		devRewardPortionCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, devPortionAmount))
		// fund community pool when rewards address is empty.
		if w.Address != emptyAddressReceiver {
			devRewardsAddr, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return math.Int{}, err
			}
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, devRewardsAddr, devRewardPortionCoins)
			if err != nil {
				return math.Int{}, err
			}

			vestedAmount = vestedAmount.Add(devPortionAmount)
		}
	}
	// send remaining tokens to team reserve
	vestedTokens := sdk.NewCoin(params.MintDenom, vestedAmount)
	remainingCoins := totalDevRewards.Sub(vestedTokens)
	if remainingCoins.IsPositive() {
		reserve, err := sdk.AccAddressFromBech32(params.TeamReserveAddress)
		if err != nil {
			panic(err)
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, reserve, sdk.Coins{remainingCoins})
		if err != nil {
			return math.Int{}, err
		}
	}

	return totalDevRewards.Amount, nil
}

func getProportions(mintedCoin sdk.Coin, ratio sdk.Dec) (sdk.Coin, error) {
	if ratio.GT(sdk.OneDec()) {
		return sdk.Coin{}, invalidRatioError{ratio}
	}
	return sdk.NewCoin(mintedCoin.Denom, sdk.NewDecFromInt(mintedCoin.Amount).Mul(ratio).TruncateInt()), nil
}
