package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gobius/bindings/arbiusrouterv1"
	"gobius/bindings/engine"
	task "gobius/common" // Renamed import to avoid conflict
	"gobius/ipfs"
	"gobius/models"
	"gobius/utils"
	"math/rand"
	"sync"
	"time"

	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
)

// Core Structs
type SubmitTaskParams struct {
	Version        uint8
	Owner          common.Address
	Model          [32]byte
	Fee            *big.Int
	Input          []byte
	Incentive      *big.Int
	OriginMethod   string
	OriginContract common.Address
}

type TaskSubmitted struct {
	TaskId task.TaskId
	TxHash common.Hash
}

// Interfaces
type MiningStrategy interface {
	Start() error
	Stop()
	Name() string
}

type SolutionSubmitter interface {
	Start(ctx context.Context, enableBatchProcessing bool) error // New: Initialize and potentially start batch loops
	SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error
	SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error
	SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error
	SignalCommitmentNow(validator common.Address, taskId task.TaskId, commitment [32]byte) error // New: Immediate submission
	SubmitSolutionNow(validator common.Address, taskId task.TaskId, cid []byte) error            // New: Immediate submission
	Stop() error                                                                                 // Added Stop method
}

type TaskProducer interface {
	TaskChannel() <-chan *TaskSubmitted
	Start(ctx context.Context) error
	Stop()
	Name() string
}

// TaskInfoCache defines the method needed by SolutionEventProducer to get TxHash from TaskId.
type TaskInfoCache interface {
	GetTxHash(taskIdStr string) (common.Hash, bool)
}

// Utilities

type goroutineRunner struct {
	sync.WaitGroup
}

func (gr *goroutineRunner) Go(f func()) {
	gr.Add(1)
	go func() {
		defer gr.Done()
		f()
	}()
}

func (gr *goroutineRunner) GoWithContext(ctx context.Context, f func(ctx context.Context)) {
	gr.Add(1)
	go func() {
		defer gr.Done()
		f(ctx)
	}()
}

var ErrProducerStopped = errors.New("producer stopped")
var ErrTxDecodePermanent = errors.New("permanent transaction decoding error")

// Subscription Manager
const (
	defaultMaxSubscriptionBackoff     = 30 * time.Second
	defaultInitialSubscriptionBackoff = 1 * time.Second
)

type ConnectFunc func(ctx context.Context) (event.Subscription, error)

// Task handler function signature now includes its own context
type taskHandlerFunc func(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context)

type subscriptionManager struct {
	logger         zerolog.Logger
	parentCtx      context.Context
	ctx            context.Context
	cancel         context.CancelFunc
	connectFunc    ConnectFunc
	eventName      string
	maxBackoff     time.Duration
	initialBackoff time.Duration
	mu             sync.Mutex
	subscription   event.Subscription
	stopOnce       sync.Once
	goroutineRunner
}

func NewSubscriptionManager(parentCtx context.Context, log zerolog.Logger, eventName string, connectFunc ConnectFunc) *subscriptionManager {
	ctx, cancel := context.WithCancel(parentCtx)
	return &subscriptionManager{
		logger:          log.With().Str("component", "subscriptionmanager").Str("event", eventName).Logger(),
		parentCtx:       parentCtx,
		ctx:             ctx,
		cancel:          cancel,
		connectFunc:     connectFunc,
		eventName:       eventName,
		maxBackoff:      defaultMaxSubscriptionBackoff,
		initialBackoff:  defaultInitialSubscriptionBackoff,
		goroutineRunner: goroutineRunner{},
	}
}

func (sm *subscriptionManager) Start() {
	sm.logger.Info().Msg("starting")
	sm.Go(sm.manageLoop)
}

func (sm *subscriptionManager) Stop() {
	sm.stopOnce.Do(func() {
		sm.logger.Info().Msg("stopping")
		if sm.cancel != nil {
			sm.cancel()
		}
		sm.mu.Lock()
		if sm.subscription != nil {
			sm.subscription.Unsubscribe()
			sm.subscription = nil
		}
		sm.mu.Unlock()
		sm.Wait() // Wait for manageLoop
		sm.logger.Info().Msg("stopped")
	})
}

// manageLoop is the core goroutine managing the subscription lifecycle.
func (sm *subscriptionManager) manageLoop() {
	sm.logger.Info().Msg("management loop started")
	defer sm.logger.Info().Msg("management loop stopped")

	currentBackoff := sm.initialBackoff
	var errChan <-chan error

	if !sm.connectWithTimeout() {
		sm.logger.Warn().Msg("initial connection failed, will retry in loop")
	}

	for {
		sm.mu.Lock()
		actSub := sm.subscription
		sm.mu.Unlock()

		if actSub == nil { // Subscription is down
			sm.logger.Warn().Dur("wait", currentBackoff).Msg("subscription down, attempting reconnect")
			select {
			case <-time.After(currentBackoff):
				currentBackoff = (currentBackoff * 2) + time.Duration(rand.Intn(500))*time.Millisecond
				if currentBackoff > sm.maxBackoff {
					currentBackoff = sm.maxBackoff
				}
				if !sm.connectWithTimeout() {
					continue // Reconnect failed, loop again
				}
				// Reconnect succeeded
				currentBackoff = sm.initialBackoff
				sm.mu.Lock()
				if sm.subscription != nil {
					errChan = sm.subscription.Err()
				} else {
					errChan = nil // Should not happen if connect succeeded
				}
				sm.mu.Unlock()
			case <-sm.ctx.Done():
				sm.logger.Info().Msg("shutting down during reconnect backoff")
				return
			}
		} else { // Subscription is active
			errChan = actSub.Err()
		}

		// Wait for an error or context cancellation
		select {
		case <-sm.ctx.Done(): // Owner context cancelled
			sm.logger.Info().Msg("parent context cancelled, shutting down loop")
			// Unsubscribe is handled in Stop(), which should have cancelled ctx
			return
		case err, ok := <-errChan:
			if !ok || err != nil { // Channel closed or error received
				if !ok && err == nil { // Don't log error if channel just closed cleanly
					err = errors.New("subscription error channel closed unexpectedly")
				}
				if err != nil {
					sm.logger.Warn().Err(err).Msg("subscription error detected")
				}
				// Subscription is broken, reset it
				sm.mu.Lock()
				if sm.subscription != nil {
					sm.subscription.Unsubscribe()
					sm.subscription = nil
				}
				sm.mu.Unlock()
				// Loop will now enter the reconnect logic
			}
		}
	}
}

// connectWithTimeout attempts to establish the subscription using the connectFunc.
func (sm *subscriptionManager) connectWithTimeout() bool {
	// Use manager's internal context for the connection attempt itself
	connectCtx, cancel := context.WithTimeout(sm.ctx, 20*time.Second) // TODO: configurable timeout
	defer cancel()

	sm.logger.Info().Msg("attempting to connect subscription...")
	// The connectFunc performs the contract Watch call.
	// It should use the appropriate context (likely sm.parentCtx) for the WatchOpts.
	newSub, err := sm.connectFunc(connectCtx)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			sm.logger.Warn().Err(err).Msg("connection attempt timed out or was cancelled")
		} else {
			sm.logger.Error().Err(err).Msg("connection function failed")
		}
		return false
	}
	if newSub == nil { // Should not happen if connectFunc adheres to convention
		sm.logger.Error().Msg("connection function returned nil subscription without error")
		return false
	}

	sm.logger.Info().Msg("connection successful")
	sm.mu.Lock()
	// Check if Stop was called while we were connecting
	if sm.ctx.Err() != nil {
		sm.logger.Warn().Msg("context cancelled during successful connection, unsubscribing immediately")
		sm.mu.Unlock()
		newSub.Unsubscribe() // Clean up the newly created sub
		return false
	}
	// Assign the new subscription
	sm.subscription = newSub
	sm.mu.Unlock()
	return true
}

type baseStrategy struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	services     *Services
	submitter    SolutionSubmitter
	gpuPool      *GPUPool
	logger       zerolog.Logger
	numWorkers   int
	stopOnce     sync.Once
	strategyName string

	// Shared Resources
	txTaskCache        *lru.Cache // TaskIdStr -> common.Hash, TxHashStr -> *SubmitTaskParams
	decodingInProgress sync.Map   // map[common.Hash]chan *SubmitTaskParams

	// ABI Methods
	submitMethod              *abi.Method
	bulkSubmitMethod          *abi.Method
	submitMethodOnRouter      *abi.Method
	submitTaskWithTokenMethod *abi.Method
	submitTaskWithETHMethod   *abi.Method

	// Event Subscription Management (Owned by BaseStrategy)
	taskSubManager    *subscriptionManager                 // Owns TaskSubmitted sub
	taskSubmittedSink chan *engine.EngineTaskSubmitted     // Receives raw TaskSubmitted events
	verifySubManager  *subscriptionManager                 // Owns SolutionSubmitted sub
	verifySink        chan *engine.EngineSolutionSubmitted // Receives raw SolutionSubmitted events

	// Shared channel for baseline sampled validation tasks (populated by EventManager)
	sampledValidationChan chan *TaskSubmitted

	goroutineRunner // Manages worker goroutines ONLY
}

