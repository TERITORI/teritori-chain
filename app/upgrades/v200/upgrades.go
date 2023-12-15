package v200

import (
	"github.com/TERITORI/teritori-chain/app/keepers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("start to run module migrations...")

		// Mint module params update
		params := keepers.MintKeeper.GetParams(ctx)
		params.BlocksPerYear = 5733818
		params.TotalBurntAmount = []sdk.Coin{sdk.NewInt64Coin("utori", 118550_000000)}
		keepers.MintKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
