package cli

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	appparams "github.com/TERITORI/teritori-chain/app/params"
	airdroptypes "github.com/TERITORI/teritori-chain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func parseEvmosOrbitalApeAirdropAmount(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}

func saveAllocation(path string, allocations []airdroptypes.AirdropAllocation) {
	records := [][]string{{"address", "amount", "claimed"}}
	for _, allocation := range allocations {
		decimal := sdk.NewDec(1000_000)
		records = append(records, []string{
			allocation.Address,
			sdk.NewDecFromInt(allocation.Amount.Amount).Quo(decimal).String(),
			sdk.NewDecFromInt(allocation.ClaimedAmount.Amount).Quo(decimal).String(),
		})
	}

	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			panic(err)
		}
	}
}

// FetchAndRemoveAirdropCmd removes airdrop allocation and stores the result
func FetchAndRemoveAirdropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch-and-remove-airdrop [airdrop_file_path] [onchain_result_store_path] [start_index] [msgs_per_tx]",
		Short: "Fetch and remove airdrop allocation",
		Long: `Fetch and remove airdrop allocation.
Example:
	teritorid tx airdrop fetch-and-remove-airdrop evmos_orbital_ape.csv stored_result.csv 0 500 --from=validator --keyring-backend=test --chain-id=teritori-1 --home=$HOME/.teritorid/ --yes --broadcast-mode=block --gas=10000000
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocations := []airdroptypes.AirdropAllocation{}
			allocationRecords := parseEvmosOrbitalApeAirdropAmount(args[0])
			for _, line := range allocationRecords[1:] {
				airdropAddr, amountStr := line[0], line[1]
				amountDec := sdk.MustNewDecFromStr(amountStr)
				amount := amountDec.Mul(sdk.NewDec(1000_000)).TruncateInt()

				params := &airdroptypes.QueryAllocationRequest{Address: airdropAddr}

				allocation := airdroptypes.AirdropAllocation{
					Chain:         "evmos",
					Address:       airdropAddr,
					Amount:        sdk.NewCoin(appparams.BaseCoinUnit, amount),
					ClaimedAmount: sdk.NewInt64Coin(appparams.BaseCoinUnit, 0),
				}

				// query allocation
				queryClient := airdroptypes.NewQueryClient(clientCtx)
				res, err := queryClient.Allocation(context.Background(), params)
				if err == nil && res.Allocation != nil {
					allocation = *res.Allocation
				}

				allocations = append(allocations, allocation)
			}

			// save allocations as csv file
			saveAllocation(args[1], allocations)

			startIndex, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			threshold, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{}
			for index, allocation := range allocations {
				if index < startIndex {
					continue
				}
				if allocation.ClaimedAmount.Equal(allocation.Amount) {
					continue
				}
				allocation.ClaimedAmount = allocation.Amount
				msg := airdroptypes.NewMsgSetAllocation(
					clientCtx.FromAddress.String(),
					allocation,
				)
				msgs = append(msgs, msg)
				if len(msgs) >= threshold || index+1 == len(allocations) {
					err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
					if err != nil {
						return err
					}
					fmt.Println("executed until index", index)
					msgs = []sdk.Msg{}
				}
			}

			fmt.Println("finalized execution of setting airdrop as claimed")

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