func (bs *baseStrategy) Go(f func()) { bs.goroutineRunner.Go(f) }

// GetTxHash implements TaskInfoCache for SolutionEventProducer.
func (bs *baseStrategy) GetTxHash(taskIdStr string) (common.Hash, bool) {
	// Look specifically for the common.Hash type cached by EventManager
	if val, found := bs.txTaskCache.Get(taskIdStr); found {
		if hashVal, ok := val.(common.Hash); ok {
			bs.logger.Trace().Str("task", taskIdStr).Str("txHash", hashVal.Hex()).Msg("cache hit for taskid->txhash")
			return hashVal, true
		}
		bs.logger.Warn().Str("task", taskIdStr).Msgf("cache contained non-hash type (%T) for taskid->txhash lookup", val)
	} else {
		bs.logger.Trace().Str("task", taskIdStr).Msg("cache miss for taskid->txhash lookup")
	}
	return common.Hash{}, false
}

func newBaseStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool, strategyName string) (*baseStrategy, error) {
	ctx, cancel := context.WithCancel(appCtx)
	numWorkers := gpuPool.NumGPUs() * services.Config.NumWorkersPerGPU
	logger := services.Logger.With().Str("strategy", strategyName).Logger()

	const cacheSize = 200_000 // TODO: Configurable
	taskCache, err := lru.New(cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LRU cache: %w", err)
	}

	// Load ABIs and methods once at initialization
	engineAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get engine ABI: %w", err)
	}
	submitMethod, ok := engineAbi.Methods["submitTask"]
	if !ok {
		return nil, errors.New("submitTask method not found in engine ABI")
	}
	bulkSubmitMethod, ok := engineAbi.Methods["bulkSubmitTask"]
	if !ok {
		return nil, errors.New("bulkSubmitTask method not found in engine ABI")
	}

	routerAbi, err := arbiusrouterv1.ArbiusRouterV1MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get router ABI: %w", err)
	}
	submitMethodOnRouter, ok := routerAbi.Methods["submitTask"]
	if !ok {
		return nil, errors.New("submitTask method not found in router ABI")
	}
	submitTaskWithTokenMethod, ok := routerAbi.Methods["submitTaskWithToken"]
	if !ok {
		return nil, errors.New("submitTaskWithToken method not found in router ABI")
	}
	submitTaskWithETHMethod, ok := routerAbi.Methods["submitTaskWithETH"]
	if !ok {
		return nil, errors.New("submitTaskWithETH method not found in router ABI")
	}
	// abi methods loaded

	bufferSize := numWorkers * 2
	if bufferSize < 10 {
		bufferSize = 10
	} // Minimum buffer

	bs := &baseStrategy{
		ctx:                       ctx,
		cancelFunc:                cancel,
		services:                  services,
		submitter:                 submitter,
		gpuPool:                   gpuPool,
		logger:                    logger,
		numWorkers:                numWorkers,
		strategyName:              strategyName,
		txTaskCache:               taskCache,
		taskSubmittedSink:         make(chan *engine.EngineTaskSubmitted, 1024), // Generous buffer to handle bursts
		verifySink:                make(chan *engine.EngineSolutionSubmitted, bufferSize),
		sampledValidationChan:     make(chan *TaskSubmitted, bufferSize),
		submitMethod:              &submitMethod,
		bulkSubmitMethod:          &bulkSubmitMethod,
		submitMethodOnRouter:      &submitMethodOnRouter,
		submitTaskWithTokenMethod: &submitTaskWithTokenMethod,
		submitTaskWithETHMethod:   &submitTaskWithETHMethod,
		goroutineRunner:           goroutineRunner{},
	}

	// Setup subscriptions
	bs.taskSubManager = NewSubscriptionManager(bs.ctx, bs.logger, "TaskSubmitted", func(cCtx context.Context) (event.Subscription, error) {
		cl := bs.services.OwnerAccount.Client.Client
		if cl == nil {
			return nil, errors.New("client nil for task sub")
		}
		co, err := engine.NewEngine(bs.services.Config.BaseConfig.EngineAddress, cl)
		if err != nil {
			return nil, fmt.Errorf("new engine failed: %w", err)
		}
		bn, err := cl.BlockNumber(cCtx)
		if err != nil {
			return nil, fmt.Errorf("block number failed: %w", err)
		}
		return co.WatchTaskSubmitted(&bind.WatchOpts{Start: &bn, Context: bs.ctx}, bs.taskSubmittedSink, nil, nil, nil)
	})

	bs.verifySubManager = NewSubscriptionManager(bs.ctx, bs.logger, "SolutionSubmitted", func(cCtx context.Context) (event.Subscription, error) {
		cl := bs.services.OwnerAccount.Client.Client
		if cl == nil {
			return nil, errors.New("client nil for solution sub")
		}
		co := bs.services.Engine.Engine
		if co == nil {
			return nil, errors.New("engine service nil for solution sub")
		}
		bn, err := cl.BlockNumber(cCtx)
		if err != nil {
			return nil, fmt.Errorf("block number failed: %w", err)
		}
		return co.WatchSolutionSubmitted(&bind.WatchOpts{Start: &bn, Context: bs.ctx}, bs.verifySink, nil, nil)
	})

	return bs, nil
}

