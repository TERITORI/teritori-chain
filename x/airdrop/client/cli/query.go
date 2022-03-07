package cli

import (
	"context"
	"fmt"

	"github.com/POPSmartContract/nxtpop-chain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryAllocation(),
	)

	return queryCmd
}

func GetCmdQueryAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions [addr]",
		Short: "Query allocations of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryAllocationRequest{Address: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Allocation(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Allocation)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
