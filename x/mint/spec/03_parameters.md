# Network Parameters

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