// decodeTaskTransaction: Decodes raw tx data (no caching).
func (bs *baseStrategy) decodeTaskTransaction(tx *types.Transaction) (*SubmitTaskParams, error) {
	if tx.To() == nil {
		return nil, errors.New("transaction recipient is nil")
	}
	txTo := *tx.To()
	isEngine := txTo == bs.services.Config.BaseConfig.EngineAddress
	isRouter := txTo == bs.services.Config.BaseConfig.ArbiusRouterAddress

	if !isEngine && !isRouter {
		return nil, fmt.Errorf("transaction not sent to known contract (engine=%s, router=%s)",
			bs.services.Config.BaseConfig.EngineAddress.Hex(),
			bs.services.Config.BaseConfig.ArbiusRouterAddress.Hex())
	}
	if len(tx.Data()) < 4 {
		return nil, errors.New("transaction data too short for method signature")
	}

	methodSig := tx.Data()[:4]
	var params []interface{}
	var err error
	var selectedMethod *abi.Method
	originMethod := "unknown"
	originContract := txTo
	incentive := big.NewInt(0)

	if isEngine {
		if bytes.Equal(bs.submitMethod.ID, methodSig) {
			selectedMethod = bs.submitMethod
			originMethod = "submitTask"
		} else if bytes.Equal(bs.bulkSubmitMethod.ID, methodSig) {
			selectedMethod = bs.bulkSubmitMethod
			originMethod = "bulkSubmitTask"
		}
	} else {
		if bytes.Equal(bs.submitMethodOnRouter.ID, methodSig) {
			selectedMethod = bs.submitMethodOnRouter
			originMethod = "submitTaskOnRouter"
		} else if bytes.Equal(bs.submitTaskWithTokenMethod.ID, methodSig) {
			selectedMethod = bs.submitTaskWithTokenMethod
			originMethod = "submitTaskWithToken"
		} else if bytes.Equal(bs.submitTaskWithETHMethod.ID, methodSig) {
			selectedMethod = bs.submitTaskWithETHMethod
			originMethod = "submitTaskWithETH"
		}
	}

	if selectedMethod == nil {
		return nil, fmt.Errorf("transaction to %s has unknown method signature: %s", txTo.Hex(), hex.EncodeToString(methodSig))
	}

	params, err = selectedMethod.Inputs.Unpack(tx.Data()[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack %s: %w", originMethod, err)
	}

	// Extract common parameters
	if len(params) < 5 {
		return nil, fmt.Errorf("unpacked fewer parameters (%d) than expected (min 5) for %s", len(params), originMethod)
	}
	version, okV := params[0].(uint8)
	owner, okO := params[1].(common.Address)
	model, okM := params[2].([32]byte)
	fee, okF := params[3].(*big.Int)
	input, okI := params[4].([]byte)
	if !okV || !okO || !okM || !okF || !okI {
		return nil, fmt.Errorf("type assertion failed for common parameters in %s", originMethod)
	}

	// Extract incentive (Index 5 for Router Methods)
	if isRouter && len(params) > 5 {
		if inc, okInc := params[5].(*big.Int); okInc {
			incentive = inc
		} else {
			bs.logger.Warn().Str("method", originMethod).Msg("type assertion failed for incentive parameter")
		}
	}

	submitParams := &SubmitTaskParams{
		Version:        version,
		Owner:          owner,
		Model:          model,
		Fee:            fee,
		Input:          input,
		Incentive:      incentive,
		OriginMethod:   originMethod,
		OriginContract: originContract,
	}
	return submitParams, nil
}

// decodeTransactionWithCache fetches and decodes, using cache and decode-in-progress map.
func (bs *baseStrategy) decodeTransactionWithCache(txHash common.Hash) (*SubmitTaskParams, error) {
	txHashStr := txHash.String()
	decodeLog := bs.logger.With().Str("tx_hash", txHashStr).Str("routine", "decodeTransactionWithCache").Logger()

	// 1. Check cache for final result (*SubmitTaskParams)
	if cached, found := bs.txTaskCache.Get(txHashStr); found {
		if params, ok := cached.(*SubmitTaskParams); ok {
			decodeLog.Debug().Msg("cache hit (full params)")
			return params, nil
		}
		// Ignore other types (like initial common.Hash mapping from TaskId)
		decodeLog.Debug().Msgf("cache contained non-params type (%T), proceeding to decode", cached)
	} else {
		decodeLog.Debug().Msg("cache miss")
	}

	// 2. Check if decode is already in progress by another goroutine
	waitChan := make(chan *SubmitTaskParams, 1)
	val, loaded := bs.decodingInProgress.LoadOrStore(txHash, waitChan)

	if loaded { // Another goroutine is already decoding this txHash
		close(waitChan) // Close the channel we created but didn't store/use
		existingChan, ok := val.(chan *SubmitTaskParams)
		if !ok {
			decodeLog.Error().Msgf("invalid type in decodingInProgress map: %T", val)
			return nil, fmt.Errorf("internal error: invalid type in decoding map for %s", txHashStr)
		}

		decodeLog.Info().Msg("decode already in progress, waiting...")
		waitCtx, cancel := context.WithTimeout(bs.ctx, 60*time.Second) // TODO: Configurable wait timeout
		defer cancel()
		select {
		case params, chanOk := <-existingChan:
			if !chanOk {
				decodeLog.Error().Msg("decode channel closed unexpectedly while waiting")
				return nil, fmt.Errorf("decode channel closed unexpectedly for %s", txHashStr)
			}
			if params != nil {
				decodeLog.Info().Msg("decode finished, got params from waiter")
				return params, nil // Success
			} else {
				decodeLog.Warn().Msg("decode failed (waited for result)")
				return nil, fmt.Errorf("decode failed for tx %s (waited)", txHashStr)
			}
		case <-waitCtx.Done():
			decodeLog.Warn().Msg("timeout waiting for decode result")
			return nil, fmt.Errorf("timeout waiting for decode tx %s", txHashStr)
		}
	}

	// 3. We are the first goroutine for this txHash, initiate the decode
	var params *SubmitTaskParams
	decodeLog.Info().Msg("initiating decode...")
	defer func() {
		decodeLog.Debug().Msg("cleaning up decodingInProgress map entry")
		bs.decodingInProgress.Delete(txHash)
		waitChan <- params // Send params (if successful) or nil (if failed)
		close(waitChan)    // Signal completion/failure
	}()

	// Fetch transaction (RPC call)
	fetchCtx, cancelFetch := context.WithTimeout(bs.ctx, 30*time.Second) // TODO: Configurable fetch timeout
	tx, _, err := bs.services.OwnerAccount.Client.Client.TransactionByHash(fetchCtx, txHash)
	cancelFetch()

	if err != nil {
		// Handle permanent vs transient fetch errors
		if errors.Is(err, ethereum.NotFound) {
			decodeLog.Warn().Err(err).Msg("transaction not found")
			return nil, err
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			decodeLog.Warn().Err(err).Msg("context error fetching transaction")
			return nil, err
		}
		decodeLog.Error().Err(err).Msg("transient error fetching transaction")
		return nil, err // Assume others are transient
	}
	if tx == nil {
		decodeLog.Error().Msg("transaction not found (nil tx returned)")
		return nil, ethereum.NotFound
	}

	// Decode the transaction data
	params, err = bs.decodeTaskTransaction(tx)
	if err != nil {
		params = nil // Ensure params is nil on error
		decodeLog.Error().Err(err).Msg("decode failed")
		// Don't wrap ErrTxDecodePermanent here, let caller check if needed
		return nil, err // Return the original decode error
	}

	// Success: Cache result (keyed by TxHash) and return
	bs.txTaskCache.Add(txHashStr, params)
	decodeLog.Info().Str("origin", params.OriginMethod).Msg("decode successful, params cached")
	return params, nil
}

// solveTask: Executes model inference and optionally submits.
func (bs *baseStrategy) solveTask(ctx context.Context, taskId task.TaskId, params *SubmitTaskParams, gpu *task.GPU, validateOnly bool, submitImmediately bool /* New Param */) ([]byte, error) {
	taskIdStr := taskId.String()
	solveLog := bs.logger.With().Str("task", taskIdStr).Str("origin", params.OriginMethod).Logger()

	// Check model match
	modelIdHex := "0x" + common.Bytes2Hex(params.Model[:])
	configuredModel := bs.services.Config.Strategies.Model
	if modelIdHex != configuredModel {
		solveLog.Warn().Str("task_model", modelIdHex).Str("configured_model", configuredModel).Msg("skipping task, model mismatch")
		return nil, nil // Skipped, not an error
	}

	model := models.ModelRegistry.GetModel(modelIdHex)
	if model == nil {
		solveLog.Error().Str("model", modelIdHex).Msg("model specified in task not found or enabled")
		return nil, fmt.Errorf("model %s not found or enabled", modelIdHex)
	}

	// Unmarshal and Hydrate Input
	var inputMap map[string]interface{}
	if err := json.Unmarshal(params.Input, &inputMap); err != nil {
		solveLog.Error().Err(err).Str("input_raw", string(params.Input)).Msg("could not unmarshal task input json")
		return nil, err
	}

	hydratedInput, err := model.HydrateInput(inputMap, taskId.TaskId2Seed())
	if err != nil {
		solveLog.Error().Err(err).Msg("could not hydrate input")
		return nil, err
	}
	// solveLog.Debug().Interface("hydrated_input", hydratedInput).Msg("hydrated input") // Verbose

	// Inference
	var cid []byte
	var inferenceElapsed time.Duration // Keep track of time for logging
	if bs.services.Config.EvilMode {
		cid, _ = hex.DecodeString("12206666666666666666666666666666666666666666666666666666666666666666") // Example evil CID
		duration := time.Duration(bs.services.Config.EvilModeMinTime+rand.Intn(bs.services.Config.EvilModeRandInt)) * time.Millisecond
		solveLog.Warn().Dur("sleep", duration).Msg("evil mode: simulating inference")
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			solveLog.Warn().Msg("context cancelled during evil mode sleep")
			return nil, ctx.Err()
		}
		inferenceElapsed = duration // Approximate
	} else {
		solveLog.Debug().Msg("sending task to gpu for inference")
		inferenceStart := time.Now()
		if gpu.Mock {
			data, errM := gpu.GetMockCid(taskIdStr, hydratedInput)
			if errM == nil {
				cid, err = ipfs.GetIPFSHashFast(data) // Use fast path if mock provides raw data
			} else {
				err = errM // Propagate mock error
			}
		} else {
			cid, err = model.GetCID(ctx, gpu, taskIdStr, hydratedInput)
		}
		inferenceElapsed = time.Since(inferenceStart)

		if err != nil {
			// Handle context cancellation specifically
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				solveLog.Info().Err(err).Msg("inference cancelled or timed out")
				// Do not mark GPU error for context issues
			} else {
				solveLog.Error().Err(err).Msg("inference failed")
				gpu.IncrementErrorCount() // Mark GPU error for actual failures
			}
			return nil, err // Propagate error
		}
		solveLog.Info().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", inferenceElapsed.String()).Msg("inference successful")
	}

	if validateOnly {
		return cid, nil // Return CID without submitting
	}

	validatorAddr := bs.services.Validators.GetNextValidatorAddress() // Keep internal selection
	solveLog = solveLog.With().Str("validator", validatorAddr.Hex()).Logger()

	// Generate Commitment
	commitment, err := utils.GenerateCommitment(validatorAddr, taskId, cid)
	if err != nil {
		solveLog.Error().Err(err).Msg("error generating commitment hash")
		return nil, err // Cannot proceed
	}

	// Signal Commitment (Conditional)
	if submitImmediately {
		solveLog.Info().Msg("signalling commitment immediately...")
		err = bs.submitter.SignalCommitmentNow(validatorAddr, taskId, commitment)
	} else {
		solveLog.Info().Msg("queueing commitment...")
		err = bs.submitter.SignalCommitment(validatorAddr, taskId, commitment) // Queues via storage
	}
	if err != nil {
		solveLog.Error().Err(err).Msg("failed to signal/queue commitment")
		// Consider if this is fatal or if we should still try submitting
		return nil, err
	}
	solveLog.Info().Msg("commitment signalled/queued")

	// Optimistically submit IPFS CID (best effort)
	// Check incentive threshold before submitting CID for potential claim
	minAiusThreshold := bs.services.Config.IPFS.MinAiusIncentiveThreshold // Get threshold
	incentiveExists := params.Incentive != nil && params.Incentive.Cmp(big.NewInt(0)) > 0

	if incentiveExists {
		submitCidForIncentive := true // Assume we submit unless threshold check fails
		if minAiusThreshold > 0 {     // Only check if threshold is enabled
			incentiveAmountFloat := bs.services.Config.BaseConfig.BaseToken.ToFloat(params.Incentive)
			if incentiveAmountFloat < minAiusThreshold {
				solveLog.Info().
					Float64("incentive_aius", incentiveAmountFloat).
					Float64("threshold_aius", minAiusThreshold).
					Msg("incentive below threshold, skipping IPFS CID submission")
				submitCidForIncentive = false
			}
		}

		if submitCidForIncentive {
			err = bs.submitter.SubmitIpfsCid(validatorAddr, taskId, cid)
			if err != nil {
				solveLog.Warn().Err(err).Msg("ipfs cid submission failed (non-fatal)")
				// Reset err to nil as this is non-fatal for the solve process
				err = nil
			}
		}
	}

	// Submit Solution (Conditional)
	if submitImmediately {
		solveLog.Info().Msg("submitting solution immediately...")
		err = bs.submitter.SubmitSolutionNow(validatorAddr, taskId, cid)
	} else {
		solveLog.Info().Msg("queueing solution...")
		err = bs.submitter.SubmitSolution(validatorAddr, taskId, cid) // Queues via storage
	}
	if err != nil {
		solveLog.Error().Err(err).Msg("solution submission/queuing failed")
		return nil, err // This is critical
	}
	solveLog.Info().Msg("solution submitted/queued successfully")

	return cid, nil // Full success
}

