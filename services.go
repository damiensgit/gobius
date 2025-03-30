package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/bindings/arbiusrouterv1"
	"gobius/bindings/basetoken"
	"gobius/bindings/bulktasks"
	"gobius/bindings/engine"
	"gobius/bindings/voter"
	"gobius/client"
	"gobius/config"
	"gobius/erc20"
	"gobius/ipfs"
	"gobius/metrics"
	"gobius/paraswap"
	"gobius/storage"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type servicesKey struct{}

type Services struct {
	Basetoken          *basetoken.BaseToken
	Engine             *EngineWrapper
	Voter              *voter.Voter
	ArbiusRouter       ipfs.ArbiusRouterContract
	BulkTasks          *bulktasks.BulkTasks
	Eth                *erc20.TokenERC20
	OwnerAccount       *account.Account
	SenderOwnerAccount *account.Account
	Clients            []*client.Client
	Config             *config.AppConfig
	Logger             *zerolog.Logger
	TaskStorage        *storage.TaskStorageDB
	AutoMineParams     *SubmitTaskParams
	Paraswap           *paraswap.ParaswapManager
	TaskTracker        *metrics.TaskTracker
	IpfsOracle         ipfs.OracleClient
}

func NewApplicationContext(rpc *client.Client, senderrpc *client.Client, clients []*client.Client, sql *sql.DB, logger *zerolog.Logger, cfg *config.AppConfig, ipfsOracle ipfs.OracleClient, appContext, appQuit context.Context) (context.Context, error) {

	ownerAccount, err := account.NewAccount(cfg.Blockchain.PrivateKey, rpc, appContext, cfg.Blockchain.CacheNonce)
	if err != nil {
		return nil, err
	}

	senderOwnerAccount, err := account.NewAccount(cfg.Blockchain.PrivateKey, senderrpc, appContext, cfg.Blockchain.CacheNonce)
	if err != nil {
		return nil, err
	}

	// TODO: need a cleaner way to handle this nonce update on first use - maybe move to NewAccount
	senderOwnerAccount.UpdateNonce()

	baseTokenContract, err := basetoken.NewBaseToken(cfg.BaseConfig.BaseTokenAddress, senderrpc.Client)
	if err != nil {
		return nil, err
	}

	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpc.Client)
	if err != nil {
		return nil, err
	}

	voterContract, err := voter.NewVoter(cfg.BaseConfig.VoterAddress, rpc.Client)
	if err != nil {
		return nil, err
	}

	// if we are on sepolia or mainnet we can use a real contract otherwise we mock (determined by config)
	var arbiusRouter ipfs.ArbiusRouterContract
	if cfg.BaseConfig.ArbiusRouterAddress == (common.Address{}) {
		arbiusRouter = ipfs.NewMockRouterContract()
	} else {
		arbiusRouterContract, err := arbiusrouterv1.NewArbiusRouterV1(cfg.BaseConfig.ArbiusRouterAddress, rpc.Client)
		if err != nil {
			return nil, err
		}
		arbiusRouter = ipfs.NewRouterContractWrapper(arbiusRouterContract, rpc)
	}

	var bulkTasksContract *bulktasks.BulkTasks

	if cfg.BaseConfig.TestnetType > 0 {
		// TODO: this is a hack to check if the contract exists on testnet
		if err := checkContractExists(cfg.BaseConfig.BulkTasksAddress, senderrpc.Client, logger); err != nil {

			// on local testnet we can deploy the contract
			// TODO: this is effectively disabled for now until we have a better way to handle it / rewrite contract
			if cfg.BaseConfig.TestnetType > 0 {
				ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				// get the gas price
				gasPrice, err := ownerAccount.Client.GetGasPrice()
				if err != nil {
					return nil, err
				}

				opts := ownerAccount.GetOpts(1_000_000, gasPrice, nil, nil)
				// deploy the contract
				var bulkTasksContractAddress common.Address
				var tx *types.Transaction

				logger.Info().Msg("deploying BulkTasks contract")

				bulkTasksContractAddress, tx, bulkTasksContract, err = bulktasks.DeployBulkTasks(opts, senderrpc.Client, cfg.BaseConfig.BaseTokenAddress, cfg.BaseConfig.EngineAddress)
				if err != nil {
					return nil, err
				}

				logger.Info().Msg("waiting for BulkTasks contract deployment")

				receipt, err := bind.WaitMined(ctxTimeout, senderrpc.Client, tx)
				if err != nil {
					return nil, err
				}

				if receipt.Status != types.ReceiptStatusSuccessful {
					return nil, fmt.Errorf("BulkTasks contract deployment failed")
				}

				logger.Info().Msgf("BulkTasks contract deployed at address %s (you must update the config with this address)", bulkTasksContractAddress.Hex())

				cfg.BaseConfig.BulkTasksAddress = bulkTasksContractAddress
			} else {
				return nil, fmt.Errorf("BulkTasks contract not deployed or invalid at address %s: %w",
					cfg.BaseConfig.BulkTasksAddress.Hex(), err)
			}
		}
	}

	if bulkTasksContract == nil {
		bulkTasksContract, err = bulktasks.NewBulkTasks(cfg.BaseConfig.BulkTasksAddress, senderrpc.Client)
		if err != nil {
			return nil, err
		}
	}

	// cache the min claim time
	minClaimSolTimeBig, err := engineContract.MinClaimSolutionTime(nil)
	if err != nil {
		return nil, err
	}

	// A convienience wrapper to represent ETH
	eth := erc20.NewTokenERC20(common.HexToAddress("0x0"), 18, "ETH", "ETH")

	// 120 = jitter offset in seconds for the min claim time
	minclaimTime := time.Duration(minClaimSolTimeBig.Uint64()+120) * time.Second

	logger.Info().Msgf("Minimum claim time is %s", minclaimTime)

	ts := storage.NewTaskStorageDB(appContext, sql, minclaimTime, logger)

	engineWrapper := NewEngineWrapper(engineContract, voterContract, logger)

	paraswapManager := paraswap.NewParaswapManager(
		senderOwnerAccount,
		baseTokenContract,
		cfg.BaseConfig.BaseToken,
		logger)

	taskMetrics := metrics.NewTaskTracker(appQuit)

	input, err := json.Marshal(cfg.Strategies.Automine.Input)
	if err != nil {
		return nil, err
	}
	st := &SubmitTaskParams{
		Version: uint8(cfg.Strategies.Automine.Version),
		Owner:   ownerAccount.Address,
		Model:   cfg.Strategies.Automine.ModelAsBytes,
		Fee:     cfg.Strategies.Automine.Fee,
		Input:   input,
	}

	services := &Services{
		Basetoken:          baseTokenContract,
		Engine:             engineWrapper,
		Voter:              voterContract,
		BulkTasks:          bulkTasksContract,
		Eth:                eth,
		OwnerAccount:       ownerAccount,
		SenderOwnerAccount: senderOwnerAccount,
		Clients:            clients,
		Config:             cfg,
		Logger:             logger,
		TaskStorage:        ts,
		AutoMineParams:     st,
		Paraswap:           paraswapManager,
		TaskTracker:        taskMetrics,
		IpfsOracle:         ipfsOracle,
		ArbiusRouter:       arbiusRouter,
	}

	ctx := context.WithValue(appContext, servicesKey{}, services)

	return ctx, nil
}

// checkContractExists verifies that a contract exists at the specified address
func checkContractExists(address common.Address, client *ethclient.Client, logger *zerolog.Logger) error {
	// Check if the address has code (only contracts have code)
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return fmt.Errorf("failed to get code at address: %w", err)
	}

	if len(code) == 0 {
		return fmt.Errorf("no contract code found at address")
	}

	return nil
}
