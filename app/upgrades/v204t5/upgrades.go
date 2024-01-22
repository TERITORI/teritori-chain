package v204t5

import (
	"github.com/TERITORI/teritori-chain/app/keepers"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

//nolint:all
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("start to run module migrations...")

		cp := tmtypes.DefaultConsensusParams().ToProto()
		keepers.ConsensusParamsKeeper.Set(ctx, &tmproto.ConsensusParams{
			Block:     cp.Block,
			Validator: cp.Validator,
			Evidence:  cp.Evidence,
			Version:   cp.Version,
		})

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
