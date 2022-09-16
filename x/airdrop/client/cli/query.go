package cli

import (
	"context"
	"fmt"

	"github.com/TERITORI/teritori-chain/x/airdrop/types"
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
		GetCmdQueryParams(),
	)

	return queryCmd
}

func GetCmdQueryAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocation [addr]",
		Short: "Query allocations of an address",
		Args:  cobra.ExactArgs(1),
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

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query params",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryParamsRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
