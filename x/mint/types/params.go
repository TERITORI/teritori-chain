package types

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyMintDenom                            = []byte("MintDenom")
	KeyGenesisBlockProvisions               = []byte("GenesisBlockProvisions")
	KeyReductionPeriodInBlocks              = []byte("ReductionPeriodInBlocks")
	KeyReductionFactor                      = []byte("ReductionFactor")
	KeyPoolAllocationRatio                  = []byte("PoolAllocationRatio")
	KeyDeveloperRewardsReceiver             = []byte("DeveloperRewardsReceiver")
	KeyMintingRewardsDistributionStartBlock = []byte("MintingRewardsDistributionStartBlock")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams returns new mint module parameters initialized to the given values.
func NewParams(
	mintDenom string, genesisBlockProvisions sdk.Dec,
	ReductionFactor sdk.Dec, reductionPeriodInBlocks int64, distrProportions DistributionProportions,
	weightedDevRewardsReceivers []WeightedAddress, MintingRewardsDistributionStartBlock int64,
) Params {
	return Params{
		MintDenom:                            mintDenom,
		GenesisBlockProvisions:               genesisBlockProvisions,
		ReductionPeriodInBlocks:              reductionPeriodInBlocks,
		ReductionFactor:                      ReductionFactor,
		DistributionProportions:              distrProportions,
		WeightedDeveloperRewardsReceivers:    weightedDevRewardsReceivers,
		MintingRewardsDistributionStartBlock: MintingRewardsDistributionStartBlock,
	}
}

// DefaultParams returns the default minting module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:               sdk.DefaultBondDenom,
		GenesisBlockProvisions:  sdk.NewDec(5000000),
		ReductionPeriodInBlocks: 156,                      // 3 years
		ReductionFactor:         sdk.NewDecWithPrec(5, 1), // 0.5
		DistributionProportions: DistributionProportions{
			Staking:          sdk.NewDecWithPrec(4, 1), // 0.4
			PoolIncentives:   sdk.NewDecWithPrec(3, 1), // 0.3
			DeveloperRewards: sdk.NewDecWithPrec(2, 1), // 0.2
			CommunityPool:    sdk.NewDecWithPrec(1, 1), // 0.1
		},
		WeightedDeveloperRewardsReceivers:    []WeightedAddress{},
		MintingRewardsDistributionStartBlock: 0,
	}
}

// Validate validates mint module parameters. Returns nil if valid,
// error otherwise
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateGenesisBlockProvisions(p.GenesisBlockProvisions); err != nil {
		return err
	}
	if err := validateReductionPeriodInBlocks(p.ReductionPeriodInBlocks); err != nil {
		return err
	}
	if err := validateReductionFactor(p.ReductionFactor); err != nil {
		return err
	}
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}
	if err := validateWeightedDeveloperRewardsReceivers(p.WeightedDeveloperRewardsReceivers); err != nil {
		return err
	}
	if err := validateMintingRewardsDistributionStartBlock(p.MintingRewardsDistributionStartBlock); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyGenesisBlockProvisions, &p.GenesisBlockProvisions, validateGenesisBlockProvisions),
		paramtypes.NewParamSetPair(KeyReductionPeriodInBlocks, &p.ReductionPeriodInBlocks, validateReductionPeriodInBlocks),
		paramtypes.NewParamSetPair(KeyReductionFactor, &p.ReductionFactor, validateReductionFactor),
		paramtypes.NewParamSetPair(KeyPoolAllocationRatio, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(KeyDeveloperRewardsReceiver, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramtypes.NewParamSetPair(KeyMintingRewardsDistributionStartBlock, &p.MintingRewardsDistributionStartBlock, validateMintingRewardsDistributionStartBlock),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateGenesisBlockProvisions(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("genesis block provision must be non-negative")
	}

	return nil
}

func validateReductionPeriodInBlocks(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateReductionFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("reduction factor cannot be greater than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("reduction factor cannot be negative")
	}

	return nil
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Staking.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.PoolIncentives.IsNegative() {
		return errors.New("pool incentives distribution ratio should not be negative")
	}

	if v.DeveloperRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}

	if v.CommunityPool.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	totalProportions := v.Staking.Add(v.PoolIncentives).Add(v.DeveloperRewards).Add(v.CommunityPool)

	if !totalProportions.Equal(sdk.NewDec(1)) {
		return errors.New("total distributions ratio should be 1")
	}

	return nil
}

func validateWeightedDeveloperRewardsReceivers(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateMintingRewardsDistributionStartBlock(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("start block must be non-negative")
	}

	return nil
}
