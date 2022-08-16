package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/TERITORI/teritori-chain/x/nftstaking/types"
)

func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the nftstaking module",
	}
	queryCmd.AddCommand(
		GetCmdQueryNftStaking(),
		GetCmdQueryNftStakingsByOwner(),
		GetCmdQueryNftStakings(),
		GetCmdQueryAccessInfos(),
		GetCmdQueryAccessInfo(),
		GetCmdQueryAllNftTypePerms(),
		GetCmdQueryNftTypePerms(),
		GetCmdQueryHasPermission(),
	)

	return queryCmd
}

func GetCmdQueryNftStaking() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking [identifier]",
		Short: "Get an nft staking by identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryNftStaking(context.Background(), &types.QueryNftStakingRequest{
				Identifier: args[0],
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

func GetCmdQueryNftStakingsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stakings_by_owner [owner]",
		Short: "Get nft stakings by owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryNftStakingsByOwner(context.Background(), &types.QueryNftStakingsByOwnerRequest{
				Owner: args[0],
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

func GetCmdQueryAllNftTypePerms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-nfttype-perms [address]",
		Short: "Get all permissions allocated to nfts by type",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAllNftTypePerms(context.Background(), &types.QueryAllNftTypePermsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryNftTypePerms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfttype-perms [nft_type]",
		Short: "Get permissions allocated to a nft type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryNftTypePerms(context.Background(), &types.QueryNftTypePermsRequest{
				NftType: types.NftType(types.NftType_value[args[0]]),
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

func GetCmdQueryHasPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "has-permission [address] [permission]",
		Short: "Get if an address has permission",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryHasPermission(context.Background(), &types.QueryHasPermissionRequest{
				Address:    args[0],
				Permission: args[1],
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
