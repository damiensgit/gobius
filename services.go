package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"gobius/account"
	"gobius/arbius/basetoken"
	"gobius/arbius/bulktasks"
	"gobius/arbius/engine"
	"gobius/client"
	"gobius/config"
	"gobius/erc20"
	"gobius/metrics"
	"gobius/paraswap"
	"gobius/storage"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

type servicesKey struct{}

type Services struct {
	Basetoken          *basetoken.Basetoken
	Engine             *EngineWrapper
	BulkTasks          *bulktasks.Bulktasks
	Eth                *erc20.TokenERC20
	OwnerAccount       *account.Account
	SenderOwnerAccount *account.Account
	Clients            []*client.Client
	Config             *config.AppConfig
	Logger             *zerolog.Logger
	TaskStorage        *storage.TaskStorageDB
	AutoMineParams     *SubmitTaskParams
	//Redis          *redis.Client
	Paraswap    *paraswap.ParaswapManager
	TaskTracker *metrics.TaskTracker
}

func NewApplicationContext(rpc *client.Client, senderrpc *client.Client, clients []*client.Client, sql *sql.DB, logger *zerolog.Logger, cfg *config.AppConfig, appContext context.Context) (context.Context, error) {

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

	baseTokenContract, err := basetoken.NewBasetoken(cfg.BaseConfig.BaseTokenAddress, senderrpc.Client)
	if err != nil {
		return nil, err
	}

	engineContract, err := engine.NewEngine(cfg.BaseConfig.EngineAddress, rpc.Client)
	if err != nil {
		return nil, err
	}

	bulkTasksContract, err := bulktasks.NewBulktasks(cfg.BaseConfig.BulkTasksAddress, senderrpc.Client)
	if err != nil {
		return nil, err
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

	// TODO: fix this context - using appcontext isnt suitable as it is cancelled on quit signal
	// which means redis calls will fail even as we try to exit
	ts := storage.NewTaskStorageDB(appContext, sql, minclaimTime, logger)

	engineWrapper := NewEngineWrapper(engineContract, logger)

	paraswapManager := paraswap.NewParaswapManager(
		senderOwnerAccount,
		baseTokenContract,
		cfg.BaseConfig.BaseToken,
		logger,
	)

	taskMetrics := metrics.NewTaskTracker()

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
	}

	ctx := context.WithValue(appContext, servicesKey{}, services)

	return ctx, nil
}
