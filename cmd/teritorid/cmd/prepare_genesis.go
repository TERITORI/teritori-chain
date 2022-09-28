package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	appparams "github.com/TERITORI/teritori-chain/app/params"
	"github.com/TERITORI/teritori-chain/x/airdrop/types"
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
	teritorid prepare-genesis teritori-1 cosmos_aidrop.csv crew3_airdrop.csv evmos_orbital_ape.csv
	- Check input genesis:
		file is at ~/.teritorid/config/genesis.json
`,
		Args: cobra.ExactArgs(4),
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
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, chainID, args[1], args[2], args[3])
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

func combineAirdropAllocations(allocations1, allocations2 []airdroptypes.AirdropAllocation) []types.AirdropAllocation {
	usedAllocation := map[string]bool{}
	allocationIndex := map[string]int{}

	for index, allo := range allocations1 {
		usedAllocation[allo.Address] = true
		allocationIndex[allo.Address] = index
	}

	allocations := allocations1
	for _, allo := range allocations2 {
		if usedAllocation[allo.Address] {
			allocations[allocationIndex[allo.Address]].Amount = allocations[allocationIndex[allo.Address]].Amount.Add(allo.Amount)
		} else {
			allocations = append(allocations, allo)
		}
	}
	return allocations
}

func PrepareGenesis(clientCtx client.Context, appState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, chainID, cosmosAirdropPath, crew3AirdropPath, evmosOrbitalApePath string) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
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
	crew3Allocations, totalCrew3AirdropAllocation := parseCosmosAirdropAmount(crew3AirdropPath)
	cosmosAllocations = combineAirdropAllocations(cosmosAllocations, crew3Allocations)
	totalCosmosAirdropAllocation = totalCosmosAirdropAllocation.Add(totalCrew3AirdropAllocation)
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

	// send 10 TORI to genesis validators
	genesisValidators := []string{
		"tori1534tslwra4hrvt8k8tdwh5aghmc74hvtl7xfnc", // gopher
		"tori1uechzauku6mhj2je8jmyrkq6d0ydm3g4q470r0", // metahuahua
		"tori1t7cyvydpp4lklprksnrjy2y3xzv3q2l0n4qqvn", // activenodes
		"tori1utr8j9685hfxyza3wnu8pa9lpu8360knym0vjc", // alxvoy
		"tori17vx59a9897ltpyw6dwr7jvcjk8wyxhc08svtzj", // aurie
		"tori19pg5t5adjese84q5azjuv46rtz7jt2yqrrfueg", // berty
		"tori16czx9ukfcelsxmjpyt90fprjx5qjw5anhpr7vk", // chillvalidation
		"tori1c22uwrtvadcp2a8rjn2l00kmuuqdcu2t3qttty", // crosnest
		"tori1jfw63tylcc4gscayv68prsu275q6te4wmpwmp7", // dibugnodes
		"tori1a7taydvzhkd5vrndlykqtj7nsk2erdp2c06ljz", // ericet
		"tori15n624eajd04jjhnlvza2fvft3lmf69aed78yde", // interblockchainservices
		"tori1gspqsgxm8s4e9uza78met67x2eted2cdcfqefg", // goldenratiostaking
		"tori1ttgzvn4lwkqe33drcvjxrefu8j9u5x9qvf22vc", // hashquark
		"tori1lgy98shrs4uyrqnmgh38su3gm08uh3srhsmlzf", // ibrahimarslan
		"tori1xpyql3vw67h8l99n3sswy5ev94ntwt9ce8fdj9", // icosmosdao
		"tori1shtyw4f5pdhvx7gsrsknwrryy9ulqvvyjflv2j", // kalianetwork
		"tori1rly8ah6hffkt28hy3ka8ue2h32mqknyxka5tgk", // landeros
		"tori140l6y2gp3gxvay6qtn70re7z2s0gn57zwdyda3", // lavenderfive
		"tori1kunzrdg6u8gql4faj33lstghhqdtp59e25th6f", // lesnik_utsa
		"tori1e6ajryqxefpxuhjg2y9wk4y2dzq48uz4gvuy4q", // maethstro
		"tori18wjuryzyuwpg5f0wukgjey3za28s4fm9mujyrh", // munris
		"tori1p5z27dj7zrxue8pe5t0m39q8mmgavdclwerqqw", // n0ok
		"tori18t2j2kc08su2l2dafcanq43yxj9akpwpl66y3z", // nodejumper
		"tori1nrgahzmlr4nrnumlu0ud99qslsdvay8ah5k6c8", // noderunners
		"tori1phzay7cf4ayk9dsvt0q5nlc8qehlwlpxs3p957", // nodesblocks
		"tori1w3wse8cx2al5947ke0hnd2tgphjt43dyq6dvay", // nodesguru
		"tori1nuh2h60wlvzvk58xll3d8gz2wpqjt6gw8tgsc0", // nysa_network
		"tori1azdfljp04ptlazs95e5gscweavmaszw58fwpth", // oni
		"tori18je2ph09a7flemkkzmvenz2eeyw5pdge93gxyt", // orbitalapes
		"tori1gp957czryfgyvxwn3tfnyy2f0t9g2p4p8e683y", // polkachu
		"tori1fyyl63zqylda0qrkqdzeyag28eyh9swrkusuls", // rhino
		"tori1qy38xmcrnht0kt5c5fryvl8llrpdwer6ce9dx3", // romanv
		"tori1267l9z6yeua438mct5ee2mnm53yn3n9wlqvdl7", // samourai-world
		"tori167xwmhtrn7n8ftexu6luhh4luvhpy435nuxmma", // silknodes
		"tori1xu736l4vt6l2pg9k2yk66fq7zq6y4aj5rfw97d", // stakelab
		"tori1cm3hmw63a9wugawf9jn2jv0savynkgu9ra9ggm", // stakingcabin
		"tori1sqk72uwf6tg867ssuu7whxfu9pfcyrpe9u76c4", // stavr
		"tori1lrq8sl2jq7246yjplutv5lul8ykrhqcr9frjz3", // stingray
		"tori1gtz5v838vf7ucnn0jnqr3crs5099g9p2fnpe9p", // teritori-core-1
		"tori1vxmq5epj83z8en5h0zul624nrmfxzmhkmwmmtl", // teritori-core-2
		"tori1x6vfjy754fvzrlug2kxsp6s54yfj753sheqpay", // web34ever
		"tori1dfnzup7nppxvlpwnmzjnuet0tn4t9cnwhqx54s", // wetez
		"tori1tjh6wpj6d9kpkfrcyglksevkhhtk9gm7auaxy3", // whispernode
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
