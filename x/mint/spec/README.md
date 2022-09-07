# Mint

The `mint` module is responsible for creating tokens in a flexible way to reward
validators, incentivize providing pool liquidity, provide funds for governance,
and pay developers to maintain.

The module is also responsible for reducing the token creation and distribution by a set period
until it reaches its maximum supply (see `reduction_factor` and `reduction_period_in_blocks`)

The module uses time basis blocks supported by the `blocks` module.

## Contents

1. **[Concept](01_concept.md)**
2. **[State](02_state.md)**
3. **[Parameters](03_parameters.md)**
4. **[Events](04_events.md)**
5. **[Queries](05_queries.md)**