// handleSingleTask: Common logic for processing high-priority single tasks.
func (bs *baseStrategy) handleSingleTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()
	workerLogger.Info().Msg("processing single task event")
	gpu.SetStatus("Mining Single Task")

	params, err := bs.decodeTransactionWithCache(ts.TxHash)
	if err != nil || params == nil {
		workerLogger.Error().Err(err).Msg("failed to get params for single task")
		gpu.SetStatus("Error")
		return
	}

	// Safety check (should be pre-filtered)
	if params.OriginMethod == "bulkSubmitTask" {
		workerLogger.Error().Msg("bulk task received on single task channel unexpectedly, ignoring")
		gpu.SetStatus("Idle")
		return
	}

	// Solve and submit the task IMMEDIATELY
	solveStart := time.Now()
	// Call solveTask with submitImmediately: true
	cid, err := bs.solveTask(taskCtx, ts.TaskId, params, gpu, false, true /* submitImmediately */)
	solveElapsed := time.Since(solveStart)

	taskFailed := err != nil
	taskSkipped := cid == nil && err == nil

	// Set GPU status based on outcome
	if taskFailed && gpu.IsEnabled() && !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		gpu.SetStatus("Error")
	} else {
		gpu.SetStatus("Idle") // Set Idle for success, skip, cancellation, or disabled GPU
	}

	// Log results
	if taskFailed {
		workerLogger.Error().Err(err).Msg("single task solve failed")
		// Event tasks are generally not requeued
	} else if taskSkipped {
		workerLogger.Info().Msg("single task skipped")
	} else {
		workerLogger.Info().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", solveElapsed.String()).Msg("single task solved successfully")
		bs.gpuPool.AddSolveTime(solveElapsed)
	}
}

// handleValidationTask: Common logic for validating sampled solutions.
func (bs *baseStrategy) handleValidationTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Str("routine", "handleValidationTask").Logger()
	workerLogger.Info().Msg("validating sampled task")
	gpu.SetStatus("Validating") // Set status

	params, err := bs.decodeTransactionWithCache(ts.TxHash)
	if err != nil || params == nil {
		workerLogger.Error().Err(err).Msg("failed to get params for validation task")
		gpu.SetStatus("Error") // Mark error if decode fails
		return
	}

	// 2. Solve Locally (validateOnly=true)
	solveStart := time.Now()
	// submitImmediately flag is irrelevant when validateOnly is true, pass false
	ourCidBytes, err := bs.solveTask(taskCtx, ts.TaskId, params, gpu, true, false /* submitImmediately */)
	solveElapsed := time.Since(solveStart)

	taskFailed := err != nil
	taskSkipped := ourCidBytes == nil && err == nil

	// Handle solve outcome
	if taskFailed {
		workerLogger.Error().Err(err).Msg("validation: local solve failed")
		if gpu.IsEnabled() && !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			gpu.SetStatus("Error")
		} else {
			gpu.SetStatus("Idle")
		}
		return // Cannot proceed if solve failed
	}
	if taskSkipped {
		workerLogger.Info().Msg("validation: local solve skipped (e.g., model mismatch)")
		gpu.SetStatus("Idle")
		return // Cannot proceed if solve skipped
	}
	ourCid := common.Bytes2Hex(ourCidBytes)
	workerLogger.Info().Str("our_cid", ourCid).Str("elapsed", solveElapsed.String()).Msg("task solved locally for validation")

	// Compare with on-chain data
	res, err := bs.services.Engine.GetSolution(ts.TaskId)
	if err != nil {
		workerLogger.Warn().Err(err).Msg("failed get on-chain solution")
		gpu.SetStatus("Idle")
		return
	}
	if res.Blocktime == 0 {
		workerLogger.Warn().Msg("no on-chain solution found")
		gpu.SetStatus("Idle")
		return
	}
	solversCid := common.Bytes2Hex(res.Cid[:])
	solverAddr := res.Validator
	workerLogger.Info().Str("solver_cid", solversCid).Str("solver", solverAddr.Hex()).Msg("comparing CIDs")

	contestData, err := bs.services.Engine.GetContestation(ts.TaskId)
	if err != nil {
		workerLogger.Error().Err(err).Msg("failed get contestation status")
		gpu.SetStatus("Idle")
		return
	}
	contestationExists := contestData.Validator != (common.Address{})

	// Decide if action is needed
	actionNeeded := false
	isMatch := ourCid == solversCid
	if !isMatch {
		actionNeeded = true // Mismatch always needs action
	} else if contestationExists {
		actionNeeded = true // Match + Contest needs NAY vote
	} else {
		workerLogger.Info().Msg("match, no contest, no action needed.")
	}

	if !actionNeeded {
		gpu.SetStatus("Idle")
		return
	}

	// Select validator if action is needed
	selectedValidator, err := bs.selectEligibleValidator(taskCtx, ts.TaskId)
	if err != nil {
		workerLogger.Warn().Err(err).Msg("no eligible validator found, cannot take action")
		gpu.SetStatus("Idle")
		return
	}
	workerLogger.Info().Str("validator", selectedValidator.Account.Address.Hex()).Msg("selected validator for action")

	// Perform Action (Vote/Contest)
	if !isMatch { // Mismatch
		workerLogger.Warn().Msgf("=== MISMATCH %s ===", task.TaskId(ts.TaskId).String())
		workerLogger.Warn().Msgf(" Our:%s", ourCid)
		workerLogger.Warn().Msgf(" Solv:%s (%s)", solversCid, solverAddr.Hex())
		workerLogger.Warn().Msgf("=======================")
		if contestationExists {
			workerLogger.Warn().Msg("voting YEA on existing contestation...")
			err = selectedValidator.VoteOnContestation(ts.TaskId, true)
		} else {
			workerLogger.Warn().Msg("submitting NEW contestation...")
			err = selectedValidator.SubmitContestation(ts.TaskId)
		}
	} else { // Match (and contestation must exist if actionNeeded is true)
		workerLogger.Info().Msg("voting NAY on existing contestation...")
		err = selectedValidator.VoteOnContestation(ts.TaskId, false)
	}

	if err != nil {
		workerLogger.Error().Err(err).Msg("validation action failed")
	} else {
		workerLogger.Info().Msg("validation action submitted successfully")
	}
	gpu.SetStatus("Idle") // Reset status after action attempt
}

// selectEligibleValidator: Helper to find the best validator based on stake.
func (bs *baseStrategy) selectEligibleValidator(ctx context.Context, taskId task.TaskId) (*Validator, error) {
	var selectedValidator *Validator
	var bestStake = big.NewInt(0)
	eligibleFound := false

	if bs.services.Validators == nil || len(bs.services.Validators.validators) == 0 {
		return nil, errors.New("no validators configured")
	}

	valLog := bs.logger.With().Str("task", taskId.String()).Logger()
	valLog.Debug().Msg("selecting eligible validator...")

	for _, validator := range bs.services.Validators.validators {
		validatorAddr := validator.ValidatorAddress()
		vLog := valLog.With().Str("validator", validatorAddr.Hex()).Logger()

		// Check eligibility first
		isEligible, _, err := validator.IsEligibleToVote(ctx, taskId)
		if err != nil {
			vLog.Warn().Err(err).Msg("eligibility check failed")
			continue // Try next validator
		}
		if !isEligible {
			// vLog.Trace().Msg("not eligible") // Verbose
			continue
		}

		// Eligible, check stake
		vLog.Trace().Msg("eligible, checking stake...")
		staked, err := bs.services.Engine.GetValidatorStaked(validatorAddr)
		if err != nil {
			vLog.Warn().Err(err).Msg("stake check failed")
			continue // Try next validator
		}

		// Compare stake
		if staked.Cmp(bestStake) > 0 {
			// bal := bs.services.Config.BaseConfig.BaseToken.ToFloat(staked) // Optional logging
			// vLog.Debug().Float64("stake", bal).Msg("found potentially better validator")
			bestStake = staked
			selectedValidator = validator
			eligibleFound = true
		}
	}

	if !eligibleFound {
		valLog.Warn().Msg("no eligible validators found")
		return nil, errors.New("no eligible validator found")
	}
	valLog.Info().Str("selectedValidator", selectedValidator.Account.Address.Hex()).Str("stake", bestStake.String()).Msg("selected best eligible validator")
	return selectedValidator, nil
}

