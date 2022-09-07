# Concepts

The `x/mint` module is designed to handle the regular printing of new tokens within a chain. The design taken is to

- Mint new tokens once per block
- To have a "Reductioning factor" every period, which reduces the number of rewards per block. (default: period is 1 year. The next period's rewards are 2/3 of the prior period's rewards)

## Reduction factor

This is a generalization over the Bitcoin-style halvenings. Every year, the number of rewards issued per week will reduce by a governance-specified factor, instead of a fixed `1/2`. So `RewardsPerBlockNextPeriod = ReductionFactor * CurrentRewardsPerBlock)`.
When `ReductionFactor = 1/2`, the Bitcoin halvenings are recreated. We default to having a reduction factor of `2/3` and thus reduce rewards at the end of every year by `33%`.

The implication of this is that the total supply is finite, according to the following formula:

`Total Supply = InitialSupply + BlocksPerPeriod * { {InitialRewardsPerBlock} / {1 - ReductionFactor} }`
