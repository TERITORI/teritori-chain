package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/NXTPOP/teritori-chain/x/nftstaking/types"
)

func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the nftstaking module",
	}
	queryCmd.AddCommand(
		GetCmdQueryNftStakings(),
	)

	return queryCmd
}

func GetCmdQueryNftStakings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stakings",
		Short: "Get all nft stakings",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryNftStakings(context.Background(), &types.QueryNftStakingsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
