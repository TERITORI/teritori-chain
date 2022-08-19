# Mint

The `mint` module is responsible for creating tokens in a flexible way to reward
validators, incentivize providing pool liquidity, provide funds for governance,
and pay developers to maintain.

The module is also responsible for reducing the token creation and distribution by a set period
until it reaches its maximum supply (see `reduction_factor` and `reduction_period_in_blocks`)

The module uses time basis blocks supported by the `blocks` module.

## Contents

1. **[Concept](#concepts)**
2. **[State](#state)**
3. **[Parameters](#network-parameters)**
4. **[Events](#events)**
5. **[Transactions](#transaction)**
6. **[Queries](#queries)**

## Concepts

The `x/mint` module is designed to handle the regular printing of new
tokens within a chain. The design taken is to

- Mint new tokens once per block
- To have a "Reductioning factor" every period, which reduces the number of
  rewards per block. (default: period is 1 year. The next period's rewards are 2/3 of the prior
  period's rewards)

### Reduction factor

This is a generalization over the Bitcoin-style halvenings. Every year, the number
of rewards issued per week will reduce by a governance-specified
factor, instead of a fixed `1/2`. So
`RewardsPerBlockNextPeriod = ReductionFactor * CurrentRewardsPerBlock)`.
When `ReductionFactor = 1/2`, the Bitcoin halvenings are recreated. We
default to having a reduction factor of `2/3` and thus reduce rewards
at the end of every year by `33%`.

The implication of this is that the total supply is finite, according to
the following formula:

`Total Supply = InitialSupply + BlocksPerPeriod * { {InitialRewardsPerBlock} / {1 - ReductionFactor} }`

## State

### Minter

### Params

### LastReductionBlock

Last reduction block stores the block number when the last reduction of
coin mint amount per block has happened.

### NextBlockProvisions

The target block provision is recalculated on each reduction period
(1 year). At the time of the reduction, the current provision is
multiplied by the reduction factor (default `2/3`), to calculate the
provisions for the next block. Consequently, the rewards of the next
period will be lowered by a `1` - reduction factor.

### BlockProvision

Calculate the provisions generated for each block based on current block
provisions. The provisions are then minted by the `mint` module's
`ModuleMinterAccount`. These rewards are transferred to a
`FeeCollector`, which handles distributing the rewards per the chain's needs.
This fee collector is specified as the `auth` module's `FeeCollector` `ModuleAccount`.

## Network Parameters

// distribution_proportions defines the proportion of the minted denom
DistributionProportions distribution_proportions = 5 [ (gogoproto.nullable) = false ];
// address to receive developer rewards
repeated WeightedAddress weighted_developer_rewards_receivers = 6 [(gogoproto.nullable) = false];
// usage incentive address
string usage_incentive_address = 7;
// grants program address
string grants_program_address = 8;
// team reserve funds address
string team_reserve_address = 9;
// start block to distribute minting rewards
int64 minting_rewards_distribution_start_block = 10;
}

The minting module contains the following parameters:

| Key                                        | Type         | Example                                |
| ------------------------------------------ | ------------ | -------------------------------------- |
| mint_denom                                 | string       | "utori"                                |
| genesis_block_provisions                   | string (dec) | "500000000"                            |
| reduction_period_in_blocks                 | int64        | 156                                    |
| reduction_factor                           | string (dec) | "0.6666666666666"                      |
| distribution_proportions.grants_program    | string (dec) | "0.4"                                  |
| distribution_proportions.community_pool    | string (dec) | "0.3"                                  |
| distribution_proportions.usage_incentive   | string (dec) | "0.2"                                  |
| distribution_proportions.staking           | string (dec) | "0.1"                                  |
| distribution_proportions.developer_rewards | string (dec) | "0.1"                                  |
| weighted_developer_rewards_receivers       | array        | [{"address": "torixx", "weight": "1"}] |
| usage_incentive_address                    | string       | "torixx"                               |
| grants_program_address                     | string       | "torixx"                               |
| team_reserve_address                       | string       | "torixx"                               |
| minting_rewards_distribution_start_block   | int64        | 10                                     |

Below are all the network parameters for the `mint` module:

- **`mint_denom`** - Token type being minted
- **`genesis_block_provisions`** - Amount of tokens generated at the block to the distribution categories (see distribution_proportions)
- **`reduction_period_in_blocks`** - How many blocks must occur before implementing the reduction factor
- **`reduction_factor`** - What the total token issuance factor will reduce by after the reduction period passes (if set to 66.66%, token issuance will reduce by 1/3)
- **`distribution_proportions`** - Categories in which the specified proportion of newly released tokens are distributed to
  - **`grants_program`** - Proportion of minted funds to grants program account
  - **`community_pool`** - Proportion of minted funds to be set aside for the community pool
  - **`usage_incentive`** - Proportion of minted funds to usage incentive account
  - **`staking`** - Proportion of minted funds to pay staking incentive
  - **`developer_rewards`** - Proportion of minted funds to pay developers for their past and future work
- **`grants_program_address`** - Address to receive gran program tokens
- **`usage_incentive_address`** - Address to receive usage incentive
- **`team_reserve_address`** - Address to receive team reserve tokens
- **`weighted_developer_rewards_receivers`** - Addresses that developer rewards will go to. The weight attached to an address is the percent of the developer rewards that the specific address will receive
- **`minting_rewards_distribution_start_block`** - What block will start the rewards distribution to the aforementioned distribution categories

**Notes**

1. `mint_denom` defines denom for minting token - utori
2. `genesis_block_provisions` provides minting tokens per block at genesis.
3. `reduction_period_in_blocks` defines the number of blocks to pass to reduce the mint amount
4. `reduction_factor` defines the reduction factor of tokens at every `reduction_period_in_blocks`
5. `distribution_proportions` defines distribution rules for minted tokens, when the developer
   rewards address is empty, it distributes tokens to the community pool.
6. `weighted_developer_rewards_receivers` provides the addresses that receive developer
   rewards by weight
7. `minting_rewards_distribution_start_block` defines the start block of minting to make sure
   minting start after initial pools are set

## Queries

### params

Query all the current mint parameter values

```sh
query mint params
```
