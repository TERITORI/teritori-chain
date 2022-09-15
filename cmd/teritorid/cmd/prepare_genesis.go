package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	appparams "github.com/TERITORI/teritori-chain/app/params"
	airdroptypes "github.com/TERITORI/teritori-chain/x/airdrop/types"
	minttypes "github.com/TERITORI/teritori-chain/x/mint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

// PrepareGenesisCmd returns prepare-genesis cobra Command.
func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis [chain_id] [airdrop_file_path]",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
Example:
	teritorid prepare-genesis teritori-1 cosmos_aidrop.csv
	- Check input genesis:
		file is at ~/.teritorid/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.Codec
			cdc := depCdc
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			chainID := args[0]

			// run Prepare Genesis
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, chainID, args[1])
			if err != nil {
				return err
			}

			// validate genesis state
			if err = mbm.ValidateGenesis(cdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func parseAirdropAmount(path string) ([]airdroptypes.AirdropAllocation, sdk.Coin) {
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

	totalAmount := sdk.ZeroInt()
	allocations := []airdroptypes.AirdropAllocation{}
	for _, line := range records[1:] {
		cosmosAddr, amountStr := line[0], line[1]
		amountDec := sdk.MustNewDecFromStr(amountStr)
		amountInt := amountDec.Mul(sdk.NewDec(1000_000)).TruncateInt()

		allocations = append(allocations, airdroptypes.AirdropAllocation{
			Chain:         "cosmos",
			Address:       cosmosAddr,
			Amount:        sdk.NewCoin(appparams.BaseCoinUnit, amountInt),
			ClaimedAmount: sdk.NewInt64Coin(appparams.BaseCoinUnit, 0),
		})
		totalAmount = totalAmount.Add(amountInt)
	}

	return allocations, sdk.NewCoin(appparams.BaseCoinUnit, totalAmount)
}

func PrepareGenesis(clientCtx client.Context, appState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, chainID, path string) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	depCdc := clientCtx.Codec
	cdc := depCdc

	// chain params genesis
	genDoc.ChainID = chainID

	// mint module genesis
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = minttypes.DefaultParams()
	mintGenState.Params.MintDenom = appparams.BaseCoinUnit
	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// airdrop module genesis
	airdropGenState := airdroptypes.DefaultGenesis()
	airdropGenState.Params = airdroptypes.DefaultParams()
	allocations, totalAirdropAllocation := parseAirdropAmount(path)
	airdropGenState.Allocations = allocations
	airdropGenStateBz, err := cdc.MarshalJSON(airdropGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal airdrop genesis state: %w", err)
	}
	appState[airdroptypes.ModuleName] = airdropGenStateBz

	// bank module genesis
	bankGenState := banktypes.DefaultGenesisState()
	bankGenState.Params = banktypes.DefaultParams()

	bankGenState.Supply = sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 200_000_000_000_000)) // 200M TORI

	airdropCoins := sdk.Coins{totalAirdropAllocation}
	communityPoolCoins := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 50_000_000_000_000)) // 50M TORI

	// airdrop balances
	airdropModuleAddr := authtypes.NewModuleAddress(airdroptypes.ModuleName)
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: airdropModuleAddr.String(),
		Coins:   airdropCoins,
	})
	// distribution balances for community pool
	distrModuleAddr := authtypes.NewModuleAddress(distributiontypes.ModuleName)
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: distrModuleAddr.String(),
		Coins:   communityPoolCoins,
	})

	// strategic reserve = 200M - 50M - airdropCoins
	addrStrategicReserve := sdk.AccAddress("strategic_reserve") // TODO: update
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: addrStrategicReserve.String(),
		Coins:   bankGenState.Supply.Sub(airdropCoins).Sub(communityPoolCoins),
	})

	// TODO: send 10 TORI to genesis validators
	// TODO: send 1 TORI to cosmos airdrop addresses if required

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	// account module genesis
	authGenState := authtypes.GetGenesisStateFromAppState(depCdc, appState)
	authGenState.Params = authtypes.DefaultParams()

	genAccounts := []authtypes.GenesisAccount{authtypes.NewBaseAccount(addrStrategicReserve, nil, 0, 0)}
	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		panic(err)
	}

	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[authtypes.ModuleName] = authGenStateBz

	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(depCdc, appState)
	stakingGenState.Params = stakingtypes.DefaultParams()
	stakingGenState.Params.BondDenom = appparams.BaseCoinUnit
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = distributiontypes.DefaultParams()
	distributionGenState.FeePool.CommunityPool = sdk.NewDecCoinsFromCoins(communityPoolCoins...)
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// gov module genesis
	govGenState := govtypes.DefaultGenesisState()
	defaultGovParams := govtypes.DefaultParams()
	govGenState.DepositParams = defaultGovParams.DepositParams
	govGenState.TallyParams = defaultGovParams.TallyParams
	govGenState.VotingParams = defaultGovParams.VotingParams
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = slashingtypes.DefaultParams()
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// return appState and genDoc
	return appState, genDoc, nil
}
