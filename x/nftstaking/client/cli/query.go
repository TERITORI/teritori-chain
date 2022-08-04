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
		GetCmdQueryAccessInfos(),
		GetCmdQueryAccessInfo(),
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

func GetCmdQueryAccessInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-infos",
		Short: "Get all access information",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAccessInfos(context.Background(), &types.QueryAccessInfosRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAccessInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-info [address]",
		Short: "Get an access information for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAccessInfo(context.Background(), &types.QueryAccessInfoRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
