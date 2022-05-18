package cli

import (
	"fmt"

	"github.com/POPSmartContract/nxtpop-chain/x/nftstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

const (
	FlagNftIdentifier = "nft-identifier"
	FlagNftMetadata   = "nft-metadata"
	FlagRewardAddress = "reward-address"
	FlagRewardWeight  = "reward-weight"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nftstaking sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxRegisterNftStakingCmd(),
	)

	return txCmd
}

func GetTxRegisterNftStakingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-nft-staking [flags]",
		Short: "Register a nft staking with nft identifier and reward address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			nftIdentifier, err := cmd.Flags().GetString(FlagNftIdentifier)
			if err != nil {
				return fmt.Errorf("invalid nft identifier: %w", err)
			}

			nftMetadata, err := cmd.Flags().GetString(FlagNftMetadata)
			if err != nil {
				return fmt.Errorf("invalid nft metadata: %w", err)
			}

			rewardAddr, err := cmd.Flags().GetString(FlagRewardAddress)
			if err != nil {
				return fmt.Errorf("invalid reward address: %w", err)
			}

			rewardWeight, err := cmd.Flags().GetUint64(FlagRewardWeight)
			if err != nil {
				return fmt.Errorf("invalid reward address: %w", err)
			}

			msg := types.NewMsgRegisterNftStaking(clientCtx.FromAddress.String(), types.NftStaking{
				NftIdentifier: nftIdentifier,
				NftMetadata:   nftMetadata,
				RewardAddress: rewardAddr,
				RewardWeight:  rewardWeight,
			})
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagNftIdentifier, "", "nft identifier.")
	cmd.Flags().String(FlagNftMetadata, "", "nft metadata.")
	cmd.Flags().String(FlagRewardAddress, "", "Reward address to receive staking rewards.")
	cmd.Flags().Uint64(FlagRewardWeight, 0, "Reward weight for the nft")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
