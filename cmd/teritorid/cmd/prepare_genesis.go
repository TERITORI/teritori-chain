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
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	liquiditytypes "github.com/gravity-devs/liquidity/x/liquidity/types"
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
	teritorid prepare-genesis teritori-1 cosmos_aidrop.csv evmos_orbital_ape.csv
	- Check input genesis:
		file is at ~/.teritorid/config/genesis.json
`,
		Args: cobra.ExactArgs(3),
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
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, chainID, args[1], args[2])
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

func parseCosmosAirdropAmount(path string) ([]airdroptypes.AirdropAllocation, sdk.Coin) {
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

func parseEvmosOrbitalApeAirdropAmount(path string) ([]airdroptypes.AirdropAllocation, sdk.Coin) {
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
		evmAddr, amountStr := line[0], line[1]
		amountDec := sdk.MustNewDecFromStr(amountStr)
		amountInt := amountDec.Mul(sdk.NewDec(1000_000)).TruncateInt()

		allocations = append(allocations, airdroptypes.AirdropAllocation{
			Chain:         "evm",
			Address:       evmAddr,
			Amount:        sdk.NewCoin(appparams.BaseCoinUnit, amountInt),
			ClaimedAmount: sdk.NewInt64Coin(appparams.BaseCoinUnit, 0),
		})
		totalAmount = totalAmount.Add(amountInt)
	}

	return allocations, sdk.NewCoin(appparams.BaseCoinUnit, totalAmount)
}

func PrepareGenesis(clientCtx client.Context, appState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, chainID, cosmosAirdropPath, evmosOrbitalApePath string) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	depCdc := clientCtx.Codec
	cdc := depCdc

	// chain params genesis
	genDoc.ChainID = chainID

	// mint module genesis
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = minttypes.DefaultParams()
	mintGenState.Params.MintDenom = appparams.BaseCoinUnit
	mintGenState.Params.MintingRewardsDistributionStartBlock = 51840 // 3 days after launch - 86400s x 3 / 5s
	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// airdrop module genesis
	airdropGenState := airdroptypes.DefaultGenesis()
	airdropGenState.Params = airdroptypes.DefaultParams()
	airdropGenState.Params.Owner = "tori19ftk3lkfupgtnh38d7enc8c6jp7aljj3jmknnm" // POP's address
	cosmosAllocations, totalCosmosAirdropAllocation := parseCosmosAirdropAmount(cosmosAirdropPath)
	evmosOrbitalApeAllocations, totalEvmosAirdropAllocataion := parseEvmosOrbitalApeAirdropAmount(evmosOrbitalApePath)
	allocations := append(cosmosAllocations, evmosOrbitalApeAllocations...)
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

	airdropCoins := sdk.Coins{totalCosmosAirdropAllocation.Add(totalEvmosAirdropAllocataion)}
	communityPoolCoins := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 50_000_000_000_000)) // 50M TORI

	seenBalances := make(map[string]bool)

	moduleAccs := []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		icatypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		govtypes.ModuleName,
		liquiditytypes.ModuleName,
		ibctransfertypes.ModuleName,
		airdroptypes.ModuleName,
	}

	for _, module := range moduleAccs {
		moduleAddr := authtypes.NewModuleAddress(module)
		seenBalances[moduleAddr.String()] = true
	}

	// airdrop balance
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

	genAccounts := []authtypes.GenesisAccount{}

	addrStrategicReserve, err := sdk.AccAddressFromBech32("tori1kcuty7d5mc0rasw6mpmn4nhk99me55cheldr23")
	if err != nil {
		return nil, nil, err
	}
	genAccounts = append(genAccounts, authtypes.NewBaseAccount(addrStrategicReserve, nil, 0, 0))

	// TODO: send 10 TORI to genesis validators
	genesisValidators := []string{
		"tori1534tslwra4hrvt8k8tdwh5aghmc74hvtl7xfnc", // gopher
	}

	totalValidatorInitialCoins := sdk.NewCoins()
	validatorInitialCoins := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 10_000_000)) // 10 TORI
	for _, address := range genesisValidators {
		if seenBalances[address] {
			continue
		}

		bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
			Address: address,
			Coins:   validatorInitialCoins, // 0.1 TORI
		})
		addr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return nil, nil, err
		}
		totalValidatorInitialCoins = totalValidatorInitialCoins.Add(validatorInitialCoins...)
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(addr, nil, 0, 0))

		seenBalances[address] = true
	}

	// send 0.1 TORI to bech32 converted cosmos airdrop addresses
	totalAirdropGasCoins := sdk.NewCoins()
	airdropGasCoins := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 100_000))
	for _, allocation := range cosmosAllocations {
		_, bz, err := bech32.DecodeAndConvert(allocation.Address)
		if err != nil {
			return nil, nil, err
		}

		bech32Addr, err := bech32.ConvertAndEncode(appparams.Bech32PrefixAccAddr, bz)
		if err != nil {
			return nil, nil, err
		}

		if seenBalances[bech32Addr] {
			continue
		}

		bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
			Address: bech32Addr,
			Coins:   airdropGasCoins, // 0.1 TORI
		})
		totalAirdropGasCoins = totalAirdropGasCoins.Add(airdropGasCoins...)
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(sdk.AccAddress(bz), nil, 0, 0))

		seenBalances[bech32Addr] = true
	}

	// strategic reserve = 200M - 50M - airdropCoins
	bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
		Address: addrStrategicReserve.String(),
		Coins:   bankGenState.Supply.Sub(airdropCoins).Sub(communityPoolCoins).Sub(totalAirdropGasCoins).Sub(totalValidatorInitialCoins),
	})

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	// account module genesis
	authGenState := authtypes.GetGenesisStateFromAppState(depCdc, appState)
	authGenState.Params = authtypes.DefaultParams()

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
	govGenState.DepositParams.MinDeposit = sdk.Coins{sdk.NewInt64Coin(appparams.BaseCoinUnit, 5000_000_000)} // 5000 TORI
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
	slashingGenState.Params.SignedBlocksWindow = 10000
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = sdk.NewInt64Coin(appparams.BaseCoinUnit, 1000_000) // 1 TORI
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// return appState and genDoc
	return appState, genDoc, nil
}
