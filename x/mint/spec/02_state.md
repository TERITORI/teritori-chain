# State

## Minter

## Params

## LastReductionBlock

Last reduction block stores the block number when the last reduction of
coin mint amount per block has happened.

## NextBlockProvisions

The target block provision is recalculated on each reduction period
(1 year). At the time of the reduction, the current provision is
multiplied by the reduction factor (default `2/3`), to calculate the
provisions for the next block. Consequently, the rewards of the next
period will be lowered by a `1` - reduction factor.

## BlockProvision

Calculate the provisions generated for each block based on current block
provisions. The provisions are then minted by the `mint` module's
`ModuleMinterAccount`. These rewards are transferred to a
`FeeCollector`, which handles distributing the rewards per the chain's needs.
This fee collector is specified as the `auth` module's `FeeCollector` `ModuleAccount`.
