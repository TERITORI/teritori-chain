package cli

import (
	"fmt"

	"github.com/TERITORI/teritori-chain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxClaimAllocationCmd(),
		GetTxSetAllocationCmd(),
		GetTxDepositTokensCmd(),
		AllocateFurtherAirdropCmd(),
	)

	return txCmd
}

// GetTxClaimAllocationCmd implement cli command for MsgClaimAllocation
func GetTxClaimAllocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-allocation [native_chain_address] [signature]",
		Short: "Claim reward allocation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgClaimAllocation(
				args[0],
				clientCtx.FromAddress,
				args[1],
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxSetAllocationCmd implement cli command for MsgSetAllocation
func GetTxSetAllocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-allocation [chain] [native_chain_address] [amount] [claimed_amount]",
		Short: "Set allocation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			claimedAmount, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAllocation(
				clientCtx.FromAddress.String(),
				types.AirdropAllocation{
					Chain:         args[0],
					Address:       args[1],
					Amount:        amount,
					ClaimedAmount: claimedAmount,
				},
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxTransferModuleOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer-module-ownership [newOwner] [flags]",
		Long: "Transfer module ownership to new address",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferModuleOwnership(
				clientCtx.GetFromAddress(),
				args[0],
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetTxDepositTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deposit-tokens [amount] [flags]",
		Long: "Deposit tokens to airdrop module",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgDepositTokens(
				clientCtx.GetFromAddress(),
				amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