// gpuWorker: Core worker loop. Selects tasks from channels based on priority.
func (bs *baseStrategy) gpuWorker(workerId int, gpu *task.GPU,
	highPrioChan <-chan *TaskSubmitted, // Can be nil
	strategyChan <-chan *TaskSubmitted, // Can be nil
	strategyHandler taskHandlerFunc, // Can be nil
	validationChan <-chan *TaskSubmitted, // Can be nil
	validationHandler taskHandlerFunc, // Can be nil
) {
	workerLogger := bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Logger()
	workerLogger.Info().Msg("worker started")
	ticker := time.NewTicker(1 * time.Minute) // For periodic checks
	defer ticker.Stop()

	for {
		// Context & GPU status
		select {
		case <-bs.ctx.Done():
			workerLogger.Info().Msg("shutting down worker (context cancelled)")
			return
		default: // Don't block here
		}
		if !gpu.IsEnabled() {
			workerLogger.Debug().Msg("gpu disabled, pausing...")
			select {
			case <-time.After(15 * time.Second): // Wait before re-checking
				continue
			case <-bs.ctx.Done():
				return // Exit if stopped while paused
			}
		}

		// Determine task context
		var taskCtx context.Context
		if bs.services.Config.Miner.WaitForTasksOnShutdown {
			taskCtx = context.Background()
		} else {
			taskCtx = bs.ctx
		}

		// Prioritized task selection (non-blocking first try)
		select {
		case ts, ok := <-highPrioChan: // Highest priority
			if !ok {
				highPrioChan = nil
			} else if ts != nil {
				workerLogger.Info().Str("source", "high-prio").Msg("handling task")
				bs.handleSingleTask(workerId, gpu, ts, taskCtx)
				continue
			}
		case ts, ok := <-validationChan: // Second priority
			if !ok {
				validationChan = nil
			} else if ts != nil {
				if validationHandler != nil {
					workerLogger.Info().Str("source", "validation").Msg("handling task")
					validationHandler(workerId, gpu, ts, taskCtx)
				} else {
					workerLogger.Error().Msg("no validation handler")
				}
				continue
			}
		case ts, ok := <-strategyChan: // Third priority
			if !ok {
				strategyChan = nil
			} else if ts != nil {
				if strategyHandler != nil {
					workerLogger.Info().Str("source", "strategy").Msg("handling task")
					strategyHandler(workerId, gpu, ts, taskCtx)
				} else {
					workerLogger.Error().Msg("no strategy handler")
				}
				continue
			}
		case <-ticker.C: // Allow periodic checks if no task waiting
			// workerLogger.Trace().Msg("periodic check (no task)")
			continue
		default:
			// Blocking wait
			select {
			case <-bs.ctx.Done():
				continue // Re-check context at loop top
			case ts, ok := <-highPrioChan:
				if !ok {
					highPrioChan = nil
				} else if ts != nil {
					workerLogger.Info().Str("source", "high-prio-block").Msg("handling task")
					bs.handleSingleTask(workerId, gpu, ts, taskCtx)
					continue
				}
			case ts, ok := <-validationChan:
				if !ok {
					validationChan = nil
				} else if ts != nil {
					if validationHandler != nil {
						workerLogger.Info().Str("source", "validation-block").Msg("handling task")
						validationHandler(workerId, gpu, ts, taskCtx)
					} else {
						workerLogger.Error().Msg("no validation handler (blocking)")
					}
					continue
				}
			case ts, ok := <-strategyChan:
				if !ok {
					strategyChan = nil
				} else if ts != nil {
					if strategyHandler != nil {
						workerLogger.Info().Str("source", "strategy-block").Msg("handling task")
						strategyHandler(workerId, gpu, ts, taskCtx)
					} else {
						workerLogger.Error().Msg("no strategy handler (blocking)")
					}
					continue
				}
			case <-ticker.C:
				// workerLogger.Trace().Msg("periodic check interrupted blocking wait")
				continue // Loop back to check status/context
			}
		}

		// Check if all sources are permanently closed
		if highPrioChan == nil && strategyChan == nil && validationChan == nil {
			workerLogger.Info().Msg("all task sources finished, worker stopping.")
			return
		}
	}
}

// start launches workers, passing the necessary channels and handlers.
func (bs *baseStrategy) start(
	strategyProducer TaskProducer, // Producer for regular strategy tasks (e.g., StorageProducer)
	strategyHandler taskHandlerFunc, // Handler for tasks from strategyProducer
	// No validationProducer argument needed anymore
	validationHandler taskHandlerFunc, // Handler for baseline validation tasks
	// Channels passed directly:
	highPrioChan <-chan *TaskSubmitted, // Channel for single tasks (from EventManager)
	sampledValidationChan <-chan *TaskSubmitted, // Baseline validation tasks (from baseStrategy)
) error {
	bs.logger.Info().Int("workers", bs.numWorkers).Msgf("starting workers for %s", bs.strategyName)

	// Start stateful producers (e.g., StorageProducer) BEFORE workers
	if strategyProducer != nil {
		if err := strategyProducer.Start(bs.ctx); err != nil {
			bs.logger.Error().Err(err).Str("producer", strategyProducer.Name()).Msg("failed to start strategy producer")
			return err
		}
		bs.logger.Info().Str("producer", strategyProducer.Name()).Msg("strategy producer started")
	}

	// Get strategy producer channel (if producer exists)
	var strategyTaskChan <-chan *TaskSubmitted
	if strategyProducer != nil {
		strategyTaskChan = strategyProducer.TaskChannel()
	}

	// Launch workers
	gpus := bs.gpuPool.GetGPUs()
	if bs.numWorkers > 0 && len(gpus) == 0 {
		bs.logger.Warn().Msg("workers requested, but no GPUs available/configured")
	}
	for i := 0; i < bs.numWorkers; i++ {
		workerId := i
		gpuIndex := workerId / bs.services.Config.NumWorkersPerGPU
		if gpuIndex >= len(gpus) {
			bs.logger.Error().Int("workerId", i).Msg("gpu index out of bounds")
			continue
		}
		gpu := gpus[gpuIndex]

		// Launch worker goroutine managed by baseStrategy's WaitGroup
		bs.Go(func() {
			bs.gpuWorker(
				workerId, gpu,
				highPrioChan,          // Direct channel from EventManager (can be nil)
				strategyTaskChan,      // Direct channel from Strategy Producer (can be nil)
				strategyHandler,       // Can be nil
				sampledValidationChan, // Direct channel from BaseStrategy (can be nil)
				validationHandler,     // Can be nil
			)
		})
	}
	bs.logger.Info().Msg("workers launched")
	return nil
}

// Stop cancels context, stops owned subscriptions, waits for workers/loops.
func (bs *baseStrategy) Stop() {
	bs.stopOnce.Do(func() {
		bs.logger.Info().Msgf("stopping base strategy components for %s", bs.strategyName)
		if bs.cancelFunc != nil {
			bs.cancelFunc() // Signal workers/loops using this context
		}

		// Stop subscription managers FIRST (they wait internally)
		if bs.taskSubManager != nil {
			bs.taskSubManager.Stop()
		}
		if bs.verifySubManager != nil {
			bs.verifySubManager.Stop()
		}

		// Wait for workers and sub manager loops
		bs.Wait()

		// Close shared validation channel *after* event manager and workers stopped
		// It's safe to close a nil channel.
		if bs.sampledValidationChan != nil {
			close(bs.sampledValidationChan)
		}

		bs.logger.Info().Msgf("base strategy components stopped for %s", bs.strategyName)
	})
}

// EventManager: Central dispatcher for raw blockchain events. Routes tasks to correct places.
type EventManager struct {
	bs         *baseStrategy // Access to sinks, cache, decode methods
	logger     zerolog.Logger
	ctx        context.Context // Derived from baseStrategy context
	cancelFunc context.CancelFunc
	stopOnce   sync.Once

	// Output channel ONLY for single tasks (fed to ListenerProducer)
	singleTasksOut chan *TaskSubmitted

	sampleRate      uint64 // Baseline verification sample rate
	goroutineRunner        // Manages internal routing loops
}

