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
	Basetoken    *basetoken.BaseToken
	Engine       *EngineWrapper
	Voter        *voter.Voter
	ArbiusRouter ipfs.ArbiusRouterContract
	BulkTasks    *bulktasks.BulkTasks
	//Eth            *erc20.TokenERC20
	OwnerAccount   *account.Account
	Config         *config.AppConfig
	Logger         zerolog.Logger
	TaskStorage    *storage.TaskStorageDB
	AutoMineParams *SubmitTaskParams
	Paraswap       *paraswap.ParaswapManager
	OracleProvider IPriceOracle
	LeverOracle    ILeverOracle
	TaskTracker    *metrics.TaskTracker
	IpfsOracle     ipfs.OracleClient
	Validators     *Validators
}

func NewApplicationContext(rpc *client.Client, sql *sql.DB, logger zerolog.Logger, cfg *config.AppConfig, ipfsOracle ipfs.OracleClient, appContext, appQuit context.Context) (context.Context, *Services, error) {

	ownerAccount, err := account.NewAccount(cfg.Blockchain.PrivateKey, rpc, appContext, cfg.Blockchain.CacheNonce, logger)
	if err != nil {
		return nil, nil, err
	}

	// TODO: need a cleaner way to handle this nonce update on first use - maybe move to NewAccount
	ownerAccount.UpdateNonce()

	baseTokenContract, err := basetoken.NewBaseToken(cfg.BaseConfig.BaseTokenAddress, rpc.Client)
	if err != nil {
		return nil, nil, err
	}

	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpc.Client)
	if err != nil {
		return nil, nil, err
	}

	voterContract, err := voter.NewVoter(cfg.BaseConfig.VoterAddress, rpc.Client)
	if err != nil {
		return nil, nil, err
	}

	// if we are on sepolia or mainnet we can use a real contract otherwise we mock (determined by config)
	var arbiusRouter ipfs.ArbiusRouterContract
	if cfg.BaseConfig.ArbiusRouterAddress == (common.Address{}) {
		arbiusRouter = ipfs.NewMockRouterContract()
	} else {
		arbiusRouterContract, err := arbiusrouterv1.NewArbiusRouterV1(cfg.BaseConfig.ArbiusRouterAddress, rpc.Client)
		if err != nil {
			return nil, nil, err
		}
		arbiusRouter = ipfs.NewRouterContractWrapper(arbiusRouterContract, rpc)
	}

	var bulkTasksContract *bulktasks.BulkTasks

	if cfg.BaseConfig.TestnetType > 0 {
		// TODO: this is a hack to check if the contract exists on testnet
		err := checkContractExists(cfg.BaseConfig.BulkTasksAddress, rpc.Client)

		if err != nil {

			// on local testnet we can deploy the contract
			if cfg.BaseConfig.TestnetType > 0 {
				ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				// get the gas price
				gasPrice, err := ownerAccount.Client.GetGasPrice()
				if err != nil {
					return nil, nil, err
				}

				opts := ownerAccount.GetOpts(1_000_000, gasPrice, nil, nil)
				// deploy the contract
				var bulkTasksContractAddress common.Address
				var tx *types.Transaction

				logger.Info().Msg("deploying BulkTasks contract")

				bulkTasksContractAddress, tx, bulkTasksContract, err = bulktasks.DeployBulkTasks(opts, rpc.Client, cfg.BaseConfig.BaseTokenAddress, cfg.BaseConfig.EngineAddress)
				if err != nil {
					return nil, nil, err
				}

				logger.Info().Msg("waiting for BulkTasks contract deployment")

				receipt, err := bind.WaitMined(ctxTimeout, rpc.Client, tx)
				if err != nil {
					return nil, nil, err
				}

				if receipt.Status != types.ReceiptStatusSuccessful {
					return nil, nil, fmt.Errorf("BulkTasks contract deployment failed")
				}

				logger.Info().Msgf("BulkTasks contract deployed at address %s (you must update the config with this address)", bulkTasksContractAddress.Hex())

				cfg.BaseConfig.BulkTasksAddress = bulkTasksContractAddress
			} else {
				return nil, nil, fmt.Errorf("BulkTasks contract not deployed or invalid at address %s: %w",
					cfg.BaseConfig.BulkTasksAddress.Hex(), err)
			}
		}
	}

	if bulkTasksContract == nil {
		bulkTasksContract, err = bulktasks.NewBulkTasks(cfg.BaseConfig.BulkTasksAddress, rpc.Client)
		if err != nil {
			return nil, nil, err
		}
	}

	// cache the min claim time
	minClaimSolTimeBig, err := engineContract.MinClaimSolutionTime(nil)
	if err != nil {
		return nil, nil, err
	}

	// A convienience wrapper to represent ETH
	eth := erc20.NewTokenERC20(common.HexToAddress("0x0"), 18, "ETH", "ETH")

	// 120 = jitter offset in seconds for the min claim time
	minclaimTime := time.Duration(minClaimSolTimeBig.Uint64()+120) * time.Second

	logger.Info().Msgf("minimum claim time is %s", minclaimTime)

	ts := storage.NewTaskStorageDB(appContext, sql, minclaimTime, logger)

	engineWrapper := NewEngineWrapper(engineContract, voterContract, logger)

	ratelimitEth, err := engineWrapper.Engine.SolutionRateLimit(nil)
	if err != nil {
		return nil, nil, err
	}

	ratelimit := eth.ToFloat(ratelimitEth)

	var validators Validators
	for _, pk := range cfg.ValidatorConfig.PrivateKeys {
		va, err := NewValidator(cfg, engineWrapper, baseTokenContract, logger, appContext, pk, ownerAccount.Client, ratelimit)
		if err != nil {
			return nil, nil, err
		}
		validators.validators = append(validators.validators, va)
	}

	paraswapCacheTTL, err := time.ParseDuration(cfg.ParaswapCacheTTL)
	if err != nil {
		logger.Warn().Err(err).Str("value", cfg.ParaswapCacheTTL).Msg("failed to parse paraswap_cache_ttl, using default 5m")
		paraswapCacheTTL = 5 * time.Minute // Default
	}

	paraswapTimeout, err := time.ParseDuration(cfg.ParaswapTimeout)
	if err != nil {
		logger.Warn().Err(err).Str("value", cfg.ParaswapTimeout).Msg("failed to parse paraswap_timeout, using default 30s")
		paraswapTimeout = 30 * time.Second // Default
	}

	paraswapManager := paraswap.NewParaswapManager(
		ownerAccount,
		baseTokenContract,
		cfg.BaseConfig.BaseToken,
		logger,
		paraswapCacheTTL,
		paraswapTimeout,
	)

	var oracleProvider IPriceOracle
	var leverOracle ILeverOracle
	// Try onchain oracle if configured
	if cfg.PriceOracleContract != (common.Address{}) {
		oracleContract, err := NewOnChainOracle(rpc.Client, cfg.PriceOracleContract, eth, cfg.BaseConfig.BaseToken, logger)
		if err != nil {
			// If onchain oracle fails, return the error
			return nil, nil, fmt.Errorf("failed to create onchain oracle: %w", err)
		}
		// Success: override defaults
		oracleProvider = oracleContract
		if cfg.Claim.UseLever {
			leverOracle = oracleContract
		} else {
			leverOracle = NewMinClaimLeverOracle(cfg.Claim.ClaimMinReward)
		}
	} else {
		// otherwise use paraswap for general price oracle
		// NOTE: the api has rate limiting so not suitable for high volume/block update
		oracleProvider = paraswapManager
		// and use the basic lever oracle
		leverOracle = NewMinClaimLeverOracle(cfg.Claim.ClaimMinReward)
	}

	taskMetrics := metrics.NewTaskTracker(appQuit)

	input, err := json.Marshal(cfg.Strategies.Automine.Input)
	if err != nil {
		return nil, nil, err
	}

	st := &SubmitTaskParams{
		Version: uint8(cfg.Strategies.Automine.Version),
		Owner:   ownerAccount.Address,
		Model:   cfg.Strategies.Automine.ModelAsBytes,
		Fee:     cfg.Strategies.Automine.Fee,
		Input:   input,
	}

	services := &Services{
		Basetoken: baseTokenContract,
		Engine:    engineWrapper,
		Voter:     voterContract,
		BulkTasks: bulkTasksContract,
		//Eth:            eth,
		OwnerAccount:   ownerAccount,
		Config:         cfg,
		Logger:         logger,
		TaskStorage:    ts,
		AutoMineParams: st,
		Paraswap:       paraswapManager,
		OracleProvider: oracleProvider,
		LeverOracle:    leverOracle,
		TaskTracker:    taskMetrics,
		IpfsOracle:     ipfsOracle,
		ArbiusRouter:   arbiusRouter,
		Validators:     &validators,
	}

	ctx := context.WithValue(appContext, servicesKey{}, services)

	return ctx, services, nil
}

// checkContractExists verifies that a contract exists at the specified address
func checkContractExists(address common.Address, client *ethclient.Client) error {
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
