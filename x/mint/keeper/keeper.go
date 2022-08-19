package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/TERITORI/teritori-chain/x/mint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the mint store.
type Keeper struct {
	cdc                 codec.BinaryCodec
	storeKey            sdk.StoreKey
	paramSpace          paramtypes.Subspace
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	communityPoolKeeper types.CommunityPoolKeeper
	hooks               types.MintHooks
	feeCollectorName    string
}

type invalidRatioError struct {
	ActualRatio sdk.Dec
}

func (e invalidRatioError) Error() string {
	return fmt.Sprintf("mint allocation ratio (%s) is greater than 1", e.ActualRatio)
}

type insufficientDevVestingBalanceError struct {
	ActualBalance         sdk.Int
	AttemptedDistribution sdk.Int
}

func (e insufficientDevVestingBalanceError) Error() string {
	return fmt.Sprintf("developer vesting balance (%s) is smaller than requested distribution of (%s)", e.ActualBalance, e.AttemptedDistribution)
}

const emptyWeightedAddressReceiver = ""

// NewKeeper creates a new mint Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, ck types.CommunityPoolKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                 cdc,
		storeKey:            key,
		paramSpace:          paramSpace,
		accountKeeper:       ak,
		bankKeeper:          bk,
		communityPoolKeeper: ck,
		feeCollectorName:    feeCollectorName,
	}
}

// _____________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Set the mint hooks.
func (k *Keeper) SetHooks(h types.MintHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set mint hooks twice")
	}

	k.hooks = h

	return k
}

// GetLastReductionBlockNum returns last reduction block number.
func (k Keeper) GetLastReductionBlockNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastReductionBlockKey)
	if b == nil {
		return 0
	}

	return int64(sdk.BigEndianToUint64(b))
}

// SetLastReductionBlockNum set last reduction block number.
func (k Keeper) SetLastReductionBlockNum(ctx sdk.Context, blockNum int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastReductionBlockKey, sdk.Uint64ToBigEndian(uint64(blockNum)))
}

// GetTeamVestingMonthInfo returns month information for team vesting
func (k Keeper) GetTeamVestingMonthInfo(ctx sdk.Context) types.TeamVestingMonthInfo {
	store := ctx.KVStore(k.storeKey)

	monthInfo := types.TeamVestingMonthInfo{}
	bz := store.Get(types.TeamVestingMonthInfoKey)
	if bz == nil {
		return monthInfo
	}

	k.cdc.MustUnmarshal(bz, &monthInfo)
	return monthInfo
}

// SetLastReductionBlockNum set last reduction block number.
func (k Keeper) SetTeamVestingMonthInfo(ctx sdk.Context, monthInfo types.TeamVestingMonthInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&monthInfo)
	store.Set(types.TeamVestingMonthInfoKey, bz)
}

// get the minter.
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter.
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// _____________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// _____________________________________________________________________

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
	k.hooks.AfterDistributeMintedCoin(ctx)

	return err
}

// distributeToModule distributes mintedCoin multiplied by proportion to the recepient account.
func (k Keeper) distributeToAddress(ctx sdk.Context, recipientAddr string, mintedCoin sdk.Coin, proportion sdk.Dec) (sdk.Int, error) {
	distributionCoin, err := getProportions(mintedCoin, proportion)
	if err != nil {
		return sdk.Int{}, err
	}

	recipient, err := sdk.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return sdk.Int{}, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.NewCoins(distributionCoin)); err != nil {
		return sdk.Int{}, err
	}
	return distributionCoin.Amount, nil
}

// distributeToModule distributes mintedCoin multiplied by proportion to the recepientModule account.
func (k Keeper) distributeToModule(ctx sdk.Context, recipientModule string, mintedCoin sdk.Coin, proportion sdk.Dec) (sdk.Int, error) {
	distributionCoin, err := getProportions(mintedCoin, proportion)
	if err != nil {
		return sdk.Int{}, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, sdk.NewCoins(distributionCoin)); err != nil {
		return sdk.Int{}, err
	}
	return distributionCoin.Amount, nil
}

func (k Keeper) distributeDeveloperRewards(ctx sdk.Context, totalMintedCoin sdk.Coin, developerRewardsProportion sdk.Dec, developerRewardsReceivers []types.WeightedAddress) (sdk.Int, error) {
	monthlyPercentages := []float32{4.79, 5.30, 5.83, 6.38, 6.95, 7.51, 8.06, 8.59, 9.08, 9.52, 9.89, 10.18, 15.57, 15.72, 15.72, 15.57, 15.27, 14.83, 14.28, 13.62, 12.89, 12.09, 11.27, 10.42, 14.36, 13.12, 11.92, 10.78, 9.70, 8.69, 7.76, 6.91, 6.14, 5.43, 4.80, 4.23, 5.59, 4.91, 4.31, 3.78, 3.31, 2.89, 2.53, 2.21, 1.93, 1.68, 1.47}
	monthPercentage := float32(0)

	totalDevRewards, err := getProportions(totalMintedCoin, developerRewardsProportion)
	if err != nil {
		return sdk.Int{}, err
	}

	monthInfo := k.GetTeamVestingMonthInfo(ctx)
	if len(monthlyPercentages) > int(monthInfo.MonthsSinceGenesis) {
		monthPercentage = monthlyPercentages[monthInfo.MonthsSinceGenesis]
	}

	vestedTokens, err := getProportions(totalMintedCoin, sdk.NewDec(int64(monthPercentage*100)).QuoInt(sdk.NewInt(100)))
	if err != nil {
		return sdk.Int{}, err
	}

	if vestedTokens.Amount.GT(totalDevRewards.Amount) {
		vestedTokens.Amount = totalDevRewards.Amount
	}

	remainingCoins := totalDevRewards.Sub(vestedTokens)

	// allocate developer rewards to addresses by weight
	for _, w := range developerRewardsReceivers {
		devPortionCoin, err := getProportions(vestedTokens, w.Weight)
		if err != nil {
			return sdk.Int{}, err
		}

		if devPortionCoin.IsZero() {
			continue
		}
		devRewardPortionCoins := sdk.NewCoins(devPortionCoin)
		// fund community pool when rewards address is empty.
		if w.Address == emptyWeightedAddressReceiver {
			remainingCoins = remainingCoins.Add(devPortionCoin)
		} else {
			devRewardsAddr, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return sdk.Int{}, err
			}
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, devRewardsAddr, devRewardPortionCoins)
			if err != nil {
				return sdk.Int{}, err
			}
		}
	}

	// send remaining tokens to team reserve
	params := k.GetParams(ctx)

	if remainingCoins.IsPositive() {
		reserve, err := sdk.AccAddressFromBech32(params.TeamReserveAddress)
		if err != nil {
			panic(err)
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, reserve, sdk.Coins{remainingCoins})
		if err != nil {
			return sdk.Int{}, err
		}
	}

	return totalDevRewards.Amount, nil
}

func getProportions(mintedCoin sdk.Coin, ratio sdk.Dec) (sdk.Coin, error) {
	if ratio.GT(sdk.OneDec()) {
		return sdk.Coin{}, invalidRatioError{ratio}
	}
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt()), nil
}