func NewEventManager(bs *baseStrategy) *EventManager {
	ctx, cancel := context.WithCancel(bs.ctx)
	logger := bs.logger.With().Str("component", "eventmanager").Logger()
	sampleRate := uint64(bs.services.Config.VerificationSampleRate)
	if sampleRate == 0 {
		logger.Warn().Msg("verification sample rate 0, disabling baseline sampling")
		sampleRate = ^uint64(0)
	}
	bufferSize := bs.numWorkers * 2
	if bufferSize < 10 {
		bufferSize = 10
	}

	return &EventManager{
		bs:             bs,
		logger:         logger,
		ctx:            ctx,
		cancelFunc:     cancel,
		singleTasksOut: make(chan *TaskSubmitted, bufferSize), // Output for single tasks
		// sampledValidationTasksOut removed - routes directly to bs.sampledValidationChan
		sampleRate:      sampleRate,
		goroutineRunner: goroutineRunner{}, // Own runner for internal loops
	}
}

// Start launches event processing loops and base subscriptions.
func (em *EventManager) Start() {
	em.logger.Info().Msg("starting event manager...")
	// Start base subscriptions (they feed the sinks we read from)
	if em.bs.taskSubManager != nil {
		em.bs.taskSubManager.Start()
	} else {
		em.logger.Error().Msg("task subscription manager is nil!")
	}
	if em.bs.verifySubManager != nil {
		em.bs.verifySubManager.Start()
	} else {
		em.logger.Error().Msg("verify subscription manager is nil!")
	}
	// Start event processing/routing loops
	em.Go(em.routeTaskSubmittedEvents)
	em.Go(em.routeSolutionSubmittedEvents)
	em.logger.Info().Msg("event manager started")
}

// Stop signals internal loops, waits, and closes output channels.
func (em *EventManager) Stop() {
	em.stopOnce.Do(func() {
		em.logger.Info().Msg("stopping event manager...")
		if em.cancelFunc != nil {
			em.cancelFunc()
		} // Signal internal loops
		em.Wait()                // Wait for routing loops
		close(em.singleTasksOut) // Close owned output channel
		// Don't close bs.sampledValidationChan here
		em.logger.Info().Msg("event manager stopped")
	})
}

// routeTaskSubmittedEvents processes raw task events: caches ID->Hash, decodes, routes single tasks.
func (em *EventManager) routeTaskSubmittedEvents() {
	loopLog := em.logger.With().Str("loop", "routeTaskSubmitted").Logger()
	loopLog.Info().Msg("started task routing loop")
	defer loopLog.Info().Msg("stopped task routing loop")

	for {
		select {
		case <-em.ctx.Done():
			return // Exit on stop signal
		case event, ok := <-em.bs.taskSubmittedSink:
			if !ok {
				loopLog.Info().Msg("task sink closed")
				return
			}
			if event == nil {
				continue
			}

			taskId, txHash := event.Id, event.Raw.TxHash
			taskIdStr := task.TaskId(taskId).String()
			eventLog := loopLog.With().Str("task", taskIdStr).Str("txHash", txHash.String()).Logger()

			// 1. Cache TaskId -> TxHash mapping immediately
			em.bs.txTaskCache.Add(taskIdStr, txHash)
			eventLog.Trace().Msg("cached taskid->txhash mapping")
			// TODO: Persist mapping to DB here

			// 2. Decode synchronously (decodeTransactionWithCache handles internal locking/cache)
			eventLog.Debug().Msg("requesting decode...")
			params, err := em.bs.decodeTransactionWithCache(txHash) // Blocking call

			if err != nil {
				eventLog.Warn().Err(err).Msg("decode failed, cannot route task")
				continue
			} // Skip routing
			if params == nil {
				eventLog.Error().Msg("decode returned nil params without error")
				continue
			}
			// eventLog.Info().Str("origin", params.OriginMethod).Msg("decode successful")
			// Change log level based on origin method to reduce spam
			if params.OriginMethod == "bulkSubmitTask" {
				eventLog.Debug().Str("origin", params.OriginMethod).Msg("decode successful (bulk)")
			} else {
				eventLog.Info().Str("origin", params.OriginMethod).Msg("decode successful")
			}

			// 3. Route *single* tasks to the output channel
			if params.OriginMethod != "bulkSubmitTask" {
				ts := &TaskSubmitted{TaskId: taskId, TxHash: txHash}
				select {
				case em.singleTasksOut <- ts:
					eventLog.Info().Msg("routed single task")
				case <-em.ctx.Done():
					eventLog.Warn().Msg("context cancelled while routing single task")
					return // Exit loop
				default:
					eventLog.Warn().Int("q_len", len(em.singleTasksOut)).Msg("single task output channel full, discarding")
				}
			} else {
				eventLog.Debug().Msg("ignoring bulk task for routing")
			}
		}
	}
}

// routeSolutionSubmittedEvents processes raw solution events: samples, looks up hash, routes baseline validation tasks.
func (em *EventManager) routeSolutionSubmittedEvents() {
	loopLog := em.logger.With().Str("loop", "routeSolutionSubmitted").Logger()
	loopLog.Info().Uint64("rate", em.sampleRate).Msg("started solution sampling loop")
	defer loopLog.Info().Msg("stopped solution sampling loop")

	if em.sampleRate == 0 || em.sampleRate == ^uint64(0) {
		loopLog.Warn().Msg("baseline sampling disabled, loop exiting.")
		return
	}

	ownValidatorAddresses := make(map[common.Address]struct{})
	if em.bs.services.Validators != nil {
		for _, v := range em.bs.services.Validators.validators {
			ownValidatorAddresses[v.ValidatorAddress()] = struct{}{}
		}
	}

	for {
		select {
		case <-em.ctx.Done():
			return // Exit on stop signal
		case event, ok := <-em.bs.verifySink:
			if !ok {
				loopLog.Info().Msg("verify sink closed")
				return
			}
			if event == nil {
				continue
			}

			taskId, solverAddr := event.Task, event.Addr
			taskIdStr := task.TaskId(taskId).String()
			eventLog := loopLog.With().Str("task", taskIdStr).Str("solver", solverAddr.Hex()).Logger()

			if _, isOwn := ownValidatorAddresses[solverAddr]; isOwn {
				eventLog.Trace().Msg("skipping own solution")
				continue
			} // Skip own

			// Baseline sampling logic
			hashData := append(taskId[:], event.Raw.BlockHash.Bytes()...)
			hash := sha256.Sum256(hashData)
			sampleValue := binary.BigEndian.Uint64(hash[:8])
			if (sampleValue % em.sampleRate) != 0 {
				eventLog.Trace().Uint64("val", sampleValue).Msg("solution not sampled (baseline)")
				continue
			} // Not sampled
			eventLog.Info().Uint64("val", sampleValue).Str("solution_tx", event.Raw.TxHash.Hex()).Msg("solution SELECTED for baseline sampling")

			// Get TxHash from cache (using interface method provided by bs)
			txHash, found := em.bs.GetTxHash(taskIdStr)
			if !found {
				eventLog.Warn().Msg("txhash cache miss for sampled taskid, skipping baseline validation")
				continue
			}
			eventLog.Debug().Str("txHash", txHash.String()).Msg("found txhash for sampled task")

			// Route directly to BASE strategy's shared validation channel
			ts := &TaskSubmitted{TaskId: taskId, TxHash: txHash}
			select {
			case em.bs.sampledValidationChan <- ts:
				eventLog.Info().Msg("routed sampled task to baseline validation channel")
			case <-em.ctx.Done():
				eventLog.Warn().Msg("context cancelled while routing sampled task")
				return
			default:
				eventLog.Warn().Int("q_len", len(em.bs.sampledValidationChan)).Msg("baseline validation channel full, discarding task")
			}
		}
	}
}

// Producer Implementations

// StorageProducer: Reads from TaskStorage (Stateful)
type StorageProducer struct {
	services *Services
	logger   zerolog.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	stopOnce sync.Once
	taskChan chan *TaskSubmitted
	poolSize int
	goroutineRunner
}

func NewStorageProducer(strategyCtx context.Context, services *Services, poolSize int) *StorageProducer {
	ctx, cancel := context.WithCancel(strategyCtx)
	if poolSize <= 0 {
		poolSize = 10
	}
	p := &StorageProducer{services: services,
		logger:          services.Logger.With().Str("producer", "storage").Logger(),
		ctx:             ctx,
		cancel:          cancel,
		taskChan:        make(chan *TaskSubmitted, poolSize),
		poolSize:        poolSize,
		goroutineRunner: goroutineRunner{},
	}
	return p
}

func (p *StorageProducer) Name() string { return "storage" }

func (p *StorageProducer) TaskChannel() <-chan *TaskSubmitted { return p.taskChan }

func (p *StorageProducer) Start(ctx context.Context) error {
	p.Go(p.storageQueuePollerLoop)
	return nil
}

func (p *StorageProducer) Stop() {
	p.stopOnce.Do(func() {
		p.logger.Info().Msg("storage producer stopping")
		if p.cancel != nil {
			p.cancel()
		}
		p.Wait()
		close(p.taskChan)
		p.logger.Info().Msg("stopped")
	})
}

func (p *StorageProducer) storageQueuePollerLoop() {
	p.logger.Info().Msg("storage queue poller loop started")
	defer p.logger.Info().Msg("storage queue poller loop stopped")
	emptyPollInt, errPollInt, backpressureInt := 1*time.Second, 5*time.Second, 100*time.Millisecond
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}
		if len(p.taskChan) >= p.poolSize {
			select {
			case <-time.After(backpressureInt):
				continue
			case <-p.ctx.Done():
				return
			}
		}
		taskId, txHash, err := p.services.TaskStorage.PopTask(p.services.Config.PopTaskRandom)
		if err != nil {
			sleepDur := errPollInt
			if errors.Is(err, sql.ErrNoRows) {
				sleepDur = emptyPollInt
			} else {
				p.logger.Error().Err(err).Msg("storage queue error")
			}
			select {
			case <-time.After(sleepDur):
				continue
			case <-p.ctx.Done():
				return
			}
		}
		ts := &TaskSubmitted{TaskId: taskId, TxHash: txHash}
		taskStr := taskId.String()
		p.logger.Debug().Str("task", taskStr).Msg("popped task from storage")
		select {
		case p.taskChan <- ts:
		case <-p.ctx.Done():
			p.logger.Warn().Str("task", taskStr).Msg("poller stopping, requeueing")
			_, _ = p.services.TaskStorage.RequeueTaskIfNoCommitmentOrSolution(taskId)
			return
		}
	}
}

// ListenerProducer: Relays single tasks from EventManager (Stateless)
type ListenerProducer struct{ eventMgr *EventManager }

func NewListenerProducer(eventMgr *EventManager) *ListenerProducer {
	return &ListenerProducer{eventMgr: eventMgr}
}

func (p *ListenerProducer) Name() string { return "listener" }

func (p *ListenerProducer) TaskChannel() <-chan *TaskSubmitted { return p.eventMgr.singleTasksOut }

func (p *ListenerProducer) Start(ctx context.Context) error { return nil } // No-op

func (p *ListenerProducer) Stop() {} // No-op

// SolutionEventProducer: Specialized producer for SolutionSamplerStrategy (Stateful)
type SolutionEventProducer struct {
	logger       zerolog.Logger
	cache        TaskInfoCache                          // Interface for looking up TxHash
	solutionSink <-chan *engine.EngineSolutionSubmitted // Source of events
	ctx          context.Context                        // Context from the owning strategy
	cancel       context.CancelFunc                     // To stop internal loops
	stopOnce     sync.Once
	taskChan     chan *TaskSubmitted // Output channel for validation tasks
	bufferSize   int

	// Reservoir Sampling state
	sampleMutex       sync.Mutex
	tasksSamples      []*TaskSubmitted
	sampleIndex       int
	maxTaskSampleSize int
	dispatchTicker    *time.Ticker

	goroutineRunner // Manages internal loops (processor, dispatcher)
}

func NewSolutionEventProducer(
	strategyCtx context.Context,
	log zerolog.Logger,
	solutionEventSink <-chan *engine.EngineSolutionSubmitted, // Comes from bs.verifySink
	cache TaskInfoCache, // bs implements this
	bufferSize int,
	sampleSize int,
	dispatchInterval time.Duration,
) *SolutionEventProducer {
	ctx, cancel := context.WithCancel(strategyCtx)
	if bufferSize <= 0 {
		bufferSize = 20
	}
	if sampleSize <= 0 {
		sampleSize = 50
	}
	if dispatchInterval <= 0 {
		dispatchInterval = 1 * time.Minute
	}

	p := &SolutionEventProducer{
		logger:            log.With().Str("producer", "solutionevent").Logger(),
		cache:             cache,
		solutionSink:      solutionEventSink,
		ctx:               ctx,
		cancel:            cancel,
		taskChan:          make(chan *TaskSubmitted, bufferSize),
		bufferSize:        bufferSize,
		maxTaskSampleSize: sampleSize,
		tasksSamples:      make([]*TaskSubmitted, 0, sampleSize), // Initial slice
		dispatchTicker:    time.NewTicker(dispatchInterval),
		goroutineRunner:   goroutineRunner{}, // Own runner for internal loops
	}
	return p
}

func (p *SolutionEventProducer) Name() string { return "solutionevent" }

func (p *SolutionEventProducer) TaskChannel() <-chan *TaskSubmitted { return p.taskChan }

func (p *SolutionEventProducer) Start(ctx context.Context) error {
	p.logger.Info().Msg("starting")
	if p.solutionSink == nil {
		return errors.New("solution event sink is nil")
	}
	if p.cache == nil {
		return errors.New("task info cache is nil")
	}
	p.Go(p.processSolutionSubmittedEvents) // Reads from sink, adds to sample pool
	p.Go(p.sampleDispatcherLoop)           // Periodically dispatches samples to taskChan
	return nil
}

func (p *SolutionEventProducer) Stop() {
	p.stopOnce.Do(func() {
		p.logger.Info().Msg("stopping")
		if p.dispatchTicker != nil {
			p.dispatchTicker.Stop()
		}
		if p.cancel != nil {
			p.cancel()
		} // Signal internal loops
		p.Wait()          // Wait for loops
		close(p.taskChan) // Close output channel
		p.logger.Info().Msg("stopped")
	})
}

// processSolutionSubmittedEvents: Reads from sink, looks up hash, performs reservoir sampling.
func (p *SolutionEventProducer) processSolutionSubmittedEvents() {
	p.logger.Info().Msg("solution processing loop started")
	defer p.logger.Info().Msg("solution processing loop stopped")
	for {
		select {
		case <-p.ctx.Done():
			return
		case event, ok := <-p.solutionSink:
			if !ok {
				p.logger.Info().Msg("solution sink closed")
				return
			}
			if event == nil {
				continue
			}
			taskIdStr := task.TaskId(event.Task).String()
			eventLog := p.logger.With().Str("task", taskIdStr).Logger()
			// Use the cache interface to get TxHash
			txHash, found := p.cache.GetTxHash(taskIdStr)
			if !found {
				eventLog.Warn().Msg("txhash not found in cache for solution event, cannot sample")
				continue
			}
			ts := &TaskSubmitted{TaskId: event.Task, TxHash: txHash}
			// Reservoir Sampling Logic
			p.sampleMutex.Lock()
			p.sampleIndex++
			if len(p.tasksSamples) < p.maxTaskSampleSize {
				p.tasksSamples = append(p.tasksSamples, ts)
				eventLog.Debug().Int("samples", len(p.tasksSamples)).Msg("added solution to sample pool")
			} else {
				j := rand.Intn(p.sampleIndex)
				if j < p.maxTaskSampleSize {
					p.tasksSamples[j] = ts
					eventLog.Debug().Int("replaced", j).Msg("replaced solution in sample pool")
				}
			}
			p.sampleMutex.Unlock()
		}
	}
}

// sampleDispatcherLoop: Periodically sends collected samples to output channel.
func (p *SolutionEventProducer) sampleDispatcherLoop() {
	p.logger.Info().Msg("sample dispatcher loop started")
	defer p.logger.Info().Msg("sample dispatcher loop stopped")
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.dispatchTicker.C:
			var newSamples []*TaskSubmitted
			p.sampleMutex.Lock()
			if len(p.tasksSamples) > 0 {
				newSamples = make([]*TaskSubmitted, len(p.tasksSamples))
				copy(newSamples, p.tasksSamples)
				p.tasksSamples = p.tasksSamples[:0]
				p.sampleIndex = 0
				p.logger.Info().Int("count", len(newSamples)).Msg("dispatching collected samples")
			}
			p.sampleMutex.Unlock()
			if len(newSamples) == 0 {
				continue
			} // Nothing to dispatch

			// Drain existing + Add new logic (simple prepend, trim if needed)
			oldWaitingTasks := make([]*TaskSubmitted, 0, p.bufferSize)
		drainingLoop:
			for {
				select {
				case task, ok := <-p.taskChan:
					if !ok {
						break drainingLoop
					}
					oldWaitingTasks = append(oldWaitingTasks, task)
				default:
					break drainingLoop
				}
			}
			combinedTasks := append(newSamples, oldWaitingTasks...) // New samples have higher priority
			numToKeep := len(combinedTasks)
			droppedCount := 0
			if numToKeep > p.bufferSize {
				droppedCount = numToKeep - p.bufferSize
				combinedTasks = combinedTasks[:p.bufferSize]
				numToKeep = p.bufferSize
				p.logger.Warn().Int("dropped", droppedCount).Msg("validation buffer full, oldest dropped")
			}
			// Refill channel
			for _, ts := range combinedTasks {
				select {
				case p.taskChan <- ts:
				case <-p.ctx.Done():
					return
				}
			} // Push back, respect context
		}
	}
}

// Strategy Implementationsf

// BulkMineStrategy: Solves storage tasks, single tasks, and baseline validation.
type BulkMineStrategy struct {
	bs       *baseStrategy
	eventMgr *EventManager
	producer TaskProducer
}

func NewBulkMineStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*BulkMineStrategy, error) {
	bs, err := newBaseStrategy(appCtx, services, submitter, gpuPool, "bulkmine")
	if err != nil {
		return nil, err
	}
	eventMgr := NewEventManager(bs)
	producer := NewStorageProducer(bs.ctx, services, bs.numWorkers*2)
	return &BulkMineStrategy{bs: bs, eventMgr: eventMgr, producer: producer}, nil
}

func (s *BulkMineStrategy) Name() string { return s.bs.strategyName }

func (s *BulkMineStrategy) Start() error {
	s.bs.logger.Info().Msg("starting BulkMineStrategy...")
	// 1. Start the submitter/batch manager, enabling batch processing
	if err := s.bs.submitter.Start(s.bs.ctx, true); err != nil { // Enable batch processing
		s.bs.logger.Error().Err(err).Msg("failed to start submitter/batch manager")
		return err
	}
	// 2. Start EventManager
	s.eventMgr.Start()
	// 3. Start workers
	return s.bs.start(s.producer, s.handleStorageTask, s.bs.handleValidationTask, s.eventMgr.singleTasksOut, s.bs.sampledValidationChan)
}

func (s *BulkMineStrategy) Stop() {
	s.bs.logger.Info().Msgf("stopping %s...", s.Name())
	s.producer.Stop()
	s.eventMgr.Stop()
	s.bs.Stop()
	s.bs.submitter.Stop() // Stop submitter
	s.bs.logger.Info().Msgf("%s stopped", s.Name())
}

func (s *BulkMineStrategy) handleStorageTask(workerId int, gpu *task.GPU, ts *TaskSubmitted, taskCtx context.Context) {
	workerLogger := s.bs.logger.With().Int("worker", workerId).Int("GPU", gpu.ID).Str("task", task.TaskId(ts.TaskId).String()).Logger()
	workerLogger.Info().Msg("processing storage task")
	gpu.SetStatus("Mining Storage Task")
	taskId := task.TaskId(ts.TaskId)

	// Helper function to requeue task if necessary
	requeueTask := func(reason string) {
		workerLogger.Warn().Str("reason", reason).Msg("requeueing storage task")
		// Attempt requeue, log error if it fails
		_, err := s.bs.services.TaskStorage.RequeueTaskIfNoCommitmentOrSolution(taskId)
		if err != nil {
			workerLogger.Error().Err(err).Str("task", task.TaskId(ts.TaskId).String()).Msg("failed to requeue task")
		}
	}

	// 1. Get Parameters (use caching decoder)
	params, err := s.bs.decodeTransactionWithCache(ts.TxHash)
	if err != nil || params == nil {
		// Decide whether to requeue based on error type
		if errors.Is(err, ErrTxDecodePermanent) || errors.Is(err, ethereum.NotFound) {
			// Permanent decode error or tx vanished - don't requeue
			workerLogger.Error().Err(err).Msg("permanent decode error or tx not found, dropping storage task")
		} else {
			// Assume transient error (RPC issue, temporary context cancel, etc.) - requeue
			requeueTask(fmt.Sprintf("transient decode error: %v", err))
		}
		gpu.SetStatus("Error") // Mark error status
		return
	}

	// 2. Solve and Submit
	solveStart := time.Now()
	// Call solveTask with submitImmediately: false
	cid, err := s.bs.solveTask(taskCtx, ts.TaskId, params, gpu, false, false /* submitImmediately */)
	solveElapsed := time.Since(solveStart)
	taskFailed := err != nil
	taskSkipped := cid == nil && err == nil

	// 3. Update Status & Handle Outcome
	if taskFailed && gpu.IsEnabled() && !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		gpu.SetStatus("Error")
	} else {
		gpu.SetStatus("Idle")
	}

	if taskFailed {
		// reqeue if context was cancelled or timed out
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			requeueTask(fmt.Sprintf("solve failed: %v", err))
		} else {
			workerLogger.Warn().Err(err).Msg("storage task solve failed, not requeueing")
		}
	} else if taskSkipped {
		workerLogger.Info().Msg("storage task skipped")
	} else {
		workerLogger.Info().Str("cid", "0x"+hex.EncodeToString(cid)).Str("elapsed", solveElapsed.String()).Msg("storage task solved successfully")
		s.bs.gpuPool.AddSolveTime(solveElapsed)
		// Do not requeue successful tasks
	}
}

// AutoMineStrategy: Inherits BulkMine.
type AutoMineStrategy struct{ *BulkMineStrategy }

func NewAutoMineStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*AutoMineStrategy, error) {
	if !services.Config.BatchTasks.Enabled {
		services.Logger.Warn().Msg("automine started but BatchTasks.Enabled is false in config")
		// Continue, but log warning
	}

	// Create the underlying BulkMineStrategy
	bulkStrategy, err := NewBulkMineStrategy(appCtx, services, submitter, gpuPool)
	if err != nil {
		return nil, fmt.Errorf("failed to create underlying bulkmine strategy for automine: %w", err)
	}

	// Override logger and strategy name for clarity
	bulkStrategy.bs.logger = services.Logger.With().Str("strategy", "automine").Logger()
	bulkStrategy.bs.strategyName = "automine"

	return &AutoMineStrategy{BulkMineStrategy: bulkStrategy}, nil
}

// ListenStrategy: Solves single tasks and baseline validation.
type ListenStrategy struct {
	bs       *baseStrategy
	eventMgr *EventManager
	producer TaskProducer
}

func NewListenStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*ListenStrategy, error) {
	bs, err := newBaseStrategy(appCtx, services, submitter, gpuPool, "listen")
	if err != nil {
		return nil, err
	}
	eventMgr := NewEventManager(bs)
	producer := NewListenerProducer(eventMgr)
	return &ListenStrategy{bs: bs, eventMgr: eventMgr, producer: producer}, nil
}

func (s *ListenStrategy) Name() string { return s.bs.strategyName }

func (s *ListenStrategy) Start() error {
	s.bs.logger.Info().Msg("starting ListenStrategy...")
	// 1. Start the submitter/batch manager, disabling batch processing
	if err := s.bs.submitter.Start(s.bs.ctx, false); err != nil { // Disable batch processing
		s.bs.logger.Error().Err(err).Msg("failed to start submitter/batch manager")
		return err
	}
	// 2. Start EventManager
	s.eventMgr.Start()
	// 3. Start workers
	return s.bs.start(nil, nil, s.bs.handleValidationTask, s.eventMgr.singleTasksOut, s.bs.sampledValidationChan)
}

func (s *ListenStrategy) Stop() {
	s.bs.logger.Info().Msgf("stopping %s strategy...", s.Name())
	s.eventMgr.Stop()
	s.bs.Stop()
	s.bs.submitter.Stop() // Stop submitter
	s.bs.logger.Info().Msgf("%s strategy stopped", s.Name())
}

// SolutionSamplerStrategy: ONLY does specialized sampling and validation.
type SolutionSamplerStrategy struct {
	bs       *baseStrategy
	producer TaskProducer
}

func NewSolutionSamplerStrategy(appCtx context.Context, services *Services, submitter SolutionSubmitter, gpuPool *GPUPool) (*SolutionSamplerStrategy, error) {
	bs, err := newBaseStrategy(appCtx, services, submitter, gpuPool, "solutionsampler")
	if err != nil {
		return nil, err
	}
	// Note: Does NOT create EventManager as it doesn't need single task routing or baseline sampling.
	// It relies on baseStrategy owning the verifySink and providing the cache lookup.
	// TODO: Configurable sample size/interval
	producer := NewSolutionEventProducer(bs.ctx, bs.logger, bs.verifySink, bs, bs.numWorkers*2, 50, 1*time.Minute)
	return &SolutionSamplerStrategy{bs: bs, producer: producer}, nil
}

func (s *SolutionSamplerStrategy) Name() string { return s.bs.strategyName }

func (s *SolutionSamplerStrategy) Start() error {
	s.bs.logger.Info().Msg("starting SolutionSamplerStrategy...")
	// 1. Start the submitter/batch manager, disabling batch processing
	if err := s.bs.submitter.Start(s.bs.ctx, false); err != nil { // Disable batch processing
		s.bs.logger.Error().Err(err).Msg("failed to start submitter/batch manager")
		return err
	}
	// Does NOT start EventManager
	// Start its own specialized producer
	if err := s.producer.Start(s.bs.ctx); err != nil {
		s.bs.logger.Error().Err(err).Msg("failed start solution event producer")
		return err
	}
	// Start workers: Only consumes from its own producer for validation.
	return s.bs.start(s.producer, s.bs.handleValidationTask, nil, nil, nil) // Pass nil for high-prio and baseline validation chans
}

func (s *SolutionSamplerStrategy) Stop() {
	s.bs.logger.Info().Msgf("stopping %s strategy...", s.Name())
	s.producer.Stop()     // Stop its specific producer
	s.bs.Stop()           // Stop base (workers/subs)
	s.bs.submitter.Stop() // Stop submitter
	s.bs.logger.Info().Msgf("%s strategy stopped", s.Name())
}
