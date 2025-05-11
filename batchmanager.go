package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gobius/account"
	"gobius/bindings/arbiusrouterv1"
	"gobius/bindings/bulktasks"
	"gobius/bindings/engine"
	"gobius/client"
	task "gobius/common"
	"gobius/config"
	"gobius/storage"
	"gobius/utils"
	"log"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/rs/zerolog"
)

// Constants for gas calculation and profit estimation
const (
	profitEstimateBatchSize = 200

	// below values are from onchain analysis using the 99% quantile of the gas used for the various functions
	// during very high inflation of gas intrinics out of gas errors are still possible
	claimTasksGasPerItem       = 100_828
	submitTasksGasPerItem      = 140_477
	signalCommitmentGasPerItem = 64_544
	submitSolutionGasPerItem   = 202_125

	baseGasLimitForClaimTasks        = 3_500_000
	baseGasLimitForSubmitTasks       = 2_500_000
	baseGasLimitForSubmitSolutions   = 3_500_000
	baseGasLimitForSignalCommitments = 1_500_000

	gasPriceAdjustmentFactor = 1_000_000_000.0

	// BulkClaimIncentiveGasPerItem is an estimated gas cost per task in a bulk claim.
	// Needs empirical measurement, starting with a guess based on single claim.
	BulkClaimIncentiveGasPerItem = 120000

	// BaseGasLimitForBulkClaimIncentive is the base gas limit for a bulk claim transaction.
	// Needs empirical measurement, starting with a guess.
	BaseGasLimitForBulkClaimIncentive = 400000
)

type CacheItem struct {
	Value      any
	LastUpdate time.Time
}

type BulkClaimData struct {
	TaskID       task.TaskId
	Signatures   []arbiusrouterv1.Signature
	ClaimAccount *account.Account // Store the account determined for this specific claim
}

type Cache struct {
	ttl   time.Duration
	items map[string]*CacheItem
	mu    sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl:   ttl,
		items: make(map[string]*CacheItem),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found { //|| time.Since(item.LastUpdate) > c.ttl {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Value:      value,
		LastUpdate: time.Now(),
	}
}

type BatchSolution struct {
	taskID task.TaskId
	cid    []byte
}

type BatchTransactionManager struct {
	services                      *Services
	cumulativeGasUsed             *GasMetrics // tracks how much gas we've spent
	signalCommitmentEvent         common.Hash
	solutionSubmittedEvent        common.Hash
	taskSubmittedEvent            common.Hash
	solutionClaimedEvent          common.Hash
	rewardsPaidEvent              common.Hash
	incentiveClaimedEvent         common.Hash
	commitments                   [][32]byte
	solutions                     []BatchSolution
	encodedTaskData               []byte
	taskAccounts                  []*account.Account
	minClaimSolutionTime          uint64
	minContestationVotePeriodTime uint64
	cache                         *Cache
	ipfsClaimBackoff              map[task.TaskId]int64 // New: Map to track backoff times
	ipfsClaimBackoffMutex         sync.Mutex            // New: Mutex for the backoff map
	goroutineRunner
	sync.Mutex
}

func NewBatchTransactionManager(services *Services, ctx context.Context) (*BatchTransactionManager, error) {

	// TODO: move these out of here!
	engineAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	// Get the event SolutionClaimed topic ID
	signalCommitmentEvent := engineAbi.Events["SignalCommitment"].ID
	solutionSubmittedEvent := engineAbi.Events["SolutionSubmitted"].ID
	taskSubmittedEvent := engineAbi.Events["TaskSubmitted"].ID

	// Get the event SolutionClaimed topic ID
	solutionClaimedEvent := engineAbi.Events["SolutionClaimed"].ID
	rewardsPaidEvent := engineAbi.Events["RewardsPaid"].ID

	arbiusRouterAbi, abiErr := arbiusrouterv1.ArbiusRouterV1MetaData.GetAbi()
	if abiErr != nil {
		panic("error getting arbiusrouterv1 abi")
	}

	incentiveClaimedEvent := arbiusRouterAbi.Events["IncentiveClaimed"].ID

	minClaimSolutionTimeBig, err := services.Engine.Engine.MinClaimSolutionTime(nil)
	if err != nil {
		return nil, err
	}

	//fmt.Println("minClaimSolutionTimeBig", minClaimSolutionTimeBig.String())

	minContestationVotePeriodTimeBig, err := services.Engine.Engine.MinContestationVotePeriodTime(nil)
	if err != nil {
		return nil, err
	}

	sampleRate, err := time.ParseDuration(services.Config.Miner.MetricsSampleRate)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := NewMetricsManager(ctx, sampleRate)

	var encodedData []byte

	encodedData, err = engineAbi.Pack("submitTask", services.AutoMineParams.Version, services.AutoMineParams.Owner, services.AutoMineParams.Model, services.AutoMineParams.Fee, services.AutoMineParams.Input)
	if err != nil {
		return nil, err
	}

	var accounts []*account.Account

	if len(services.Config.BatchTasks.PrivateKeys) > 0 {

		for _, pk := range services.Config.BatchTasks.PrivateKeys {

			account, err := account.NewAccount(pk, services.OwnerAccount.Client, ctx, services.Config.Blockchain.CacheNonce, services.Logger)
			if err != nil {
				return nil, err
			}
			account.UpdateNonce()

			accounts = append(accounts, account)
		}
	} else {
		accounts = append(accounts, services.OwnerAccount)
	}

	cache := NewCache(time.Duration(120) * time.Second)

	btm := &BatchTransactionManager{
		services:                      services,
		cache:                         cache,
		cumulativeGasUsed:             cumulativeGasUsed,
		signalCommitmentEvent:         signalCommitmentEvent,
		solutionSubmittedEvent:        solutionSubmittedEvent,
		taskSubmittedEvent:            taskSubmittedEvent,
		solutionClaimedEvent:          solutionClaimedEvent,
		rewardsPaidEvent:              rewardsPaidEvent,
		incentiveClaimedEvent:         incentiveClaimedEvent,
		encodedTaskData:               encodedData,
		taskAccounts:                  accounts,
		minClaimSolutionTime:          minClaimSolutionTimeBig.Uint64(),
		minContestationVotePeriodTime: minContestationVotePeriodTimeBig.Uint64(),
		ipfsClaimBackoff:              make(map[task.TaskId]int64), // Initialize the map
	}

	return btm, nil
}

// implementation for the metrics interface (WIP/TODO fill out)
func (tm *BatchTransactionManager) GetCurrentReward() float64 {
	reward, found := tm.cache.Get("reward")

	if !found {
		return math.NaN()
	}

	return reward.(float64)
}

func (tm *BatchTransactionManager) GetTotalTasks() int64 {
	return 11111
}

func (tm *BatchTransactionManager) GetClaims() int64 {
	return 200
}
func (tm *BatchTransactionManager) GetSolutions() int64 {
	return 300

}
func (tm *BatchTransactionManager) GetCommitments() int64 {
	return 400
}

func (tm *BatchTransactionManager) GetValidatorInfo() string {
	s := ""
	for _, v := range tm.services.Validators.validators {
		s += v.ValidatorAddress().String() + ": 100Aius\n"
	}

	return s
}

// TODO: refactor this to be more efficient and not so messy
func (tm *BatchTransactionManager) calcProfit(basefee *big.Int) (float64, float64, float64, float64, float64, error) {
	var err error
	var basePrice, ethPrice float64

	modelId := tm.services.AutoMineParams.Model
	taskFee := tm.services.AutoMineParams.Fee

	if basefee == nil {
		basefee, err = tm.services.OwnerAccount.Client.GetBaseFee()
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not get basefee!")
			return 0, 0, 0, 0, 0, err
		}
	}

	basefeeinEth := Eth.ToFloat(basefee)
	basefeeinGwei := basefeeinEth * 1000000000

	// Use the PriceOracle interface to get prices
	basePrice, ethPrice, err = tm.services.OracleProvider.GetPrices()
	if err != nil {
		if tm.services.Config.BaseConfig.TestnetType > 0 {
			tm.services.Logger.Warn().Msg("oracle failed, using default testnet prices (30, 2000)")
			basePrice, ethPrice = 30, 2000
		} else {
			// TODO:consider uing last known good price
			tm.services.Logger.Error().Err(err).Msg("could not get prices from oracle!")
			return 0, 0, 0, 0, 0, err
		}
		err = nil
	}

	tm.cache.Set("base_price", basePrice)
	tm.cache.Set("eth_price", ethPrice)

	submitTasksBatchUSD := 0.0
	submitTasksBatch := (submitTasksGasPerItem * basefeeinEth * profitEstimateBatchSize)
	submitTasksBatchUSD = submitTasksBatch * ethPrice

	signalCommitmentBatch := (signalCommitmentGasPerItem * basefeeinEth * profitEstimateBatchSize)
	signalCommitmentBatchUSD := signalCommitmentBatch * ethPrice

	submitSolutionBatch := (submitSolutionGasPerItem * basefeeinEth * profitEstimateBatchSize)
	submitSolutionBatchUSD := submitSolutionBatch * ethPrice

	claimTasksUSD := 0.0
	claimTasks := (claimTasksGasPerItem * basefeeinEth * profitEstimateBatchSize)
	claimTasksUSD = claimTasks * ethPrice

	totalCostPerBatchUSD := (submitTasksBatchUSD + signalCommitmentBatchUSD + submitSolutionBatchUSD + claimTasksUSD)

	modelReward, err := tm.services.Engine.GetModelReward(modelId)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("could not get model reward!")
		return 0, 0, 0, 0, 0, err
	}
	// for now we are taking 10% of the reward as a fee (this factors in the task owner reward and the treasury reward)
	rewardInAIUS := tm.services.Config.BaseConfig.BaseToken.ToFloat(modelReward) * 0.9

	//rewardTotal := new(big.Int).Sub(modelReward, taskFee)

	rewardInAIUSMinusFee := rewardInAIUS - tm.services.Config.BaseConfig.BaseToken.ToFloat(taskFee)

	tm.cache.Set("reward", rewardInAIUS)

	rewardsPerBatchUSD := rewardInAIUSMinusFee * basePrice * profitEstimateBatchSize

	profit := rewardsPerBatchUSD - totalCostPerBatchUSD

	tm.cumulativeGasUsed.profitEMA.Add(profit)

	tm.services.Logger.Info().
		Str("base_model_reward", fmt.Sprintf("%.8g", rewardInAIUS)).
		Str("model_reward_minus_fee", fmt.Sprintf("%.8g", rewardInAIUSMinusFee)).
		Str("eth_in_usd", fmt.Sprintf("%.4g$", ethPrice)).
		Str("aius_in_usd", fmt.Sprintf("%.4g$", basePrice)).
		Msg("💰 model reward and eth/aius price")

	tm.services.Logger.Info().
		Str("costs_in_usd", fmt.Sprintf("%.4g$", totalCostPerBatchUSD)).
		Str("rewards_in_usd", fmt.Sprintf("%.4g$", rewardsPerBatchUSD)).
		Str("profit_per_batch", fmt.Sprintf("%.4g$", profit)).
		Str("base_fee", fmt.Sprintf("%.8g", basefeeinGwei)).
		Str("profit_metrics", tm.cumulativeGasUsed.profitEMA.String()).
		Msg("💰 batch profits")

	return profit, basefeeinGwei, rewardInAIUSMinusFee, ethPrice, basePrice, nil
}

// This should only be run when not performing other batched operations and just does batch claims
func (tm *BatchTransactionManager) batchClaimPoller(appQuit context.Context, pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()

	for {
		select {
		case <-appQuit.Done():
			tm.services.Logger.Info().Msg("batch claimer shutting down")
			return
		case <-ticker.C:

			_, totalClaims, err := tm.services.TaskStorage.TotalSolutionsAndClaims()
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("failed to get total solutions")
				continue
			}

			tm.services.Logger.Info().Int64("claims", totalClaims).Msg("claims waiting to be processed")

			claimBatchSize := tm.services.Config.Claim.MaxClaims
			if claimBatchSize <= 0 {
				tm.services.Logger.Warn().Msgf("** claim batch size set to 0 so no claims will be made **")
			} else {
				claims, _, err := tm.services.TaskStorage.GetClaims(claimBatchSize)
				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("could not get claims from storage")
					continue
				}
				if len(claims) > 0 {

					if tm.services.Config.Claim.MaxGas > 0 {
						basefeeBig, err := tm.services.OwnerAccount.Client.GetBaseFee()
						if err != nil {
							tm.services.Logger.Error().Err(err).Msg("could not get basefee!")
							continue
						}
						// convert basefee to gwei
						basefeeinEth := Eth.ToFloat(basefeeBig)
						basefeeinGwei := basefeeinEth * 1000000000

						if basefeeinGwei > tm.services.Config.Claim.MaxGas {
							tm.services.Logger.Warn().Msgf("** base gas is too high to claim **")
							continue
						}
					}

					tm.processBulkClaim(tm.services.OwnerAccount, claims, tm.services.Config.Claim.MinClaims, claimBatchSize)
				}
			}
		}
	}
}

func (tm *BatchTransactionManager) processBatchBlockTrigger(appQuit context.Context) {
	var batchWG sync.WaitGroup

	headers := make(chan *types.Header)
	var newHeadSub ethereum.Subscription

	connectToHeaders := func() {
		var err error

		newHeadSub, err = tm.services.OwnerAccount.Client.Client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			tm.services.Logger.Fatal().Err(err).Msg("Failed to subscribe to new headers")
		}
	}

	connectToHeaders()
	maxBackoffHeader := time.Second * 30
	currentBackoffHeader := time.Second

	for {
		select {

		case <-appQuit.Done():
			tm.services.Logger.Info().Msg("delegated batch processor shutting down")
			return
		case <-headers:

			/*basefee, err := tm.services.OwnerAccount.Client.GetBaseFee()
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not get basefee!")
				continue
			}

			// convert basefee to gwei
			basefeeinEth := tm.services.Eth.ToFloat(basefee)
			// convert basefee to gwei
			basefeeinGwei := basefeeinEth * 1000000000

			minProfit := tm.services.Config.DelegatedMiner.MinProfit

			if basefeeinGwei <= minProfit {*/

			profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice, err := tm.calcProfit(nil)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not calculate profit, skipping batch")
				continue
			}

			//start := time.Now()
			tm.processBatch(appQuit, &batchWG, profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice)
			batchWG.Wait()
			//tm.services.Logger.Warn().Str("duration", time.Since(start).String()).Msg("batch processed")

			//}

		case err := <-newHeadSub.Err():
			if err == nil {
				continue
			}
			log.Printf("Error from newHeadSub: %v - redialling in: %s\n", err, currentBackoffHeader.String())
			newHeadSub.Unsubscribe()

			time.Sleep(currentBackoffHeader)
			currentBackoffHeader *= 2
			currentBackoffHeader += time.Duration(rand.Intn(1000)) * time.Millisecond
			if currentBackoffHeader > maxBackoffHeader {
				currentBackoffHeader = maxBackoffHeader
			}

			connectToHeaders()
		}

	}
}

func (tm *BatchTransactionManager) processValidatorStakePoller(appQuit context.Context, pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.ProcessValidatorsStakes()
		case <-appQuit.Done():
			tm.services.Logger.Info().Msg("validator stake processor shutting down")
			return
		}
	}
}
func (tm *BatchTransactionManager) processBatchPoller(appQuit context.Context, pollingtime time.Duration) {
	var batchWG sync.WaitGroup
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice, err := tm.calcProfit(nil)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not calculate profit, skipping batch")
				continue
			}

			start := time.Now()
			tm.processBatch(appQuit, &batchWG, profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice)
			batchWG.Wait()
			tm.services.Logger.Warn().Str("duration", time.Since(start).String()).Msg("batch processed")
		case <-appQuit.Done():
			tm.services.Logger.Info().Msg("delegated batch processor shutting down")
			return
		}
	}
}

func (tm *BatchTransactionManager) isIntrinsicGasTooHigh(ctx context.Context) (bool, error) {
	if !tm.services.Config.Miner.EnableIntrinsicGasCheck {
		return false, nil // Check disabled
	}

	baseline := tm.services.Config.Miner.IntrinsicGasBaseline
	if baseline == 0 {
		baseline = 22000 // Default baseline if not set
	}
	thresholdMultiplier := tm.services.Config.Miner.IntrinsicGasThresholdMultiplier
	if thresholdMultiplier <= 0 {
		thresholdMultiplier = 1.1 // Default multiplier if not set
	}

	opts := tm.services.OwnerAccount.GetOpts(0, nil, nil, nil)
	opts.NoSend = true
	opts.Value = big.NewInt(1)

	// Estimate gas for a simple transfer
	gasEstimate, err := tm.services.OwnerAccount.SendTransactionWithOpts(opts, &common.Address{}, []byte{})
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to estimate intrinsic gas cost")
		// Decide how to handle errors - fail open (assume normal) or fail closed (assume high)?
		// Let's fail open for now to avoid blocking unnecessarily due to transient estimation errors.
		return false, fmt.Errorf("failed to estimate intrinsic gas: %w", err)
	}

	gasLimit := gasEstimate.Gas()

	threshold := uint64(float64(baseline) * thresholdMultiplier)
	isHigh := gasLimit > threshold

	if isHigh {
		tm.services.Logger.Warn().
			Uint64("estimated_gas", gasLimit).
			Uint64("baseline", baseline).
			Float64("multiplier", thresholdMultiplier).
			Uint64("threshold", threshold).
			Msg("🚨 high intrinsic gas detected")
	} else {
		tm.services.Logger.Debug().
			Uint64("estimated_gas", gasLimit).
			Uint64("threshold", threshold).
			Msg("intrinsic gas check passed")
	}

	return isHigh, nil
}
func (tm *BatchTransactionManager) processBatch(
	appQuit context.Context,
	wg *sync.WaitGroup,
	profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice float64) {

	if err := appQuit.Err(); err != nil {
		return
	}

	paused, err := tm.services.Engine.IsPaused()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get paused status")
		return
	}

	if paused {
		tm.services.Logger.Warn().Msg("engine is paused, skipping batch")
		time.Sleep(10 * time.Second)
		return
	}

	if tm.services.Config.Miner.EnableIntrinsicGasCheck {
		isHigh, checkErr := tm.isIntrinsicGasTooHigh(appQuit) // Use appQuit context
		if checkErr != nil {
			// Logged within the check function, decide if we need to return or proceed
			tm.services.Logger.Error().Err(checkErr).Msg("error checking intrinsic gas, proceeding cautiously")
			// Potentially return here if strict safety is needed
		}
		if isHigh {
			tm.services.Logger.Warn().Msg("intrinsic gas cost too high") //, skipping batch processing")
			// Potentially add a short sleep here if desired
			// time.Sleep(5 * time.Second)
			//return // Skip the rest of the batch processing
		}
	}

	totalTasks, err := tm.services.TaskStorage.TotalTasks()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get total task count")
		return
	}

	// DEFAULTS:
	profitMode := tm.services.Config.Miner.ProfitMode
	minProfit := tm.services.Config.Miner.MinProfit
	maxProfit := tm.services.Config.Miner.MaxProfit
	claimMaxBatchSize := tm.services.Config.Claim.MaxClaims
	claimMinBatchSize := tm.services.Config.Claim.MinClaims

	minProfitFmt := fmt.Sprintf("%.4g$", minProfit)
	isProfitable := true
	usegwei := false
	switch profitMode {

	case "fixed":
		// just use a fixed min profit set above in defaults
	case "gwei":
		profitLevel = baseFee
		usegwei = true
		minProfitFmt = fmt.Sprintf("%.4g gwei", minProfit)
	case "ema":
		minProfit = tm.cumulativeGasUsed.profitEMA.Average()
		// just take profit at current ema
	case "randomauto":
		// use the ema and max value to decide a band
		maxProfit = tm.cumulativeGasUsed.profitEMA.MaxPrice() * 0.95
		minProfit = tm.cumulativeGasUsed.profitEMA.Average() * 0.9 //10% less than ema for breathing room
		if maxProfit < minProfit {
			maxProfit = minProfit
			minProfit = minProfit * 0.9
		}
		fallthrough

	default:
		// error
		tm.services.Logger.Error().Str("profit_mode", profitMode).Msg("invalid profit mode!")
		return
	}

	if usegwei {
		isProfitable = profitLevel <= minProfit
	} else {
		isProfitable = profitLevel >= minProfit
	}

	totalSolutions, totalClaims, err := tm.services.TaskStorage.TotalSolutionsAndClaims()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get total solutions")
		return
	}

	// Task creation decision logic
	makeTasks := false
	taskBatchCount := 0
	taskBatchSize := 0
	batchCfg := tm.services.Config.BatchTasks
	claimQueuePreventsTasks := false // Initialize here

	// 0. Check if Batch Task Creation is enabled globally
	if !batchCfg.Enabled {
		tm.services.Logger.Info().Msg("batch task creation is globally disabled (batchtasks.enabled=false)")
	} else {
		// 1. Check if claim queue prevents task creation (only if globally enabled)
		claimQueuePreventsTasks = batchCfg.MaxClaimQueue > 0 && totalClaims >= int64(batchCfg.MaxClaimQueue)
		if claimQueuePreventsTasks {
			tm.services.Logger.Info().Int64("claims", totalClaims).Int("max_claims", batchCfg.MaxClaimQueue).Msg("claim queue is full, task creation disabled")
		} else {
			// 2. Check standard task creation conditions (only if globally enabled and claim queue allows)
			shouldCreateStandardTasks := isProfitable &&
				totalTasks < int64(batchCfg.MinTasksInQueue) &&
				batchCfg.BatchSize > 0

			if shouldCreateStandardTasks {
				makeTasks = true
				taskBatchCount = batchCfg.NumberOfBatches
				taskBatchSize = batchCfg.BatchSize
				tm.services.Logger.Info().Msg("standard task creation conditions met")
			} else {
				// 3. Check hoard mode conditions (only if standard conditions not met, globally enabled, and claim queue allows)
				shouldActivateHoardMode := batchCfg.HoardMode &&
					baseFee <= batchCfg.HoardMinGasPrice &&
					totalTasks < int64(batchCfg.HoardMaxQueueSize)

				if shouldActivateHoardMode {
					makeTasks = true
					taskBatchCount = batchCfg.HoardModeNumberOfBatches
					taskBatchSize = batchCfg.HoardModeBatchSize
					tm.services.Logger.Warn().Msgf("** task hoard mode activated **")
				}
			}
		}
	}

	// Log final decision factors
	tm.services.Logger.Info().
		Bool("make_tasks", makeTasks).
		Bool("claim_queue_prevents", claimQueuePreventsTasks).
		Bool("is_profitable", isProfitable).
		Int64("total_tasks", totalTasks).
		Int("min_tasks_in_queue", batchCfg.MinTasksInQueue).
		Int("batch_size_to_use", taskBatchSize).   // Log the size actually being used
		Int("batch_count_to_use", taskBatchCount). // Log the count actually being used
		Msg("final task creation decision")

	if makeTasks {

		sendTasks := func(account *account.Account, wg *sync.WaitGroup) {
			defer wg.Done()

			// Calculate total fee required for the batch
			feePerTask := tm.services.AutoMineParams.Fee
			totalFee := new(big.Int).Mul(feePerTask, big.NewInt(int64(taskBatchSize)))

			// Get the account's AIUS balance
			aiusBalanceAsBig, err := tm.services.Basetoken.BalanceOf(nil, account.Address)
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("account", account.Address.String()).Msg("failed to get AIUS balance for task submission")
				return
			}

			aiusBalanceAsFloat := tm.services.Config.BaseConfig.BaseToken.ToFloat(aiusBalanceAsBig)

			// Check if balance is sufficient
			// TODO: make level this configurable e.g. some min balance level
			if aiusBalanceAsFloat < tm.services.Config.ValidatorConfig.MinBasetokenThreshold {
				tm.services.Logger.Warn().
					Str("account", account.Address.String()).
					Str("balance", fmt.Sprintf("%.8g", aiusBalanceAsFloat)).
					Str("required", fmt.Sprintf("%.8g", tm.services.Config.BaseConfig.BaseToken.ToFloat(totalFee))).
					Int("batch_size", taskBatchSize).
					Msgf("** task fee will exceed min balance threshold of %.8g, skipping task batch **", tm.services.Config.ValidatorConfig.MinBasetokenThreshold)
				return
			}

			// log the amount of aius reqiured for this batch in a human readable format
			feeTransferAsfloat := tm.services.Config.BaseConfig.BaseToken.ToFloat(totalFee)

			tm.services.Logger.Warn().Str("fee_transfer", fmt.Sprintf("%.4g", feeTransferAsfloat)).Int("batch_size", taskBatchSize).Str("account", account.Address.String()).Msgf("** task queue is low - sending batch **")
			receipt, err := tm.BulkTasks(account, taskBatchSize)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error sending batch tasks")
			} else {
				tm.services.Logger.Info().Int("batch_size", taskBatchSize).Str("txhash", receipt.TxHash.String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch tasks tx accepted!")
			}
		}

		switch tm.services.Config.BatchTasks.BatchMode {
		case "normal":
			//  get random index into taskAccounts
			accountIndex := rand.Intn(len(tm.taskAccounts))
			for i := 0; i < taskBatchCount; i++ {

				// get the next available account
				account := tm.taskAccounts[accountIndex%len(tm.taskAccounts)]
				accountIndex++

				wg.Add(1)
				if tm.services.Config.Miner.ConcurrentBatches {
					go sendTasks(account, wg)
					if len(tm.taskAccounts) == 1 {
						time.Sleep(500 * time.Millisecond)
					}
				} else {
					//if err := appQuit.Err(); err != nil {
					//return
					//}
					sendTasks(account, wg)
				}
			}
		case "account":
			for _, acc := range tm.taskAccounts {
				wg.Add(1)
				//if tm.services.Config.DelegatedMiner.ConcurrentBatches {
				go sendTasks(acc, wg)
				//} else {
				//sendTasks(acc, wg)
				//}
			}
		default:
			tm.services.Logger.Error().Str("batch_mode", tm.services.Config.BatchTasks.BatchMode).Msg("invalid task batch mode")
		}

		if tm.services.Config.BatchTasks.OnlyTasks {
			tm.services.Logger.Info().Msg("only tasks sent, skipping other batch operations")
			return
		}
	}
	totalCommitments, err := tm.services.TaskStorage.TotalCommitments()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get total commitments")
		return
	}

	totalTasksCount, totalTasksGasFloat, err := tm.services.TaskStorage.GetTotalTasksGas()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get total tasks gas")
		return
	}

	// auto mine fee per task
	feePerTaskAsBig := tm.services.AutoMineParams.Fee
	feePerTaskAsFloat := tm.services.Config.BaseConfig.BaseToken.ToFloat(feePerTaskAsBig)
	totalQueued := float64(totalTasks + totalSolutions + totalCommitments + totalClaims)
	totalAiusOnFee := totalQueued * feePerTaskAsFloat
	totalAiusEarnings := totalQueued * rewardInAIUS // reward in aius factors in the task fee already duh
	totalSpentonGasInUSD := totalTasksGasFloat * ethPrice
	totalProfitOfQueueInUsd := totalAiusEarnings*basePrice - totalSpentonGasInUSD

	// total queue stats
	tm.services.Logger.Info().
		Int64("tasks", totalTasks).             // total tasks in queue
		Int64("tasks_count", totalTasksCount).  // total tasks in DB which might be higher due to stale tasks
		Int64("solutions", totalSolutions).     // total solutions in queue
		Int64("commitments", totalCommitments). // total commitments in queue
		Int64("claims", totalClaims).           // total claims in queue
		Msg("pending totals for batch queue")
	tm.services.Logger.Info().
		Str("total_aius", fmt.Sprintf("%.8g", totalAiusOnFee)).
		Str("total_gas_usd", fmt.Sprintf("%.4g$", totalSpentonGasInUSD)).
		Msg("total gas and aius spent on tasks in queue")
	tm.services.Logger.Info().
		Str("profit_aius", fmt.Sprintf("%.8g", totalAiusEarnings)).
		Str("profit_usd", fmt.Sprintf("%.4g$", totalProfitOfQueueInUsd)).
		Msg("approx total aius and usd profit from pending queue")

	if !isProfitable {
		tm.services.Logger.Info().Str("profit_mode", profitMode).Str("min_profit", minProfitFmt).Str("max_profit", fmt.Sprintf("%.4g", maxProfit)).Msg("not profitable to process batch")
		return
	}

	tm.services.Logger.Info().Str("profit_mode", profitMode).Str("min_profit", minProfitFmt).Msg("profit criteria met - processing batch")

	if claimMinBatchSize <= 0 || !tm.services.Config.Claim.Enabled {
		tm.services.Logger.Warn().Msgf("** claim batch disabled or min batch size set to 0 so no claims will be made **")
	} else {

		noOfClaimBatches := tm.services.Config.Claim.NumberOfBatches
		if noOfClaimBatches <= 0 {
			noOfClaimBatches = 1
		}

		claims, averageGas, err := tm.services.TaskStorage.GetClaims(claimMaxBatchSize * noOfClaimBatches)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not get keys from storage")
			return
		}

		totalCost := 0.0
		for _, task := range claims {
			totalCost += task.TotalCost
		}

		// calculate the cost of the claims in aius using magic numbers
		claimTasks := (claimTasksGasPerItem / gasPriceAdjustmentFactor * baseFee * float64(len(claims)))

		tm.services.Logger.Warn().Msgf("** CHEAPEST BATCH WE CAN SEND OUT **")
		tm.services.Logger.Warn().Msgf("** total cost of %d claims: %f  **", len(claims), totalCost)
		tm.services.Logger.Warn().Msgf("** average gas per task   : %f  **", averageGas)

		totalCost += claimTasks

		totalCostInUSD := totalCost * ethPrice //fmt.Sprintf("%0.4f$", totalCost*ethPrice)

		claimValue := rewardInAIUS * float64(len(claims)) * basePrice
		actualProfit := claimValue - totalCostInUSD

		tm.services.Logger.Warn().Msgf("**      total cost of mining batch : %0.4g$ (gas spent: %f)**", totalCostInUSD, totalCost)
		tm.services.Logger.Warn().Msgf("**                     batch value : %0.4g$ **", claimValue)
		tm.services.Logger.Warn().Msgf("**                          profit : %0.4g$ **", claimValue-totalCostInUSD)
		tm.services.Logger.Warn().Msgf("**                     profit/task : %0.4g$ **", (claimValue-totalCostInUSD)/float64(len(claims)))
		tm.services.Logger.Warn().Msgf("**********************************************")

		claimLen := len(claims)
		if claimLen > 0 {
			canClaim := false

			claimMinReward, err := tm.services.LeverOracle.MinClaimLever()
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not get minclaim lever from oracle!")
				return
			}

			if claimMinReward > 0 {
				if rewardInAIUS >= claimMinReward {
					tm.services.Logger.Warn().Msgf("** %.8g reward is >= claim min of reward %.8g, claim **", rewardInAIUS, claimMinReward)
					canClaim = true
				} else {
					tm.services.Logger.Warn().Msgf("** %.8g reward is below claim min reward of %.8g, skipping claim **", rewardInAIUS, claimMinReward)
					canClaim = false
				}
			} else if tm.services.Config.Claim.HoardMode {
				if int(totalClaims) < tm.services.Config.Claim.HoardMaxQueueSize {
					canClaim = false
					tm.services.Logger.Warn().Msgf("** claim hoard mode on, and queue length of %d is below threshold of %d - skipping claim **", int(totalClaims), tm.services.Config.Claim.HoardMaxQueueSize)
				} else {
					canClaim = true
					tm.services.Logger.Warn().Msgf("** claim hoard mode on, and queue length of %d is above threshold of %d - claiming **", int(totalClaims), tm.services.Config.Claim.HoardMaxQueueSize)
				}
			} else if tm.services.Config.Claim.MinBatchProfit > 0 {
				if actualProfit < tm.services.Config.Claim.MinBatchProfit {
					tm.services.Logger.Warn().Msgf("** batch profit of %.8g is below claim threshold of %.8g, skipping claim **", actualProfit, tm.services.Config.Claim.MinBatchProfit)
					canClaim = false
				} else {
					canClaim = true
					tm.services.Logger.Warn().Msgf("** batch profit of %.8g is above claim threshold of %.8g, claiming **", actualProfit, tm.services.Config.Claim.MinBatchProfit)
				}
			} else {
				canClaim = true
				tm.services.Logger.Warn().Msgf("** default claim conditions met, claiming **")
			}

			// if tm.services.Config.Claim.ClaimOnApproachMinStake && validatorBuffer < tm.services.Config.Claim.MinStakeBufferLevel {
			// 	tm.services.Logger.Warn().Msgf("** claim on approach to min stake enabled and validator buffer is below set min stake level **")
			// 	canClaim = true
			// }

			if canClaim {
				accountIndex := rand.Intn(len(tm.taskAccounts))
				sendBulkClaim := func(chunk storage.ClaimTaskSlice, account *account.Account, wg *sync.WaitGroup, batchno int) {
					defer wg.Done()
					tm.services.Logger.Info().Int("max_batch_size", claimMaxBatchSize).Int("batch_no", batchno+1).Str("address", account.Address.String()).Msgf("preparing claim batch")
					tm.processBulkClaim(account, chunk, claimMinBatchSize, claimMaxBatchSize)
				}

				for i, chunk := range claims.SplitIntoChunks(claimMaxBatchSize) {
					if err := appQuit.Err(); err != nil {
						return
					}
					// get the next available account
					account := tm.taskAccounts[accountIndex%len(tm.taskAccounts)]
					accountIndex++

					wg.Add(1)
					if tm.services.Config.Miner.ConcurrentBatches {
						go sendBulkClaim(chunk, account, wg, i)
						//time.Sleep(800 * time.Millisecond)
					} else {
						sendBulkClaim(chunk, account, wg, i)
					}
				}
			}
		}
	}

	deleteCommitments := func(_commitmentsToDelete []task.TaskId) error {
		if len(_commitmentsToDelete) > 0 {
			err := tm.services.TaskStorage.DeleteProcessedCommitments(_commitmentsToDelete)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error deleting commitment(s) from storage")
				return err
			}
			tm.services.Logger.Info().Msgf("deleted %d task commitments from storage that were committed on-chain", len(_commitmentsToDelete))
		}
		return nil
	}

	deleteSolutions := func(_solutionsToDelete []task.TaskId) error {
		if len(_solutionsToDelete) > 0 {
			err := tm.services.TaskStorage.DeleteProcessedSolutions(_solutionsToDelete)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error deleting solution(s) from storage")
				return err
			}
			tm.services.Logger.Info().Msgf("deleted %d task solutions from storage that were submitted on-chain", len(_solutionsToDelete))
		}
		return nil
	}

	// simplified this function and removed solution checks, these are now only done in the getSolutionBatchdata function
	// If a commitment is found to be on-chain, it is deleted from storage
	getCommitmentBatchdata := func(batchSize int, noChecks bool) (storage.TaskDataSlice, error) {
		var tasks storage.TaskDataSlice
		var err error

		tasks, err = tm.services.TaskStorage.GetPendingCommitments(batchSize)
		if err != nil {
			tm.services.Logger.Err(err).Msg("failed to get commitments from storage")
			return nil, err
		}

		if noChecks || len(tasks) == 0 { // Also return early if no tasks
			return tasks, nil
		}

		var commitmentsToDelete []task.TaskId
		var validTasksForBatch storage.TaskDataSlice
		var tasksToUpdateToStatus2 []task.TaskId // New list to track tasks needing status update

		// Collect all commitment hashes for a bulk call
		collectedCommitmentHashes := make([][32]byte, 0, len(tasks))
		taskMapByCommitment := make(map[[32]byte]storage.TaskData) // Helper to map results back
		for _, t := range tasks {
			// Ensure task has a commitment hash locally before checking on-chain
			if t.Commitment == ([32]byte{}) {
				tm.services.Logger.Warn().Str("taskid", t.TaskId.String()).Msg("Task in pending commitments has zero commitment hash, skipping.")
				continue // Skip tasks with invalid local state
			}
			collectedCommitmentHashes = append(collectedCommitmentHashes, t.Commitment)
			taskMapByCommitment[t.Commitment] = t
		}

		if len(collectedCommitmentHashes) == 0 {
			// All tasks might have had zero commitment hash
			return validTasksForBatch, nil
		}

		// Make a single bulk call to get commitment statuses
		var blockNumbers []*big.Int

		getCommitmentsOperation := func() (interface{}, error) {
			blocks, err := tm.services.BulkTasks.GetCommitments(nil, collectedCommitmentHashes)
			if err != nil {
				return nil, err
			}
			return blocks, nil
		}

		retryResult, err := utils.ExpRetry(tm.services.Logger, getCommitmentsOperation, 3, 1000)

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error checking on-chain commitments")
			return nil, err
		}

		// Type assert the result from ExpRetry
		var ok bool
		blockNumbers, ok = retryResult.([]*big.Int)
		if !ok {
			errInternal := errors.New("ExpRetry returned unexpected type for GetCommitments")
			tm.services.Logger.Error().Err(errInternal).Msg("internal error: type assertion failed for blockNumbers")
			return nil, errInternal
		}

		if len(blockNumbers) != len(collectedCommitmentHashes) {
			tm.services.Logger.Error().Msgf("bulk commitment check returned %d results for %d commitments, mismatch, skipping processing this batch", len(blockNumbers), len(collectedCommitmentHashes))
			return nil, errors.New("bulk commitment check result count mismatch")
		}

		for i, blockNumBig := range blockNumbers {
			commitmentHash := collectedCommitmentHashes[i]
			originalTask, found := taskMapByCommitment[commitmentHash]
			if !found {
				// Should not happen if logic is correct
				tm.services.Logger.Error().Str("commitment", hex.EncodeToString(commitmentHash[:])).Msg("Original task not found for commitment hash after bulk call")
				continue
			}

			blockNo := blockNumBig.Uint64()
			if blockNo > 0 {
				// Commitment already exists on-chain
				commitmentsToDelete = append(commitmentsToDelete, originalTask.TaskId)
				tasksToUpdateToStatus2 = append(tasksToUpdateToStatus2, originalTask.TaskId) // Mark for status update
			} else {
				// Commitment not on-chain, add to the batch to be sent
				validTasksForBatch = append(validTasksForBatch, originalTask)
			}
		}

		// Delete local commitments found to be on-chain
		if err := deleteCommitments(commitmentsToDelete); err != nil {
			tm.services.Logger.Error().Err(err).Msg("error deleting locally processed commitments that are on-chain")
			return nil, err
		}

		// Update status for tasks whose commitments were found on-chain
		// this simulates the behaviour of sendCommitments
		if len(tasksToUpdateToStatus2) > 0 {
			// Use the new function to update status only, preserving cost
			err = tm.services.TaskStorage.UpdateTaskStatusOnly(tasksToUpdateToStatus2, 2)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error updating task status to 2 for on-chain commitments")
				return nil, err
			} else {
				tm.services.Logger.Info().Msgf("updated %d tasks to committed status, matching on-chain commitment status", len(tasksToUpdateToStatus2))
			}
		}

		return validTasksForBatch, nil
	}

	getSolutionBatchdata := func(batchSize int, noChecks bool) (*Validator, storage.TaskDataSlice, error) {

		// map of validator to number of items we can send
		// loop through pending counts for each validator and do min(count, items we can send)
		// send the one with highest min
		solsPerVal, errTS := tm.services.TaskStorage.GetPendingSolutionsCountPerValidator()
		if errTS != nil {
			tm.services.Logger.Err(errTS).Msg("failed to get pending solutions count per validator from storage")
			return nil, nil, errTS
		}

		blockInfo, errbn := tm.services.OwnerAccount.Client.Client.BlockByNumber(context.Background(), nil)
		if errbn != nil {
			tm.services.Logger.Error().Err(errbn).Msg("failed to get latest block")
			return nil, nil, errbn
		}

		// Adjusted blockTime with a small safety margin to prevent hitting rate limits too soon
		adjustedBlockTime := time.Unix(int64(blockInfo.Time()), 0).Add(-1 * time.Second) // 1-second safety margin

		var validator *Validator = nil
		validatorHighestMin := int64(-1)
		for _, v := range tm.services.Validators.validators {
			lastSubmission, maxSols, err := v.MaxSubmissions(adjustedBlockTime) // Use adjustedBlockTime
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("validator", v.ValidatorAddress().String()).Msg("failed to get max submissions for validator, skipping")
				continue // Skip this validator if we can't determine max submissions
			}
			solsPending, found := solsPerVal[v.ValidatorAddress()]

			if found {
				minWeCanSend := min(maxSols, solsPending)
				tm.services.Logger.Info().Str("last_submission", lastSubmission.String()).Int64("max_submissions", maxSols).Int64("sols_pending", solsPending).Str("validator", v.ValidatorAddress().String()).Msg("validator info")

				if minWeCanSend > validatorHighestMin {
					validatorBuffer, err := v.GetValidatorStakeBuffer()
					if err != nil {
						tm.services.Logger.Error().Err(err).Msg("could not get validator stake buffer")
						return nil, nil, err
					}
					if tm.services.Config.Miner.PauseStakeBufferLevel > 0 && validatorBuffer < tm.services.Config.Miner.PauseStakeBufferLevel {
						tm.services.Logger.Warn().Float64("buffer", validatorBuffer).Msgf("** validator %s stake is at or below pause threshold **", v.ValidatorAddress().String())
					} else {
						validatorHighestMin = minWeCanSend
						validator = v
					}
				}
			}
		}

		if validator == nil {
			tm.services.Logger.Info().Msg("no solutions available to submit for any validator right now")
			return nil, nil, nil
		} else {
			tm.services.Logger.Info().Int64("highest_min", validatorHighestMin).Msgf("getting solutions for validator: %s", validator.ValidatorAddress().String())
		}

		var tasks storage.TaskDataSlice
		var err error

		batchSize = min(int(validatorHighestMin), batchSize)

		tasks, err = tm.services.TaskStorage.GetPendingSolutions(validator.ValidatorAddress(), batchSize)

		if err != nil {
			tm.services.Logger.Err(err).Msg("failed to get tasks from storage")
			return nil, nil, err
		}

		if noChecks || len(tasks) == 0 { // Also return early if no tasks
			return validator, tasks, nil
		}

		var commitmentsToDelete []task.TaskId
		var solutionsToDelete []task.TaskId
		var validTasksForBatch storage.TaskDataSlice
		var claimsToDelete []task.TaskId

		// Collect all task IDs for a bulk call
		collectedTaskIds := make([][32]byte, 0, len(tasks))
		taskMapById := make(map[[32]byte]storage.TaskData) // Helper to map results back
		for _, t := range tasks {
			collectedTaskIds = append(collectedTaskIds, t.TaskId)
			taskMapById[t.TaskId] = t
		}

		if len(collectedTaskIds) == 0 {
			// Should not happen if tasks slice was not empty
			return validator, validTasksForBatch, nil
		}

		// Make a single bulk call to get solution statuses with retry
		var bulkSolutions []bulktasks.IArbiusEngineEngineSolution

		getSolutionsOperation := func() (interface{}, error) {
			solutions, err := tm.services.BulkTasks.GetSolutions(nil, collectedTaskIds)
			if err != nil {
				return nil, err
			}
			return solutions, nil
		}

		retryResult, err := utils.ExpRetry(tm.services.Logger, getSolutionsOperation, 3, 1000)

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error checking on-chain solutions")
			return validator, nil, err // Return error, no fallback
		}

		// Type assert the result from ExpRetry
		var ok bool
		bulkSolutions, ok = retryResult.([]bulktasks.IArbiusEngineEngineSolution)
		if !ok {
			errInternal := errors.New("ExpRetry returned unexpected type for GetSolutions")
			tm.services.Logger.Error().Err(errInternal).Msg("internal error: type assertion failed for bulkSolutions")
			return validator, nil, errInternal
		}

		// Process bulk results
		if len(bulkSolutions) != len(collectedTaskIds) {
			tm.services.Logger.Error().Msgf("bulk solution check returned %d results for %d tasks, mismatch, skipping processing this batch", len(bulkSolutions), len(collectedTaskIds))
			return validator, nil, errors.New("bulk solution check result count mismatch")
		}

		for i, solInfo := range bulkSolutions {
			taskIdBytes := collectedTaskIds[i] // The TaskId from our input list
			originalTask, found := taskMapById[taskIdBytes]
			if !found {
				tm.services.Logger.Error().Str("taskid", task.TaskId(taskIdBytes).String()).Msg("original task not found for task ID after bulk solution call")
				continue
			}

			if solInfo.Blocktime > 0 {
				if solInfo.Claimed {
					tm.services.Logger.Info().Str("taskid", originalTask.TaskId.String()).Str("validator", solInfo.Validator.String()).Msg("task already claimed (bulk), ensuring local cleanup")
					claimsToDelete = append(claimsToDelete, originalTask.TaskId)
				} else {
					solutionsToDelete = append(solutionsToDelete, originalTask.TaskId)
					commitmentsToDelete = append(commitmentsToDelete, originalTask.TaskId)
					tm.services.Logger.Info().Str("taskid", originalTask.TaskId.String()).Str("validator", solInfo.Validator.String()).Uint64("blocktime", solInfo.Blocktime).Msg("task solution found (bulk) and not claimed, updating status to claimable")
					claimTime := time.Unix(int64(solInfo.Blocktime), 0)
					_, taskErr := tm.services.TaskStorage.UpsertTaskToClaimable(originalTask.TaskId, common.Hash{}, claimTime)
					if taskErr != nil {
						tm.services.Logger.Error().Err(taskErr).Str("taskid", originalTask.TaskId.String()).Msg("failed to update task status to claimable (bulk)")
						continue
					}
				}
				if solInfo.Validator.String() != validator.ValidatorAddress().String() {
					tm.services.Logger.Warn().Msgf("solution solved by another validator (bulk)! solver: %s task: %s", solInfo.Validator.String(), originalTask.TaskId.String())
				}
			} else {
				validTasksForBatch = append(validTasksForBatch, originalTask)
			}
		}

		if len(claimsToDelete) > 0 {
			tm.services.Logger.Info().Int("claimed", len(claimsToDelete)).Msg("deleting claimed tasks from storage")
			tm.services.TaskStorage.DeleteClaims(claimsToDelete)
		}

		if err := deleteCommitments(commitmentsToDelete); err != nil {
			return nil, nil, err
		}

		if err := deleteSolutions(solutionsToDelete); err != nil {
			return nil, nil, err
		}

		return validator, validTasksForBatch, nil
	}

	sendCommitments := func(batchTasks storage.TaskDataSlice, account *account.Account, wg *sync.WaitGroup, noChecks bool, minBatchSize int) error {
		defer wg.Done()

		batchCommitments, commitmentsToTaskMap := batchTasks.GetCommitments()
		batchCommitmentLen := len(batchCommitments)

		if batchCommitmentLen == 0 {
			return nil
		}
		if batchCommitmentLen < minBatchSize {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("commitments", batchCommitmentLen).Msg("available commitments below min batch size")
			return nil
		}

		tm.services.Logger.Info().Str("account", account.Address.String()).Msgf("preparing %d commitment(s)", batchCommitmentLen)

		isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode
		hardcodedGasLimit := calculateGasLimit(false, batchCommitmentLen, baseGasLimitForSignalCommitments, signalCommitmentGasPerItem)

		// NonceManagerWrapper call
		tx, err := account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			opts.GasLimit = calculateGasLimit(isEstimationMode, batchCommitmentLen, baseGasLimitForSignalCommitments, signalCommitmentGasPerItem)
			if isEstimationMode && tm.services.Config.Miner.GasEstimationMargin > 0 {
				opts.GasMargin = tm.services.Config.Miner.GasEstimationMargin
			}
			opts.NoSend = true // Prepare transaction but do not send automatically via binding
			txToSign, err := tm.services.BulkTasks.BulkSignalCommitment(opts, batchCommitments)
			if err != nil {
				return nil, err
			}
			return account.SendSignedTransaction(txToSign)
		})

		logCtx := map[string]interface{}{"batch_size": batchCommitmentLen, "account": account.Address.String()}
		logGasLimitDetails(tm.services.Logger, tx, isEstimationMode, hardcodedGasLimit, "sendCommitments", logCtx)

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error preparing/sending batch commitment")
			return err
		}
		if tx == nil {
			err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
			tm.services.Logger.Error().Err(err).Msg("error preparing/sending batch commitment")
			return err
		}

		receipt, success, _, waitErr := account.WaitForConfirmedTx(tx)

		// Process receipt and metrics (handles nil receipt)
		txCost := tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddCommitment, batchCommitmentLen, "batch commitments")

		if waitErr != nil {
			tm.services.Logger.Error().Err(waitErr).Str("txhash", tx.Hash().String()).Msg("Error waiting for commitment confirmation")
			return waitErr
		}

		if !success {
			// Transaction was mined but reverted
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments transaction reverted")
			return errors.New("batch commitments tx reverted")
		} else if receipt == nil {
			err := errors.New("assertion: batch commitments transaction returned nil receipt despite no error")
			tm.services.Logger.Error().Err(err).Msg("error sending batch commitments transaction")
			return err
		}

		var signalledCommitments = make(map[[32]byte]bool) // Use map for faster lookups
		tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments transaction accepted")
		for _, log := range receipt.Logs {
			if len(log.Topics) > 0 && log.Topics[0] == tm.signalCommitmentEvent {
				parsed, parseErr := tm.services.Engine.Engine.ParseSignalCommitment(*log)
				if parseErr != nil {
					tm.services.Logger.Error().Err(parseErr).Msg("could not parse signal commitment event")
					continue
				}
				signalledCommitments[parsed.Commitment] = true
			}
		}

		var commitmentsToUpdateAndDelete []task.TaskId
		var commitmentsToCheckInBulk [][32]byte
		taskIdMapForBulkCheck := make(map[[32]byte]task.TaskId)

		for _, commitment := range batchCommitments {
			if signalledCommitments[commitment] {
				if taskId, found := commitmentsToTaskMap[commitment]; found {
					commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
				}
			} else {
				// Not found in this tx's logs, add to bulk check list
				commitmentsToCheckInBulk = append(commitmentsToCheckInBulk, commitment)
				if taskId, found := commitmentsToTaskMap[commitment]; found {
					taskIdMapForBulkCheck[commitment] = taskId
				} else {
					tm.services.Logger.Warn().Str("commitment_hex", hex.EncodeToString(commitment[:])).Msg("Original taskId not found for commitment to be bulk checked")
				}
			}
		}

		if len(commitmentsToCheckInBulk) > 0 {
			tm.services.Logger.Debug().Int("count", len(commitmentsToCheckInBulk)).Msg("performing bulk check for commitments not found in transaction logs")
			getCommitmentsOperation := func() (interface{}, error) {
				return tm.services.BulkTasks.GetCommitments(nil, commitmentsToCheckInBulk)
			}
			retryResult, expRetryErr := utils.ExpRetry(tm.services.Logger, getCommitmentsOperation, 3, 1000)
			if expRetryErr != nil {
				tm.services.Logger.Error().Err(expRetryErr).Msg("bulk getcommitments call failed after retries")
			} else {
				blockNumbers, ok := retryResult.([]*big.Int)
				if !ok {
					tm.services.Logger.Error().Msg("assertion: returned unexpected type for getcommitments")
				} else if len(blockNumbers) != len(commitmentsToCheckInBulk) {
					tm.services.Logger.Error().Int("expected", len(commitmentsToCheckInBulk)).Int("got", len(blockNumbers)).Msg("assertion: mismatch in bulk getcommitments result length")
				} else {
					for i, bn := range blockNumbers {
						if bn.Cmp(utils.Zero) != 0 { // If block number > 0, it exists on-chain
							commitmentChecked := commitmentsToCheckInBulk[i]
							if taskId, found := taskIdMapForBulkCheck[commitmentChecked]; found {
								commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
								tm.services.Logger.Debug().Str("commitment_hex", hex.EncodeToString(commitmentChecked[:])).Str("taskid", taskId.String()).Uint64("block_no", bn.Uint64()).Msg("commitment found on-chain via bulk follow-up check")
							} else {
								// This should ideally not happen if taskIdMapForBulkCheck was populated correctly
								tm.services.Logger.Warn().Str("commitment_hex", hex.EncodeToString(commitmentChecked[:])).Msg("task id not found for commitment confirmed in bulk follow-up check.")
							}
						} else {
							commitmentChecked := commitmentsToCheckInBulk[i]
							tm.services.Logger.Warn().Str("commitment_hex", hex.EncodeToString(commitmentChecked[:])).Msg("commitment not found on-chain via bulk follow-up check (was not in logs and not on chain)")
						}
					}
				}
			}
		}

		// Delete local commitments now that they are confirmed on-chain (or found to be already confirmed)
		if err := deleteCommitments(commitmentsToUpdateAndDelete); err != nil {
			return err
		}
		unacceptedCommitment := len(batchCommitments) - len(commitmentsToUpdateAndDelete)
		if unacceptedCommitment > 0 {
			tm.services.Logger.Warn().Int("unaccepted", unacceptedCommitment).Int("processed", len(commitmentsToUpdateAndDelete)).Msg("⚠️ some commitments not confirmed on-chain or already present")
		}
		if len(commitmentsToUpdateAndDelete) > 0 {
			costPerCommitment := 0.0
			if len(commitmentsToUpdateAndDelete) > 0 && txCost.Cmp(big.NewInt(0)) > 0 {
				costPerCommitment = tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(commitmentsToUpdateAndDelete))
			}
			updateErr := tm.services.TaskStorage.UpdateTaskStatusAndCost(commitmentsToUpdateAndDelete, 2, costPerCommitment)
			if updateErr != nil {
				tm.services.Logger.Error().Err(updateErr).Msg("error updating task data in storage")
				return updateErr
			}
			tm.services.Logger.Info().Int("accepted", len(commitmentsToUpdateAndDelete)).Float64("cost_per_commit", costPerCommitment).Msg("✅ commitments accepted and updated in storage")
		}

		return nil // Successful completion
	}

	sendSolutions := func(validatorToSendSubmits *Validator, batchTasks storage.TaskDataSlice, wg *sync.WaitGroup, noChecks bool, minBatchSize int) error {

		batchSolutions, batchTaskIds := batchTasks.GetSolutions()
		batchSolutionsLen := len(batchSolutions)

		if batchSolutionsLen == 0 {
			return nil
		}
		if batchSolutionsLen < minBatchSize {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("solutions", batchSolutionsLen).Msg("available solutions below min batch size")
			return nil
		}

		var solutionsToSubmit [][]byte
		var tasksToSubmit [][32]byte
		for i, cid := range batchSolutions {
			taskid := batchTaskIds[i]
			solutionsToSubmit = append(solutionsToSubmit, cid)
			tasksToSubmit = append(tasksToSubmit, taskid)
		}

		if len(solutionsToSubmit) != len(tasksToSubmit) {
			err := errors.New("assertion: mismatched number of tasks and solutions being submitted")
			tm.services.Logger.Error().Err(err).Msg("internal error submitting solutions")
			return err
		}

		tm.services.Logger.Info().Str("account", validatorToSendSubmits.ValidatorAddress().String()).Msgf("preparing %d solution(s)", len(solutionsToSubmit))

		isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode
		hardcodedGasLimit := calculateGasLimit(false, batchSolutionsLen, baseGasLimitForSubmitSolutions, submitSolutionGasPerItem)

		// NonceManagerWrapper call
		tx, err := validatorToSendSubmits.Account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			opts.GasLimit = calculateGasLimit(isEstimationMode, batchSolutionsLen, baseGasLimitForSubmitSolutions, submitSolutionGasPerItem)
			if isEstimationMode && tm.services.Config.Miner.GasEstimationMargin > 0 {
				opts.GasMargin = tm.services.Config.Miner.GasEstimationMargin
			}
			opts.NoSend = true
			txToSign, err := tm.services.Engine.Engine.BulkSubmitSolution(opts, tasksToSubmit, solutionsToSubmit)
			if err != nil {
				return nil, err
			}
			return validatorToSendSubmits.Account.SendSignedTransaction(txToSign)
		})

		logCtx := map[string]interface{}{"batch_size": batchSolutionsLen, "account": validatorToSendSubmits.ValidatorAddress().String()}
		logGasLimitDetails(tm.services.Logger, tx, isEstimationMode, hardcodedGasLimit, "sendSolutions", logCtx)

		if err != nil {
			tm.services.Logger.Error().Err(err).Int("batch_size", batchSolutionsLen).Msg("error sending batch solutions transaction")
			return err
		}
		if tx == nil {
			err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
			tm.services.Logger.Error().Err(err).Msg("error sending batch solutions transaction")
			return err
		}

		receipt, success, _, waitErr := validatorToSendSubmits.Account.WaitForConfirmedTx(tx)

		// Process receipt and metrics
		txCost := tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddSolution, batchSolutionsLen, "batch solutions")

		if waitErr != nil {
			return waitErr
		}

		if !success {
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Msg("batch solution transaction reverted on-chain")
			return errors.New("batch solution transaction reverted on-chain")
		} else if receipt == nil {
			err := errors.New("assertion: batch solution transaction returned nil receipt despite no error")
			tm.services.Logger.Error().Err(err).Msg("error sending batch solutions transaction")
			return err
		}

		tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch solutions transaction accepted")
		var solutionsSubmittedInLogs = make(map[[32]byte]bool) // Use map for faster lookups
		for _, log := range receipt.Logs {
			if len(log.Topics) > 0 && log.Topics[0] == tm.solutionSubmittedEvent {
				parsed, parseErr := tm.services.Engine.Engine.ParseSolutionSubmitted(*log)
				if parseErr != nil {
					tm.services.Logger.Error().Err(parseErr).Msg("could not parse solution submitted event")
					continue
				}
				solutionsSubmittedInLogs[parsed.Task] = true
			}
		}

		solutionsToDelete := make([]task.TaskId, 0)
		var tasksToCheckInBulk [][32]byte

		for _, taskid := range tasksToSubmit { // tasksToSubmit is [][32]byte
			if solutionsSubmittedInLogs[taskid] {
				tm.services.TaskTracker.TaskSucceeded()
				solutionsToDelete = append(solutionsToDelete, taskid)
			} else {
				tm.services.TaskTracker.TaskFailed() // Assume failed for now, will be checked in bulk
				tasksToCheckInBulk = append(tasksToCheckInBulk, taskid)
			}
		}

		if len(tasksToCheckInBulk) > 0 {
			tm.services.Logger.Debug().Int("count", len(tasksToCheckInBulk)).Msg("performing bulk check for solutions not found in transaction logs")
			getSolutionsOperation := func() (interface{}, error) {
				return tm.services.BulkTasks.GetSolutions(nil, tasksToCheckInBulk)
			}
			retryResult, expRetryErr := utils.ExpRetry(tm.services.Logger, getSolutionsOperation, 3, 1000)

			if expRetryErr != nil {
				tm.services.Logger.Error().Err(expRetryErr).Msg("bulk getsolutions call failed after retries for follow-up check")
				// Continue with solutions found in logs, others remain unconfirmed for this round
			} else {
				bulkSolutions, ok := retryResult.([]bulktasks.IArbiusEngineEngineSolution)
				if !ok {
					tm.services.Logger.Error().Msg("assertion: returned unexpected type for getsolutions")
				} else if len(bulkSolutions) != len(tasksToCheckInBulk) {
					tm.services.Logger.Error().Int("expected", len(tasksToCheckInBulk)).Int("got", len(bulkSolutions)).Msg("assertion: mismatch in bulk getsolutions result length")
				} else {
					for i, solInfo := range bulkSolutions {
						taskIdChecked := tasksToCheckInBulk[i]
						taskIdStr := task.TaskId(taskIdChecked).String()
						if solInfo.Blocktime > 0 { // If blocktime > 0, solution exists on-chain
							solutionsToDelete = append(solutionsToDelete, taskIdChecked)
							tm.services.Logger.Info().Str("taskid", taskIdStr).Str("validator", solInfo.Validator.String()).Uint64("blocktime", solInfo.Blocktime).Msg("solution found on-chain via bulk follow-up check")
							// If it was solved by us, TaskSucceeded was already called.
							// If solved by someone else, TaskFailed was already called. This is just confirming it's on chain.
						} else {
							tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("Solution NOT found on-chain via bulk follow-up check (was not in logs and not on chain)")
							// TaskFailed was already called for these.
						}
					}
				}
			}
		}

		// do not handle error here as we want to continue processing the next batch
		deleteErr := deleteSolutions(solutionsToDelete)
		if deleteErr != nil {
			tm.services.Logger.Error().Err(deleteErr).Msg("error deleting solutions from storage")
		}

		// Reconstruct the list of task IDs that were successfully submitted (found in logs)
		// to be added to the claim queue.
		var tasksActuallySubmittedToClaim []task.TaskId
		for taskidBytes := range solutionsSubmittedInLogs {
			tasksActuallySubmittedToClaim = append(tasksActuallySubmittedToClaim, task.TaskId(taskidBytes))
		}

		unacceptedSolutions := len(tasksToSubmit) - len(solutionsSubmittedInLogs)
		if unacceptedSolutions > 0 {
			tm.services.Logger.Warn().Int("unaccepted", unacceptedSolutions).Int("total_in_batch", len(tasksToSubmit)).Msg("⚠️ some solutions submitted in batch were not confirmed via logs or subsequent bulk check")
		}

		if len(tasksActuallySubmittedToClaim) > 0 { // Use the reconstructed slice here
			gasPerSolution := 0.0
			if len(tasksActuallySubmittedToClaim) > 0 && txCost.Cmp(big.NewInt(0)) > 0 {
				gasPerSolution = tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(tasksActuallySubmittedToClaim))
			}
			_, addClaimErr := tm.services.TaskStorage.AddTasksToClaim(tasksActuallySubmittedToClaim, gasPerSolution) // Pass the correct slice
			if addClaimErr != nil {
				tm.services.Logger.Error().Err(addClaimErr).Msg("error adding tasks to claim in storage")
				return addClaimErr
			}
			tm.services.Logger.Info().Str("validator", validatorToSendSubmits.ValidatorAddress().String()).Int("accepted_in_log", len(tasksActuallySubmittedToClaim)).Float64("cost_per_sol", gasPerSolution).Msg("✅ submitted solutions (confirmed in logs) added to claim queue in storage")
		}

		return nil // Successful completion
	}

	if tm.services.Config.Miner.CommitmentsAndSolutions != config.DoNothing {
		//  get radom index into taskAccounts
		accountIndex := rand.Intn(len(tm.taskAccounts))

		//minBatchSize := tm.services.Config.Miner.MinBatchSize//
		//maxBatchSize := tm.services.Config.Miner.MaxBatchSize

		doCommitments := tm.services.Config.Miner.CommitmentsAndSolutions == config.DoCommitmentsOnly || tm.services.Config.Miner.CommitmentsAndSolutions == config.DoBoth
		doSolutions := tm.services.Config.Miner.CommitmentsAndSolutions == config.DoBoth || tm.services.Config.Miner.CommitmentsAndSolutions == config.DoSolutionsOnly
		if doCommitments {
			minBatchSize := tm.services.Config.Miner.CommitmentBatch.MinBatchSize
			maxBatchSize := tm.services.Config.Miner.CommitmentBatch.MaxBatchSize
			allBatchCommitments, err := getCommitmentBatchdata(tm.services.Config.Miner.CommitmentBatch.NumberOfBatches*maxBatchSize, tm.services.Config.Miner.NoChecks)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error getting commitment batch data")
				return
			}

			for i := 0; i < tm.services.Config.Miner.CommitmentBatch.NumberOfBatches; i++ {
				if err := appQuit.Err(); err != nil {
					return
				}

				// get the next available account
				account := tm.taskAccounts[accountIndex%len(tm.taskAccounts)]
				accountIndex++

				startIndex := i * maxBatchSize

				if startIndex >= len(allBatchCommitments) {
					break
				}

				endIndex := min(len(allBatchCommitments), (i+1)*maxBatchSize)
				batchTasks := allBatchCommitments[startIndex:endIndex]

				wg.Add(1)
				if tm.services.Config.Miner.ConcurrentBatches {
					go sendCommitments(batchTasks, account, wg, tm.services.Config.Miner.NoChecks, minBatchSize)
				} else {
					sendCommitments(batchTasks, account, wg, tm.services.Config.Miner.NoChecks, minBatchSize)
				}
			}

		}

		if doSolutions {
			minBatchSize := tm.services.Config.Miner.SolutionBatch.MinBatchSize
			maxBatchSize := tm.services.Config.Miner.SolutionBatch.MaxBatchSize

			// Revert to original: Use configured NoChecks value
			validator, batchTasks, err := getSolutionBatchdata(tm.services.Config.Miner.SolutionBatch.NumberOfBatches*maxBatchSize, tm.services.Config.Miner.NoChecks)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error getting solution batch data")
				return
			}

			// if we have no tasks or validator, we can return (not an error)
			if batchTasks == nil || validator == nil {
				return
			}
			sendSolutions(validator, batchTasks, wg, tm.services.Config.Miner.NoChecks, minBatchSize)
		}
	}
}

// TODO: Remove this function
func (tm *BatchTransactionManager) processBulkClaimFast(account *account.Account, tasks []storage.ClaimTask, batchSize int) {
	panic("processBulkClaimFast is deprecated")
	/*
		// Use the Map function to extract the TaskId from each struct
		taskIds := utils.Map(tasks, func(s storage.ClaimTask) [32]byte {
			//value, _ := task.ConvertTaskIdString2Bytes(s.ID)
			return s.ID
		})

		receipt, err := tm.BulkClaimWithAccount(account, taskIds)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("❌ error submitting bulk claim")
			return //err
		}

		// Check if receipt is nil before accessing Status
		if receipt == nil {
			tm.services.Logger.Error().Msg("❌ bulk claim returned nil receipt despite no error")
			return
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			tm.services.Logger.Warn().Str("txhash", receipt.TxHash.String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("⚠️ bulk claim transaction failed/reverted")
			return
		}

		tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Int("tasks", len(taskIds)).Uint64("block", receipt.BlockNumber.Uint64()).Msg("bulk claim transaction accepted")

		tasksClaimed := make([]task.TaskId, 0)

		var totalValidatorReward *big.Int = big.NewInt(0)

		for _, log := range receipt.Logs {
			if len(log.Topics) > 0 {
				// Check for SolutionClaimed event
				if log.Topics[0] == tm.solutionClaimedEvent {
					parsed, err := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
					if err != nil {
						tm.services.Logger.Error().Err(err).Msg("could not parse solution claimed event")
						continue
					}
					tasksClaimed = append(tasksClaimed, task.TaskId(parsed.Task))
				}

				// Check for RewardsPaid event to sum validatorReward
				if log.Topics[0] == tm.rewardsPaidEvent {
					parsed, err := tm.services.Engine.Engine.ParseRewardsPaid(*log)
					if err != nil {
						tm.services.Logger.Error().Err(err).Msg("could not parse rewards paid event")
						continue
					}
					totalValidatorReward.Add(totalValidatorReward, parsed.ValidatorReward)
				}
			}
		}

		if len(tasksClaimed) == 0 {
			tm.services.Logger.Warn().Msg("⚠️ successful bulk claim transaction but no tasks were claimed: check for contestations/cooldown issues ⚠️")
		}

		if len(tasksClaimed) > 0 {
			tm.services.Logger.Info().Int("claimed", len(tasksClaimed)).Msg("✅ successfully claimed tasks and removed from storage")
		}

		unclaimedTaskCount := len(tasks) - len(tasksClaimed)
		if unclaimedTaskCount > 0 {
			tm.services.Logger.Warn().Int("unclaimed", unclaimedTaskCount).Msg("⚠️ tasks that failed to be claimed")
		}*/
}

// returns true if the task is valid and can be claimed or false if not
func (tm *BatchTransactionManager) canTaskIdBeClaimed(
	claim storage.ClaimTask,
	combinedInfo *bulktasks.BulkTasksTaskSolutionWithContestationInfo,
	cooldownTimes map[common.Address]uint64,
) bool {
	taskIdStr := claim.ID.String()

	// First check if solution exists
	if !combinedInfo.SolutionExists {
		tm.services.Logger.Warn().Str("task", taskIdStr).Msg("no valid solution details in combined info or solution not found on-chain")
		return false
	}

	// Then check if already claimed
	if combinedInfo.SolutionClaimed {
		tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("solution already claimed")
		return false
	}

	// Then check contestation
	if combinedInfo.ContestationExists {
		tm.services.Logger.Warn().Str("task", taskIdStr).Str("contestor", combinedInfo.ContestationValidator.String()).Msg("⚠️ task was contested (details from combined bulk) ⚠️")
		return false
	}

	solTime := time.Unix(int64(combinedInfo.SolutionBlocktime), 0)
	cooldownTime := cooldownTimes[combinedInfo.SolutionValidator]
	if combinedInfo.SolutionBlocktime <= cooldownTime {
		tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("⚠️ claim is lost due to lost contestation cooldown - removing from storage ⚠️")
		return false
	}

	tm.services.Logger.Debug().Str("taskid", taskIdStr).Time("solved", solTime).Str("validator", combinedInfo.SolutionValidator.String()).Msg("solution information (using combined pre-fetched data)")

	return true
}
func (tm *BatchTransactionManager) processBulkClaim(account *account.Account, tasks []storage.ClaimTask, minbatchSize, maxbatchSize int) {
	if len(tasks) == 0 {
		tm.services.Logger.Error().Msg("assertion: process bulk claim called with no tasks to claim")
		return
	}

	// Get the CooldownTime for each validator we might be claiming for
	var validatorCooldownTimesMap = make(map[common.Address]uint64)
	for _, v := range tm.services.Validators.validators {
		cooldownTime, err := v.CooldownTime(tm.minClaimSolutionTime, tm.minContestationVotePeriodTime)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error calling MinContestationVotePeriodTime")
			return
		}
		validatorCooldownTimesMap[v.ValidatorAddress()] = cooldownTime
	}

	tasksToClaim := make([]storage.ClaimTask, 0)
	claimsToDelete := make([]task.TaskId, 0)

	// Fetch combined solution and contestation data upfront
	bulkCombinedInfoMap := make(map[task.TaskId]bulktasks.BulkTasksTaskSolutionWithContestationInfo)
	taskIdsForBulkFetch := make([][32]byte, 0, len(tasks))

	// tasks is the input []storage.ClaimTask
	for _, taskToValidate := range tasks {
		taskIdsForBulkFetch = append(taskIdsForBulkFetch, taskToValidate.ID)
	}

	tm.services.Logger.Debug().Int("count", len(taskIdsForBulkFetch)).Msg("bulk fetching combined task info for claim processing")
	getCombinedInfoOperation := func() (interface{}, error) {
		return tm.services.BulkTasks.GetBulkCombinedTaskInfo(nil, taskIdsForBulkFetch)
	}
	retryResult, expRetryErr := utils.ExpRetry(tm.services.Logger, getCombinedInfoOperation, 3, 1000)

	if expRetryErr != nil {
		tm.services.Logger.Error().Err(expRetryErr).Msg("could not bulk fetch combined task info for claim processing")
	} else {
		combinedInfos, ok := retryResult.([]bulktasks.BulkTasksTaskSolutionWithContestationInfo)
		if !ok {
			tm.services.Logger.Error().Msg("assertion: expRetry returned unexpected type for getbulkcombinedtaskinfo")
		} else if len(combinedInfos) != len(taskIdsForBulkFetch) {
			tm.services.Logger.Error().Int("expected", len(taskIdsForBulkFetch)).Int("got", len(combinedInfos)).Msg("assertion: mismatch in bulk getcombinedtaskinfo result length")
		} else {
			for _, info := range combinedInfos { // Iterate through the results and map them by their TaskId field
				if info.TaskId != ([32]byte{}) {
					bulkCombinedInfoMap[info.TaskId] = info
				} else {
					// This case should ideally not happen if the contract always returns TaskId
					tm.services.Logger.Warn().Msg("assertion: received a combinedtaskinfo struct with zero taskid from bulk call")
				}
			}
		}
	}

	// Pre-claim validation loop
	for _, taskItem := range tasks { // taskItem is storage.ClaimTask
		taskIdStr := task.TaskId(taskItem.ID).String()

		if taskItem.Time >= time.Now().Unix() {
			continue
		}

		retrievedCombinedInfo, infoFetched := bulkCombinedInfoMap[taskItem.ID]
		if !infoFetched {
			tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("no combined info available for task, skipping claim validation")
			continue
		}

		claimable := tm.canTaskIdBeClaimed(taskItem, &retrievedCombinedInfo, validatorCooldownTimesMap)
		if !claimable {
			claimsToDelete = append(claimsToDelete, taskItem.ID)
			tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("task determined not claimable (pre-tx check), adding to delete queue")
			continue
		}
		tasksToClaim = append(tasksToClaim, taskItem)
	}

	if len(claimsToDelete) > 0 {
		err := tm.services.TaskStorage.DeleteClaims(claimsToDelete)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error deleting claims pre-transaction")
			// Decide if to return or continue
		}
	}

	if len(tasksToClaim) < minbatchSize {
		tm.services.Logger.Info().Int("min_batch", minbatchSize).Int("max_batch", maxbatchSize).Int("claimable", len(tasksToClaim)).Msg("claimable less than min batch size")
		return
	}

	sort.Slice(tasksToClaim, func(i, j int) bool {
		return tasksToClaim[i].Time < tasksToClaim[j].Time
	})

	if len(tasksToClaim) > maxbatchSize {
		tasksToClaim = tasksToClaim[:maxbatchSize]
	}

	taskIdsForTx := utils.Map(tasksToClaim, func(s storage.ClaimTask) [32]byte {
		return s.ID
	})

	if len(taskIdsForTx) == 0 { // Add check here to avoid sending empty tx if all were filtered out
		tm.services.Logger.Info().Msg("No tasks to claim after validation and filtering.")
		return
	}

	receipt, err := tm.BulkClaimWithAccount(account, taskIdsForTx)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("❌ error submitting bulk claim")
		return
	}
	if receipt == nil {
		tm.services.Logger.Error().Msg("❌ bulk claim returned nil receipt")
		return
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		tm.services.Logger.Warn().Str("txhash", receipt.TxHash.String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("⚠️ bulk claim transaction failed or reverted")
		return
	}

	tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Int("tasks_sent", len(taskIdsForTx)).Uint64("block", receipt.BlockNumber.Uint64()).Msg("bulk claim transaction accepted")

	tasksClaimedInLogs := make(map[task.TaskId]bool)
	var totalValidatorReward *big.Int = big.NewInt(0) // Keep this if needed for logging rewards

	for _, log := range receipt.Logs {
		if len(log.Topics) > 0 {
			if log.Topics[0] == tm.solutionClaimedEvent {
				parsed, parseErr := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
				if parseErr != nil {
					tm.services.Logger.Error().Err(parseErr).Msg("could not parse solution claimed event")
					continue
				}
				tasksClaimedInLogs[task.TaskId(parsed.Task)] = true
			}
			if log.Topics[0] == tm.rewardsPaidEvent { // Keep if reward logging is desired
				parsedRewards, parseErrRewards := tm.services.Engine.Engine.ParseRewardsPaid(*log)
				if parseErrRewards != nil {
					tm.services.Logger.Error().Err(parseErrRewards).Msg("could not parse rewards paid event")
					continue
				}
				totalValidatorReward.Add(totalValidatorReward, parsedRewards.ValidatorReward)
			}
		}
	}

	if len(tasksClaimedInLogs) == 0 {
		tm.services.Logger.Warn().Msg("⚠️ successful bulk claim transaction but no tasks were claimed: check for contestations/cooldown issues ⚠️")
	}

	if totalValidatorReward.Cmp(big.NewInt(0)) > 0 { // Log total rewards from this batch
		totalValidatorRewardInAius := tm.services.Config.BaseConfig.BaseToken.ToFloat(totalValidatorReward)
		tm.services.Logger.Info().Float64("validator_reward_batch", totalValidatorRewardInAius).Msg("total validator reward for this bulk claim tx")
	}

	// Post-transaction check for tasks that were sent but not found in logs
	claimsToDeletePostTx := make([]task.TaskId, 0)
	for _, taskidItem := range tasksToClaim { // taskidItem is storage.ClaimTask
		if !tasksClaimedInLogs[taskidItem.ID] {
			taskIdStr := task.TaskId(taskidItem.ID).String()
			retrievedCombinedInfo, infoFetched := bulkCombinedInfoMap[taskidItem.ID]
			if !infoFetched {
				tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("no combined info available for task, skipping post-tx claim validation")
				continue
			}
			claimable := tm.canTaskIdBeClaimed(taskidItem, &retrievedCombinedInfo, validatorCooldownTimesMap)
			if !claimable {
				tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("task determined not claimable (post-tx check), adding to delete queue")
				claimsToDeletePostTx = append(claimsToDeletePostTx, taskidItem.ID)
			}
		}

	}

	// Delete tasks that were successfully claimed or determined unclaimable
	var finalClaimsToDelete []task.TaskId

	// Add tasks that were successfully claimed according to logs
	for taskId := range tasksClaimedInLogs {
		finalClaimsToDelete = append(finalClaimsToDelete, taskId)
	}

	// Add tasks that were determined unclaimable in post-tx check
	finalClaimsToDelete = append(finalClaimsToDelete, claimsToDeletePostTx...)

	if len(finalClaimsToDelete) > 0 {
		tm.services.Logger.Info().Int("count", len(finalClaimsToDelete)).Msg("deleting claimed/unclaimable tasks from storage")
		err := tm.services.TaskStorage.DeleteClaims(finalClaimsToDelete)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error deleting final claims list")
		}
	}

	unclaimedCount := len(taskIdsForTx) - len(tasksClaimedInLogs) // Compare against tasks sent in TX
	if unclaimedCount > 0 {
		tm.services.Logger.Warn().Int("unclaimed", unclaimedCount).Int("sent_in_batch", len(taskIdsForTx)).Msg("⚠️ some tasks sent in batch that were not confirmed claimed via logs")
	} else if len(tasksClaimedInLogs) > 0 {
		tm.services.Logger.Info().Int("claimed_in_log", len(tasksClaimedInLogs)).Msg("✅ successfully processed bulk claim transaction results")
	}
}

func (tm *BatchTransactionManager) BatchCommitments() error {
	tm.Lock()
	copyCommitments := make([][32]byte, len(tm.commitments))
	copy(copyCommitments, tm.commitments)
	tm.commitments = [][32]byte{}
	tm.Unlock()

	tm.services.Logger.Info().Msgf("bulk submitting %d commitment(s)", len(copyCommitments))

	tx, err := tm.services.OwnerAccount.NonceManagerWrapper(5, 425, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = uint64(27_900*len(copyCommitments) + 1_500_000)
		return tm.services.BulkTasks.BulkSignalCommitment(opts, copyCommitments)
	})

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error sending batch commitments")
		return err
	}

	receipt, success, _, _ := tm.services.OwnerAccount.WaitForConfirmedTx(tx)

	if receipt != nil {
		txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
		tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Float64("gas_per_commit", float64(receipt.GasUsed)/float64(len(copyCommitments))).Msg("**** bulk Commitment gas used *****")
		tm.cumulativeGasUsed.AddCommitment(txCost)
	}

	if !success {
		tm.services.Logger.Error().Err(err).Msg("batch commitments tx reverted")
		return err
	}

	tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments tx accepted!")

	return nil
}

func (tm *BatchTransactionManager) SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error {
	if tm.services.Config.IPFS.IncentiveClaim {
		return tm.services.TaskStorage.AddIpfsCid(taskId, cid)
	}
	return nil
}

func (tm *BatchTransactionManager) SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error {
	tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task commitment to storage")
	added, err := tm.services.TaskStorage.TryAddCommitment(validator, taskId, commitment)
	if err != nil {
		return err
	}
	if !added {
		tm.services.Logger.Warn().Str("taskid", taskId.String()).Msg("commitment already exists, skipping add")
		return nil
	}
	return nil
}

func (tm *BatchTransactionManager) SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error {
	tm.services.TaskTracker.Solved()
	tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task solution to batch")
	err := tm.services.TaskStorage.AddSolution(validator, taskId, cid)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error adding solution to storage")
	}
	return err
}

func (tm *BatchTransactionManager) SignalCommitmentNow(validator common.Address, taskId task.TaskId, commitment [32]byte) error {

	if tm.services.Config.CheckCommitment {
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("checking for existing task commitment")

		block, err := tm.services.Engine.Engine.Commitments(nil, commitment)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error getting commitment")
			return err
		}

		blockNo := block.Uint64()
		if blockNo > 0 {
			tm.services.Logger.Warn().Str("taskid", taskId.String()).Uint64("block", blockNo).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("commitment already exists for task onchain")
			return nil
		}
	}
	taskIdStr := taskId.String()

	tm.services.Logger.Info().Str("taskid", taskIdStr).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("sending commitment")

	start := time.Now()
	retries := tm.services.Config.Miner.ErrorMaxRetries
	backoff := tm.services.Config.Miner.ErrorBackoffTime
	backoffMultiplier := tm.services.Config.Miner.ErrorBackofMultiplier
	tx, err := tm.services.OwnerAccount.NonceManagerWrapper(retries, backoff, backoffMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		if tm.services.Config.Miner.EnableGasEstimationMode {
			opts.GasLimit = 0
		} else {
			opts.GasLimit = signalCommitmentGasPerItem * 10 // why 10? ask the arbitrum team, floating gas intrinsics is a bitch
		}
		return tm.services.Engine.Engine.SignalCommitment(opts, commitment)
	})
	elapsed := time.Since(start)

	if err != nil {
		tm.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Str("elapsed", elapsed.String()).Msg("❌ error preparing commitment tx")
		return err
	}
	if tx == nil {
		err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		tm.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Msg("❌ error preparing commitment tx")
		return err
	}

	tm.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("signal commitment tx sent, waiting for confirmation...")

	go func() {
		receipt, success, _, waitErr := tm.services.OwnerAccount.WaitForConfirmedTx(tx)

		// Process receipt and metrics regardless of success/failure, but handle nil receipt
		tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddCommitment, 1, "single commitment") // batchSize is 1

		if waitErr != nil {
			return
		}

		if !success {
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Msg("single commitment transaction reverted on-chain")
			return // errors.New("single commitment transaction reverted on-chain")
		}
		// Double-check receipt is not nil after success
		if receipt == nil {
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Msg("assertion: single commitment transaction returned nil receipt despite success")
			return
		}

		var signalledCommitments [][32]byte
		for _, log := range receipt.Logs {

			if len(log.Topics) > 0 && log.Topics[0] == tm.signalCommitmentEvent {
				parsed, err := tm.services.Engine.Engine.ParseSignalCommitment(*log)
				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("could not parse signal commitment event")
					continue
				}
				signalledCommitments = append(signalledCommitments, parsed.Commitment)
			}
		}

		if len(signalledCommitments) != 1 {
			tm.services.Logger.Error().Msg("ASSERT: COMMITMENTS SIGNALLLED NOT EQUAL TO 1!")
		}
	}()

	// wait a bit to hope commitment is mined before submitting solution
	duration := 2000 + rand.Intn(350)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	return nil
}

func (tm *BatchTransactionManager) SubmitSolutionNow(validator common.Address, taskId task.TaskId, cid []byte) error {

	val := tm.services.Validators.GetValidatorByAddress(validator)
	if val == nil {
		return errors.New("validator not found")
	}

	taskIdStr := taskId.String()
	tm.services.Logger.Info().Str("taskid", taskIdStr).Msg("sending single solution")

	start := time.Now()
	retries := tm.services.Config.Miner.ErrorMaxRetries
	backoff := tm.services.Config.Miner.ErrorBackoffTime
	backoffMultiplier := tm.services.Config.Miner.ErrorBackofMultiplier
	tx, err := val.Account.NonceManagerWrapper(retries, backoff, backoffMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		if tm.services.Config.Miner.EnableGasEstimationMode {
			opts.GasLimit = 0
		} else {
			opts.GasLimit = submitSolutionGasPerItem * 10 // why 10? ask the arbitrum team, floating gas intrinsics is a bitch
		}

		return tm.services.Engine.Engine.SubmitSolution(opts, taskId, cid)
	})
	elapsed := time.Since(start)

	if err != nil {
		tm.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Str("elapsed", elapsed.String()).Msg("❌ error preparing solution tx")
		return err
	}
	if tx == nil {
		err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		tm.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Msg("❌ error preparing solution tx")
		return err
	}

	tm.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("solution tx sent, waiting for confirmation...")

	go func() {
		// find out who mined the soluition and log it
		defer func() {
			res, err := tm.services.Engine.GetSolution(taskId)
			if err != nil {
				tm.services.Logger.Err(err).Msg("error getting solution information")
				return
			}
			if res.Blocktime > 0 {
				if tm.services.Validators.IsAddressValidator(res.Validator) {
					tm.services.TaskTracker.TaskSucceeded()
				} else {
					tm.services.TaskTracker.TaskFailed()
				}
				solversCid := common.Bytes2Hex(res.Cid[:])
				ourCid := common.Bytes2Hex(cid)
				if ourCid != solversCid {
					tm.services.Logger.Warn().Msg("=======================================================================")
					tm.services.Logger.Warn().Msg("  WARNING: our solution cid does not match the solvers cid!")
					tm.services.Logger.Warn().Msg("  our cid: " + ourCid)
					tm.services.Logger.Warn().Msg("  ther cid: " + solversCid)
					tm.services.Logger.Warn().Str("validator", res.Validator.String()).Msg("  solvers address")
					tm.services.Logger.Warn().Msg("========================================================================")
				}
				tm.services.Logger.Info().Str("taskid", taskIdStr).Str("validator", res.Validator.String()).Str("Cid", solversCid).Msg("solution information")
			} else {
				tm.services.TaskTracker.TaskFailed()
				tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("solution not solved")
			}
		}()

		receipt, success, _, waitErr := tm.services.OwnerAccount.WaitForConfirmedTx(tx)

		// Process receipt and metrics regardless of success/failure, but handle nil receipt
		tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddSolution, 1, "single solution") // batchSize is 1

		if waitErr != nil {
			return
		}

		if !success {
			tm.services.Logger.Error().Str("txHash", tx.Hash().String()).Msg("single solution transaction reverted on-chain")
			return // errors.New("single solution transaction reverted on-chain")
		}
		// Double-check receipt is not nil after success
		if receipt == nil {
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Msg("assertion: single solution transaction returned nil receipt despite success")
			return
		}

		tm.services.Logger.Info().Str("taskid", taskIdStr).Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("✅ solution accepted!")

		// ensure the task is added to the claimable list
		claimTime, err := tm.services.TaskStorage.UpsertTaskToClaimable(taskId, tx.Hash(), time.Now())
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error adding claim in redis")
			return
		}
		tm.services.Logger.Info().Str("taskid", taskIdStr).Str("when", claimTime.String()).Msg("added taskid claim to storage")
	}()

	return nil
}

func (tm *BatchTransactionManager) BulkClaim(taskIds [][32]byte) (*types.Receipt, error) {
	return tm.BulkClaimWithAccount(tm.services.OwnerAccount, taskIds)
}

func (tm *BatchTransactionManager) BulkClaimWithAccount(account *account.Account, taskIds [][32]byte) (*types.Receipt, error) {

	taskCount := len(taskIds)
	if taskCount == 0 {
		return nil, nil
	}

	tm.services.Logger.Info().Str("account", account.Address.String()).Msgf("preparing %d claim(s)", taskCount)

	isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode
	hardcodedGasLimit := calculateGasLimit(false, taskCount, baseGasLimitForClaimTasks, claimTasksGasPerItem)

	tx, err := account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = calculateGasLimit(isEstimationMode, taskCount, baseGasLimitForClaimTasks, claimTasksGasPerItem)
		if isEstimationMode && tm.services.Config.Miner.GasEstimationMargin > 0 {
			opts.GasMargin = tm.services.Config.Miner.GasEstimationMargin
		}
		opts.NoSend = true // Prepare transaction but do not send automatically via wrapper
		txToSign, err := tm.services.BulkTasks.ClaimSolutions(opts, taskIds)
		if err != nil {
			return nil, err
		}
		return account.SendSignedTransaction(txToSign)
	})

	logCtx := map[string]interface{}{"batch_size": taskCount, "account": account.Address.String()}
	logGasLimitDetails(tm.services.Logger, tx, isEstimationMode, hardcodedGasLimit, "BulkClaimWithAccount", logCtx)

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error preparing/sending bulk claim")
		return nil, err
	}
	if tx == nil {
		err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		tm.services.Logger.Error().Err(err).Msg("error preparing/sending bulk claim")
		return nil, err
	}

	tm.cumulativeGasUsed.ClaimValue(taskCount)

	receipt, _, _, waitErr := account.WaitForConfirmedTx(tx)

	tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddClaim, taskCount, "bulk claims")

	return receipt, waitErr
}

func (tm *BatchTransactionManager) BulkTasks(account *account.Account, count int) (*types.Receipt, error) {

	if count <= 0 {
		// No tasks to submit.
		return nil, nil
	}

	// Inner function to handle sending attempt for a single client
	sendTx := func(ctx context.Context, client *client.Client, outerOpts *bind.TransactOpts) (*types.Receipt, error) {

		isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode
		hardcodedGasLimit := calculateGasLimit(false, count, baseGasLimitForSubmitTasks, submitTasksGasPerItem)

		tm.services.Logger.Info().Str("account", account.Address.String()).Msgf("preparing %d task(s)", count)

		// NonceManagerWrapper call
		tx, err := account.NonceManagerWrapperWithContext(ctx, outerOpts, tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			opts.GasLimit = calculateGasLimit(isEstimationMode, count, baseGasLimitForSubmitTasks, submitTasksGasPerItem)
			if isEstimationMode && tm.services.Config.Miner.GasEstimationMargin > 0 {
				opts.GasMargin = tm.services.Config.Miner.GasEstimationMargin
			}
			opts.NoSend = true // Prepare transaction but do not send automatically via wrapper
			tasks := big.NewInt(int64(count))
			txToSign, err := tm.services.Engine.Engine.BulkSubmitTask(opts, tm.services.AutoMineParams.Version, tm.services.AutoMineParams.Owner, tm.services.AutoMineParams.Model, tm.services.AutoMineParams.Fee, tm.services.AutoMineParams.Input, tasks)
			if err != nil {
				return nil, err
			}
			return account.SendSignedTransaction(txToSign)
		})

		logCtx := map[string]interface{}{"batch_size": count, "account": account.Address.String()}
		logGasLimitDetails(tm.services.Logger, tx, isEstimationMode, hardcodedGasLimit, "BulkTasks", logCtx)

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error preparing/sending bulk tasks")
			return nil, err
		}
		if tx == nil {
			err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
			tm.services.Logger.Error().Err(err).Msg("error preparing/sending bulk tasks")
			return nil, err
		}

		receipt, success, _, waitErr := account.WaitForConfirmedTx(tx)

		// Process receipt and metrics
		txCost := tm.processReceiptAndMetrics(receipt, tm.cumulativeGasUsed.AddTasks, count, "bulk tasks")

		if waitErr != nil {
			// Error logged in WaitForConfirmedTx
			return receipt, waitErr
		}
		if !success {
			return receipt, nil
		}

		// Process logs and update storage
		if receipt != nil {
			tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("bulk tasks transaction succeeded")

			var submittedTasks []task.TaskId
			for _, log := range receipt.Logs {
				if len(log.Topics) > 0 && log.Topics[0] == tm.taskSubmittedEvent {
					parsed, parseErr := tm.services.Engine.Engine.ParseTaskSubmitted(*log)
					if parseErr != nil {
						tm.services.Logger.Error().Err(parseErr).Msg("could not parse task submitted event")
						continue
					}
					submittedTasks = append(submittedTasks, task.TaskId(parsed.Id))
				}
			}
			missingTasks := count - len(submittedTasks)
			if missingTasks > 0 {
				tm.services.Logger.Warn().Int("missing", missingTasks).Int("expected", count).Int("found", len(submittedTasks)).Msg("⚠️ tasks created mismatch ⚠️")
			}

			if len(submittedTasks) > 0 {
				costPerTask := 0.0
				if len(submittedTasks) > 0 && txCost.Cmp(big.NewInt(0)) > 0 {
					// Use the txCost returned by processReceiptAndMetrics
					costPerTask = tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(submittedTasks))
				}

				addErr := tm.services.TaskStorage.AddTasks(submittedTasks, tx.Hash(), costPerTask)
				if addErr != nil {
					tm.services.Logger.Error().Err(addErr).Msg("error adding tasks in storage")
					return receipt, addErr // Return specific storage error
				}
				tm.services.Logger.Info().Int("accepted", len(submittedTasks)).Float64("cost_per_task", costPerTask).Msg("✅ tasks created and added to storage queue")
			} else {
				tm.services.Logger.Warn().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("⚠️ no tasks created despite successful transaction")
			}
		}

		return receipt, nil
	}

	receipt, err := sendTx(context.Background(), nil, nil)
	return receipt, err
}

func (m *BatchTransactionManager) ProcessValidatorsStakes() {
	m.services.Logger.Info().Msg("🚀 performing validator stake and balance checks")
	// get the Eth balance
	bal, err := m.services.OwnerAccount.GetBalance()
	if err != nil {
		m.services.Logger.Err(err).Str("account", m.services.OwnerAccount.Address.String()).Msg("could not get eth balance on account")
		return
	}

	balAsFloat := Eth.ToFloat(bal)

	// check if the Ether balance is less than configured threshold
	if balAsFloat < m.services.Config.ValidatorConfig.EthLowThreshold {
		m.services.Logger.Warn().Float64("threshold", m.services.Config.ValidatorConfig.EthLowThreshold).Msg("⚠️ balance is below threshold")
	}

	for _, v := range m.services.Validators.validators {

		// get the baseTokenBalance on owner account as balance may change between checks
		baseTokenBalance, err := m.services.Basetoken.BalanceOf(nil, m.services.OwnerAccount.Address)
		if err != nil {
			m.services.Logger.Err(err).Msg("failed to get balance")
			return
		}

		ethAsFmt := fmt.Sprintf("%.8g Ξ", balAsFloat)
		baseAsFmt := fmt.Sprintf("%.8g %s", m.services.Config.BaseConfig.BaseToken.ToFloat(baseTokenBalance), m.services.Config.BaseConfig.BaseToken.Symbol)
		m.services.Logger.Info().Str("eth_bal", ethAsFmt).Str("basetoken_bal", baseAsFmt).Msg("wallet balances of owner account")

		v.ProcessValidatorStake(baseTokenBalance)
	}

}

func (m *BatchTransactionManager) TotalStaked(validator common.Address) error {
	// get the baseTokenBalance
	countBig, err := m.services.Engine.Engine.PendingValidatorWithdrawRequestsCount(nil, validator)
	if err != nil {
		m.services.Logger.Err(err).Msg("failed to get pending validator withdraw request count")
		return err
	}
	count := countBig.Uint64()

	totalStaked := big.NewInt(0)

	for i := uint64(0); i <= count; i++ {

		currCountBig := big.NewInt(int64(i))

		request, err := m.services.Engine.Engine.PendingValidatorWithdrawRequests(nil, validator, currCountBig)
		if err != nil {
			m.services.Logger.Err(err).Msg("failed to get pending validator withdraw request info")
			return err
		}

		unlockUnix := request.UnlockTime.Int64()

		if unlockUnix == 0 {
			continue
		}

		t := time.Unix(unlockUnix, 0)
		formatted := t.Format(time.DateTime)
		valAsFromat := m.services.Config.BaseConfig.BaseToken.FormatFixed(request.Amount)
		m.services.Logger.Info().Msgf("withdraw request #%d of %s, has unlock time %s", i, valAsFromat, formatted)

		totalStaked.Add(totalStaked, request.Amount)

	}
	valAsFromat := m.services.Config.BaseConfig.BaseToken.FormatFixed(totalStaked)
	m.services.Logger.Info().Msgf("Total staked for validator %s: %s", validator.String(), valAsFromat)
	return nil
}

func (tm *BatchTransactionManager) Start(ctx context.Context, enableBatchProcessing bool) error {
	tm.GoWithContext(ctx, tm.cumulativeGasUsed.Start)

	// if the validator stake / balance check is enabled, start the validator stake / balance check processor
	if tm.services.Config.ValidatorConfig.StakeCheck {
		stakeCheckInterval, err := time.ParseDuration(tm.services.Config.ValidatorConfig.StakeCheckInterval)
		if err != nil {
			return err
		}
		tm.ProcessValidatorsStakes()
		tm.Go(func() { tm.processValidatorStakePoller(ctx, stakeCheckInterval) })
	} else {
		tm.services.Logger.Warn().Msg("validator stake / balance checks are disabled!")
	}

	// if the ipfs incentive claim is enabled, start the ipfs claim processor
	if tm.services.Config.IPFS.IncentiveClaim {
		claimInterval, err := time.ParseDuration(tm.services.Config.IPFS.ClaimInterval)
		if err != nil {
			return err
		}
		tm.Go(func() { tm.startIpfsClaimProcessor(ctx, claimInterval) })
	}

	if enableBatchProcessing {
		if tm.services.Config.Miner.UsePolling {
			d, err := time.ParseDuration(tm.services.Config.Miner.PollingTime)
			if err != nil {
				return err
			}
			tm.Go(func() { tm.processBatchPoller(ctx, d) })
		} else {
			// process batch on block trigger (this hits external apis harder and not recommended)
			tm.GoWithContext(ctx, tm.processBatchBlockTrigger)
		}
	} else {
		if tm.services.Config.Claim.Enabled {
			tm.Go(func() { tm.batchClaimPoller(ctx, time.Duration(tm.services.Config.Claim.Delay)*time.Second) })
		}
	}
	return nil
}

func (tm *BatchTransactionManager) Stop() error {
	tm.services.Logger.Info().Msg("stopping batch manager")
	tm.Wait()
	tm.services.Logger.Info().Msg("batch manager stopped")
	return nil
}

// Add after the existing metrics interface implementation:

func (tm *BatchTransactionManager) GetSessionTime() string {
	sessionDuration := time.Since(tm.cumulativeGasUsed.sessionStartTime)
	hours := int(sessionDuration.Hours())
	minutes := int(sessionDuration.Minutes()) % 60
	seconds := int(sessionDuration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (tm *BatchTransactionManager) GetSolvedLastMinute() int64 {
	return tm.services.TaskTracker.GetSolvedLastMinute()
}

func (tm *BatchTransactionManager) GetSuccessCount() int64 {
	return tm.services.TaskTracker.GetSuccessCount()
}

func (tm *BatchTransactionManager) GetTotalCount() int64 {
	return tm.services.TaskTracker.GetTotalCount()
}

func (tm *BatchTransactionManager) GetSuccessRate() float64 {
	return tm.services.TaskTracker.GetSuccessRate()
}

func (tm *BatchTransactionManager) GetAverageSolutionRate() float64 {
	return tm.services.TaskTracker.GetAverageSolutionRate()
}

func (tm *BatchTransactionManager) GetAverageSolutionsPerMin() float64 {
	return tm.services.TaskTracker.GetAverageSolutionsPerMin()
}

func (tm *BatchTransactionManager) GetAverageSolvesPerMin() float64 {
	return tm.services.TaskTracker.GetAverageSolvesPerMin()
}

func (tm *BatchTransactionManager) GetTokenIncomePerMin() float64 {
	return tm.services.TaskTracker.GetAverageTasksPerPeriod() * tm.cumulativeGasUsed.rewardEMA.Average()
}

func (tm *BatchTransactionManager) GetTokenIncomePerHour() float64 {
	return tm.GetTokenIncomePerMin() * 60
}

func (tm *BatchTransactionManager) GetTokenIncomePerDay() float64 {
	return tm.GetTokenIncomePerHour() * 24
}

func (tm *BatchTransactionManager) GetIncomePerMin() float64 {
	return tm.GetTokenIncomePerMin() * tm.cumulativeGasUsed.basePriceEMA.Average()
}

func (tm *BatchTransactionManager) GetIncomePerHour() float64 {
	return tm.GetIncomePerMin() * 60
}

func (tm *BatchTransactionManager) GetIncomePerDay() float64 {
	return tm.GetIncomePerHour() * 24
}

func (tm *BatchTransactionManager) GetProfitPerMin() float64 {
	totalCostInUSD := tm.services.Config.BaseConfig.BaseToken.ToFloat(tm.cumulativeGasUsed.GetTotals()) * tm.cumulativeGasUsed.lastEthPrice
	timeSinceSessionStart := time.Since(tm.cumulativeGasUsed.sessionStartTime).Minutes()
	if timeSinceSessionStart <= 0 {
		return 0
	}
	averageCostsPerMin := totalCostInUSD / timeSinceSessionStart
	return tm.GetIncomePerMin() - averageCostsPerMin
}

func (tm *BatchTransactionManager) GetProfitPerHour() float64 {
	return tm.GetProfitPerMin() * 60
}

func (tm *BatchTransactionManager) GetProfitPerDay() float64 {
	return tm.GetProfitPerHour() * 24
}

// ProcessIpfsClaims queries the database for stored IPFS CIDs, gets signatures
// from the oracle and submits claims on-chain
// TODO: add support for querying other public IPFS providers for the Cid first to improve propagation speed before claiming
func (tm *BatchTransactionManager) processIpfsClaimsWithAccount(mainAccount *account.Account, ctx context.Context) error {

	// Configuration
	ipfsCfg := tm.services.Config.IPFS
	useBulkClaim := ipfsCfg.UseBulkClaim
	bulkClaimBatchSize := ipfsCfg.BulkClaimBatchSize
	maxSingleClaims := ipfsCfg.MaxSingleClaimsPerRun
	minAiusThreshold := ipfsCfg.MinAiusIncentiveThreshold // Use correct config

	if bulkClaimBatchSize <= 0 {
		bulkClaimBatchSize = 10 // Default batch size
	}
	if maxSingleClaims <= 0 {
		maxSingleClaims = 10 // Default single claims limit
	}

	// Get all Ipfs entries from the database that haven't been claimed (oldest first)
	// Fetch all available entries (-1 typically signifies 'all')
	ipfsEntries, err := tm.services.TaskStorage.GetIpfsCids(-1)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get unclaimed ipfs entries")
		return err
	}

	if len(ipfsEntries) == 0 {
		tm.services.Logger.Debug().Msg("no unclaimed ipfs entries to process")
		return nil
	}

	tm.services.Logger.Info().Int("fetched_count", len(ipfsEntries)).Bool("use_bulk", useBulkClaim).Int("batch_size", bulkClaimBatchSize).Int("max_single_claims", maxSingleClaims).Msg("processing ipfs claim entries")

	oracleClient := tm.services.IpfsOracle
	currentTime := time.Now().Unix() // Get current time once for the batch

	// Setup for grouping Bulk Claims by sender account
	// Map Address -> Batch List
	eligibleTasksByAccount := make(map[common.Address][]BulkClaimData)
	// Map Address -> Account Object (for easy retrieval when sending)
	accountMap := make(map[common.Address]*account.Account)
	accountMap[mainAccount.Address] = mainAccount
	for _, v := range tm.services.Validators.validators {
		accountMap[v.Account.Address] = v.Account
	}

	singleClaimsSent := 0 // Counter for single claims in this run

	for _, entry := range ipfsEntries {
		// Check if we should stop processing early within the loop
		select {
		case <-ctx.Done():
			tm.services.Logger.Info().Msg("context cancelled during ipfs claim processing loop")
			// If using bulk claim, process any pending batches before returning
			if useBulkClaim {
				for senderAddr, batch := range eligibleTasksByAccount {
					if len(batch) > 0 {
						bulkSenderAccount := accountMap[senderAddr]
						if bulkSenderAccount == nil {
							tm.services.Logger.Error().Str("sender", senderAddr.Hex()).Msg("Context Cancel: Could not find account object for bulk sender")
						} else {
							tm.services.Logger.Info().Str("sender", senderAddr.Hex()).Int("batch_size", len(batch)).Msg("Sending pending bulk claim batch due to context cancel")
							tm.sendBulkIpfsClaim(bulkSenderAccount, batch)
						}
					}
				}
			}
			return ctx.Err()
		default:
		}

		taskId := entry.TaskId
		taskIdStr := task.TaskId(taskId).String()

		// Declare claimLog here so it's always available
		var claimLog zerolog.Logger

		// Cast the local CID bytes to multihash for logging
		var cidStr string
		multihashLocal, castErr := mh.Cast(entry.Cid)
		if castErr != nil {
			// Log warning but don't create logger with invalid field
			tm.services.Logger.Warn().Err(castErr).Str("taskid", taskIdStr).Msg("failed to cast local cid bytes to multihash for logging")
			cidStr = "invalid_local_cid"
			// Create logger without the potentially problematic cid field
			claimLog = tm.services.Logger.With().Str("taskid", taskIdStr).Logger()
		} else {
			cidStr = cid.NewCidV0(multihashLocal).String()
			claimLog = tm.services.Logger.With().Str("taskid", taskIdStr).Str("local_cid", cidStr).Logger()
		}

		tm.ipfsClaimBackoffMutex.Lock()
		backoffUntil, inBackoff := tm.ipfsClaimBackoff[taskId]
		tm.ipfsClaimBackoffMutex.Unlock()

		if inBackoff && currentTime < backoffUntil {
			remainingSeconds := backoffUntil - currentTime
			claimLog.Info().Int64("remaining_seconds", remainingSeconds).Msg("task is in backoff period for IPFS claim, skipping")
			continue // Skip this task for now
		}

		// 1. Check Router Incentive Status
		incentiveAmount, err := tm.services.ArbiusRouter.Incentives(nil, taskId)
		if err != nil {
			claimLog.Error().Err(err).Msg("failed to check incentive availability on router")
			continue // Try next entry
		}

		if incentiveAmount.Cmp(big.NewInt(0)) == 0 {
			claimLog.Info().Msg("incentive already claimed or zero")
			// Cleanup storage and backoff
			if err := tm.services.TaskStorage.DeleteIpfsCid(taskId); err != nil {
				claimLog.Error().Err(err).Msg("failed to delete entry for already claimed incentive")
			}
			tm.ipfsClaimBackoffMutex.Lock()
			delete(tm.ipfsClaimBackoff, taskId)
			tm.ipfsClaimBackoffMutex.Unlock()
			continue // Move to next entry
		}
		claimLog.Debug().Str("incentive_amount", incentiveAmount.String()).Msg("incentive available")

		// Check: Minimum AIUS Threshold
		if minAiusThreshold > 0 { // Only check if threshold is set
			incentiveAmountFloat := tm.services.Config.BaseConfig.BaseToken.ToFloat(incentiveAmount)
			if incentiveAmountFloat < minAiusThreshold {
				claimLog.Info().
					Float64("incentive_aius", incentiveAmountFloat).
					Float64("threshold_aius", minAiusThreshold).
					Msg("incentive amount below minimum threshold, skipping")
				// Optionally delete from storage if we will never claim?
				// if err := tm.services.TaskStorage.DeleteIpfsCid(taskId); err != nil {
				// 	claimLog.Error().Err(err).Msg("failed to delete low-incentive entry")
				// }
				continue // Skip to the next entry
			}
		}

		// 2. Get On-Chain Solution Info from Engine
		solutionInfo, err := tm.services.Engine.GetSolution(taskId)
		if err != nil {
			claimLog.Error().Err(err).Msg("failed to get on-chain solution info from engine")
			continue // Try next entry
		}

		if solutionInfo.Blocktime == 0 {
			claimLog.Info().Msg("solution not yet confirmed on engine, will retry later")
			// No backoff here, just skip for this cycle
			continue
		}

		onChainValidator := solutionInfo.Validator
		onChainBlocktime := solutionInfo.Blocktime
		onChainCidBytes := solutionInfo.Cid // This is the CID we MUST use

		claimLog = claimLog.With().Str("onchain_validator", onChainValidator.Hex()).Uint64("onchain_blocktime", onChainBlocktime).Logger()

		// 3. Determine Sending Account based on Solver and Time Window
		var sendingAccount *account.Account
		isOurValidatorSolver := false
		validatorInstance := tm.services.Validators.GetValidatorByAddress(onChainValidator)
		if validatorInstance != nil {
			isOurValidatorSolver = true
		}

		hasPriorityWindowPassed := (currentTime >= int64(onChainBlocktime)+60)

		if isOurValidatorSolver && !hasPriorityWindowPassed {
			// Solver is ours, priority window active: Solver MUST send
			sendingAccount = validatorInstance.Account
			claimLog.Info().Str("sender", sendingAccount.Address.Hex()).Msg("Solver is our validator (priority active), will use specific account")
		} else if hasPriorityWindowPassed {
			// Priority window passed OR solver is external: Main account CAN send
			sendingAccount = mainAccount
			claimLog.Info().Str("sender", sendingAccount.Address.Hex()).Msg("Priority window passed or external solver, will use main account")
		} else { // Solver is external, priority window active
			claimLog.Info().Msg("External solver still has priority, skipping and applying backoff")
			// Add backoff because sending now would fail
			retryTime := int64(onChainBlocktime) + 61
			tm.ipfsClaimBackoffMutex.Lock()
			tm.ipfsClaimBackoff[taskId] = retryTime
			tm.ipfsClaimBackoffMutex.Unlock()
			continue
		}

		// Sanity check - should not happen based on logic above
		if sendingAccount == nil {
			claimLog.Error().Msg("Logical error: Could not determine sending account")
			continue
		}
		sendingAccountAddress := sendingAccount.Address

		// 4. Get Signatures for the CORRECT (On-Chain) CID
		onChainCidHex := hex.EncodeToString(onChainCidBytes)
		claimLog.Info().Str("onchain_cid", "0x"+onChainCidHex).Msg("getting signatures for on-chain CID")

		multihash, err := mh.Cast(onChainCidBytes)
		if err != nil {
			claimLog.Error().Err(err).Str("onchain_cid_hex", onChainCidHex).Msg("failed to cast on-chain cid bytes to multihash")
			continue
		}
		onChainCidStr := cid.NewCidV0(multihash).String()

		signatures, err := oracleClient.GetSignaturesForCID(ctx, onChainCidStr)
		if err != nil {
			if strings.Contains(err.Error(), "504 Gateway Time-out") {
				claimLog.Warn().Msg("Oracle timed out (504) getting signatures, scheduling retry after 30s")
				retryTime := time.Now().Unix() + 30 // 30-second backoff for oracle timeout
				tm.ipfsClaimBackoffMutex.Lock()
				tm.ipfsClaimBackoff[taskId] = retryTime
				tm.ipfsClaimBackoffMutex.Unlock()
			} else {
				claimLog.Error().Err(err).Str("onchain_cid_str", onChainCidStr).Msg("failed to get signatures for on-chain cid from oracle")
				// Optionally add a different backoff strategy for other persistent oracle errors here
				// Add generic backoff for other oracle errors
				retryTime := time.Now().Unix() + 60 // Longer backoff
				tm.ipfsClaimBackoffMutex.Lock()
				tm.ipfsClaimBackoff[taskId] = retryTime
				tm.ipfsClaimBackoffMutex.Unlock()
			}
			continue // Retry later
		}

		const requiredSignatures = 1              // Using default as config field missing
		if len(signatures) < requiredSignatures { // Use default min signatures
			claimLog.Warn().Int("signatures_have", len(signatures)).Int("signatures_needed", requiredSignatures).Msg("insufficient signatures returned from oracle, will retry later")
			// Consider adding a backoff here too, similar to the 504 case
			retryTime := time.Now().Unix() + 60 // Example: 60 second backoff for insufficient signatures
			tm.ipfsClaimBackoffMutex.Lock()
			tm.ipfsClaimBackoff[taskId] = retryTime
			tm.ipfsClaimBackoffMutex.Unlock()
			continue // Retry later
		}

		// Convert signatures to the format required by the contract
		arbiusSignatures := make([]arbiusrouterv1.Signature, len(signatures))
		for i, signature := range signatures {
			// TODO: Ensure oracle provides sorted signatures, or sort here.
			arbiusSignatures[i] = arbiusrouterv1.Signature{
				Signer:    signature.Signer,
				Signature: common.FromHex(signature.Signature), // Assuming signature is hex string
			}
		}

		// 5. Add to Batch or Send Single Claim
		if useBulkClaim {
			// Add task data to the batch corresponding to the determined sender account
			eligibleTasksByAccount[sendingAccountAddress] = append(eligibleTasksByAccount[sendingAccountAddress], BulkClaimData{
				TaskID:       taskId,
				Signatures:   arbiusSignatures,
				ClaimAccount: sendingAccount, // Store the account determined for this claim
			})
			currentBatchSize := len(eligibleTasksByAccount[sendingAccountAddress])
			claimLog.Info().Str("sender", sendingAccountAddress.Hex()).Int("current_batch_size", currentBatchSize).Msg("added task to bulk claim batch")

			// If this specific sender's batch is full, send it
			if currentBatchSize >= bulkClaimBatchSize {
				batchToSend := eligibleTasksByAccount[sendingAccountAddress]
				bulkSenderAccount := accountMap[sendingAccountAddress] // Get the account object
				if bulkSenderAccount == nil {
					claimLog.Error().Str("sender", sendingAccountAddress.Hex()).Msg("Could not find account object for full bulk batch sender address")
					// How to handle? For now, just log and the batch will be attempted at the end.
				} else {
					claimLog.Info().Str("sender", sendingAccountAddress.Hex()).Int("batch_size", len(batchToSend)).Msg("Sending full bulk claim batch")
					tm.sendBulkIpfsClaim(bulkSenderAccount, batchToSend)
				}
				eligibleTasksByAccount[sendingAccountAddress] = nil // Reset this specific batch
			}
		} else {
			// Send single claim immediately using the determined sendingAccount
			tm.sendSingleIpfsClaim(sendingAccount, taskId, arbiusSignatures, onChainBlocktime, claimLog)
			singleClaimsSent++ // Increment counter after attempting single claim
		}

		// Check if single claim limit is reached
		if !useBulkClaim && singleClaimsSent >= maxSingleClaims {
			tm.services.Logger.Info().Int("limit", maxSingleClaims).Msg("Reached single IPFS claim limit for this run")
			break // Stop processing more entries in this run
		}
	} // End of loop through ipfsEntries

	// After loop, send any remaining partial batches if using bulk claim
	if useBulkClaim {
		for senderAddr, remainingBatch := range eligibleTasksByAccount {
			if len(remainingBatch) > 0 {
				bulkSenderAccount := accountMap[senderAddr]
				if bulkSenderAccount == nil {
					tm.services.Logger.Error().Str("sender", senderAddr.Hex()).Int("batch_size", len(remainingBatch)).Msg("Could not find account object for final partial bulk sender address")
				} else {
					tm.services.Logger.Info().Str("sender", senderAddr.Hex()).Int("batch_size", len(remainingBatch)).Msg("Sending remaining partial bulk claim batch")
					tm.sendBulkIpfsClaim(bulkSenderAccount, remainingBatch)
				}
			}
		}
	}

	return nil
}

// StartIpfsClaimProcessor starts a background process to periodically check for and process IPFS claims
func (tm *BatchTransactionManager) startIpfsClaimProcessor(appQuit context.Context, claimInterval time.Duration) {
	ticker := time.NewTicker(claimInterval)
	defer ticker.Stop()

	tm.services.Logger.Info().Msgf("started ipfs claim processor with interval %s", claimInterval)

	for {
		select {
		case <-appQuit.Done():
			tm.services.Logger.Info().Msg("shutting down ipfs claim processor")
			return
		case <-ticker.C:
			// TODO: make this use task accounts instead of sender account
			if err := tm.processIpfsClaimsWithAccount(tm.services.OwnerAccount, appQuit); err != nil {
				if err != context.Canceled {
					tm.services.Logger.Error().Err(err).Msg("error processing ipfs claims")
				}
			}
		}
	}
}

// Helper Functions

// calculateGasLimit determines the gas limit based on mode and formula.
// Returns the calculated limit (0 for estimation, formula otherwise).
func calculateGasLimit(isEstimationMode bool, batchSize int, baseLimit, perItemGas uint64) uint64 {
	if isEstimationMode {
		return 0 // Force estimation by setting limit to 0
	}
	// Ensure non-negative batchSize for calculation
	if batchSize < 0 {
		batchSize = 0
	}
	return baseLimit + uint64(batchSize)*perItemGas
}

// logGasLimitDetails logs the comparison between actual and hardcoded gas limits.
func logGasLimitDetails(logger zerolog.Logger, tx *types.Transaction, isEstimationMode bool, hardcodedGasLimit uint64, functionName string, context map[string]interface{}) {
	if tx == nil {
		return
	}

	actualGasLimit := tx.Gas()
	logFields := map[string]interface{}{
		"tx_gas_limit":        actualGasLimit,
		"function":            functionName,
		"hardcoded_gas_limit": hardcodedGasLimit,
	}
	// Merge context fields
	for k, v := range context {
		if k != "tx_gas_limit" && k != "function" && k != "hardcoded_gas_limit" { // Avoid overriding core fields
			logFields[k] = v
		}
	}

	if isEstimationMode {
		logFields["mode"] = "estimate"
		if hardcodedGasLimit > 0 {
			diff := (float64(actualGasLimit) - float64(hardcodedGasLimit)) / float64(hardcodedGasLimit) * 100.0
			logFields["diff_percent"] = fmt.Sprintf("%.2f%%", diff)
		} else {
			logFields["diff_percent"] = "N/A"
		}
	} else {
		logFields["mode"] = "normal"
	}
	logger.Info().Fields(logFields).Msg("transaction prepared with set gas limit")
}

// processReceiptAndMetrics handles receipt processing, metrics, and logging gas used.
// It calls the provided metricFunc (e.g., AddCommitment, AddSolution) with the calculated txCost.
// Returns the calculated txCost.
func (tm *BatchTransactionManager) processReceiptAndMetrics(receipt *types.Receipt, metricFunc func(*big.Int), batchSize int, itemName string) *big.Int {
	txCost := big.NewInt(0) // Initialize cost to zero
	if receipt == nil {
		tm.services.Logger.Warn().Str("item_name", itemName).Int("batch_size", batchSize).Msg("assertion: receipt is nil in processReceiptAndMetrics")
		return txCost
	}
	if metricFunc == nil {
		tm.services.Logger.Error().Str("item_name", itemName).Msg("assertion: metricFunc is nil in processReceiptAndMetrics")
		return txCost
	}

	txCost.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
	gasPerItem := 0.0
	if batchSize > 0 {
		gasPerItem = float64(receipt.GasUsed) / float64(batchSize)
	}

	txCostEthStr := fmt.Sprintf("%.8f", tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost))

	tm.services.Logger.Info().
		Uint64("gas_used", receipt.GasUsed).
		Float64(fmt.Sprintf("gas_per_%s", itemName), gasPerItem).
		Str("tx_cost_eth", txCostEthStr).       // Log cost in ETH
		Str("txhash", receipt.TxHash.String()). // Add tx hash for context
		Msgf("**** %s gas usage *****", itemName)

	metricFunc(txCost)

	return txCost
}

// sendSingleIpfsClaim handles the logic for sending a single IPFS incentive claim.
func (tm *BatchTransactionManager) sendSingleIpfsClaim(
	claimAccount *account.Account, // The account to use for sending the tx
	taskId task.TaskId,
	arbiusSignatures []arbiusrouterv1.Signature,
	onChainBlocktime uint64, // Needed for backoff calculation on revert
	claimLog zerolog.Logger, // Use the logger passed down
) {
	currentTime := time.Now().Unix()

	claimLog.Info().Int("signature_count", len(arbiusSignatures)).Str("using_account", claimAccount.Address.Hex()).Msg("attempting single ipfs claim transaction")

	// Get Gas Opts for the specific claimAccount
	claimGasPrice, claimGasFeeCap, claimGasFeeTip, _ := claimAccount.Client.GasPriceOracle(false)
	opts := claimAccount.GetOpts(0, claimGasPrice, claimGasFeeCap, claimGasFeeTip)

	tx, err := tm.services.ArbiusRouter.ClaimIncentive(opts, taskId, arbiusSignatures)

	// Handle Pre-Send Errors (including reverts caught by estimation)
	if err != nil {
		retryTime := tm.handleIpfsClaimError(err, taskId, onChainBlocktime, claimLog)
		if retryTime > 0 {
			tm.ipfsClaimBackoffMutex.Lock()
			tm.ipfsClaimBackoff[taskId] = retryTime
			tm.ipfsClaimBackoffMutex.Unlock()
			remainingSeconds := retryTime - currentTime
			claimLog.Info().Int64("remaining_seconds", remainingSeconds).Msg("setting backoff for ipfs claim after pre-send error")
		}
		return
	}
	if tx == nil {
		claimLog.Error().Msg("assertion: single ipfs claim tx is nil but no error reported")
		return // Avoid panic on nil tx
	}

	// Handle Successful Send
	claimLog.Info().Str("txhash", tx.Hash().String()).Msg("single ipfs claim transaction sent, waiting for confirmation...")
	var success bool

	if *tx.To() == (common.Address{}) {
		claimLog.Info().Str("taskId", task.TaskId(taskId).String()).Msg("mock router contract, skipping confirmation check")
		success = true
	} else {
		_, success, _, err = claimAccount.WaitForConfirmedTx(tx) // Use claimAccount here
		if err != nil {
			claimLog.Error().Err(err).Str("txhash", tx.Hash().String()).Msg("failed waiting for single ipfs claim confirmation, will retry later")
			// Apply backoff even for wait errors
			retryTime := time.Now().Unix() + 60 // Generic 60s backoff for wait errors
			tm.ipfsClaimBackoffMutex.Lock()
			tm.ipfsClaimBackoff[taskId] = retryTime
			tm.ipfsClaimBackoffMutex.Unlock()
			return
		}
	}

	if success {
		claimLog.Info().Str("txhash", tx.Hash().String()).Msg("successfully claimed incentive transaction confirmed")
		// Delete from database only on confirmed success
		if err := tm.services.TaskStorage.DeleteIpfsCid(taskId); err != nil {
			claimLog.Error().Err(err).Msg("failed to delete entry from storage after successful claim")
		} else {
			claimLog.Info().Msg("deleted entry from storage")
		}
		// Clear any backoff on success
		tm.ipfsClaimBackoffMutex.Lock()
		delete(tm.ipfsClaimBackoff, taskId)
		tm.ipfsClaimBackoffMutex.Unlock()
	} else {
		claimLog.Warn().Str("txhash", tx.Hash().String()).Msg("single ipfs claim transaction reverted on-chain after sending, will retry later")
		// Add backoff for post-send reverts
		retryTime := int64(onChainBlocktime) + 61 // Use blocktime for TimeNotPassed reverts
		tm.ipfsClaimBackoffMutex.Lock()
		tm.ipfsClaimBackoff[taskId] = retryTime
		remainingSeconds := retryTime - currentTime
		claimLog.Info().Int64("remaining_seconds", remainingSeconds).Msg("setting backoff after post-send revert")
		tm.ipfsClaimBackoffMutex.Unlock()
	}
}

// sendBulkIpfsClaim handles the logic for sending a bulk IPFS incentive claim.
func (tm *BatchTransactionManager) sendBulkIpfsClaim(
	senderAccount *account.Account, // The single account used to send the bulk transaction
	batch []BulkClaimData,
) {
	if len(batch) == 0 {
		return
	}

	bulkLog := tm.services.Logger.With().Int("batch_size", len(batch)).Str("sender_account", senderAccount.Address.Hex()).Logger()
	bulkLog.Info().Msg("attempting bulk ipfs claim transaction")

	// Prepare arguments for bulkClaimIncentive
	taskIds := make([][32]byte, len(batch))
	var flatSignatures []arbiusrouterv1.Signature
	var sigsPerTask uint64 = 0
	if len(batch) > 0 {
		sigsPerTask = uint64(len(batch[0].Signatures)) // Assume consistent number of sigs
	}

	taskIDToBlocktimeMap := make(map[task.TaskId]uint64) // Store blocktimes for backoff on revert

	for i, data := range batch {
		taskIds[i] = data.TaskID
		if uint64(len(data.Signatures)) != sigsPerTask {
			bulkLog.Error().Str("taskid", task.TaskId(data.TaskID).String()).Int("expected_sigs", int(sigsPerTask)).Int("actual_sigs", len(data.Signatures)).Msg("Inconsistent number of signatures in bulk batch, cannot send.")
			// Handle this error - potentially skip this task or fail the batch?
			// For now, failing the batch to be safe.
			// We could add backoff to all tasks in the batch here.
			return
		}
		flatSignatures = append(flatSignatures, data.Signatures...)

		// Fetch blocktime again for reliable backoff calculation in case of revert
		solutionInfo, err := tm.services.Engine.GetSolution(data.TaskID)
		if err != nil || solutionInfo.Blocktime == 0 {
			bulkLog.Warn().Err(err).Str("taskid", task.TaskId(data.TaskID).String()).Msg("Could not get blocktime for task in bulk batch for backoff calculation")
			taskIDToBlocktimeMap[data.TaskID] = 0 // Mark as unknown
		} else {
			taskIDToBlocktimeMap[data.TaskID] = solutionInfo.Blocktime
		}
	}

	// Get Gas Opts for the senderAccount
	claimGasPrice, claimGasFeeCap, claimGasFeeTip, _ := senderAccount.Client.GasPriceOracle(false)
	opts := senderAccount.GetOpts(0, claimGasPrice, claimGasFeeCap, claimGasFeeTip)

	// Calculate gas limit for bulk claim
	isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode // Assuming global setting applies
	opts.GasLimit = calculateGasLimit(isEstimationMode, len(taskIds), BaseGasLimitForBulkClaimIncentive, BulkClaimIncentiveGasPerItem)
	if isEstimationMode && tm.services.Config.Miner.GasEstimationMargin > 0 {
		opts.GasMargin = tm.services.Config.Miner.GasEstimationMargin
	}
	opts.NoSend = true // We need to sign and send manually

	// Call the contract method via the wrapper/account
	tx, err := senderAccount.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(innerOpts *bind.TransactOpts) (interface{}, error) {
		// Apply gas settings from outer opts to inner opts
		innerOpts.GasLimit = opts.GasLimit
		innerOpts.GasMargin = opts.GasMargin
		innerOpts.NoSend = true // Ensure inner call doesn't send

		txToSign, err := tm.services.ArbiusRouter.BulkClaimIncentive(innerOpts, taskIds, flatSignatures, big.NewInt(int64(sigsPerTask)))

		if err != nil {
			return nil, err // Propagate contract preparation errors
		}
		// Send the signed transaction using the senderAccount
		return senderAccount.SendSignedTransaction(txToSign)
	})

	// Handle Pre-Send Errors (including reverts caught by estimation)
	if err != nil {
		bulkLog.Error().Err(err).Msg("failed to estimate/prepare bulk ipfs claim transaction")
		// Apply backoff to all tasks in the batch on pre-send failure
		tm.ipfsClaimBackoffMutex.Lock()
		currentTime := time.Now().Unix()
		for _, taskId := range taskIds {
			blocktime := taskIDToBlocktimeMap[taskId]
			retryTime := tm.calculateBackoffForError(err, blocktime, currentTime) // Use helper
			tm.ipfsClaimBackoff[taskId] = retryTime
			bulkLog.Debug().Str("taskid", task.TaskId(taskId).String()).Int64("backoff_until", retryTime).Msg("setting backoff for task in failed bulk batch (pre-send)")
		}
		tm.ipfsClaimBackoffMutex.Unlock()
		return
	}
	if tx == nil {
		bulkLog.Error().Msg("assertion: bulk ipfs claim tx is nil but no error reported")
		return // Avoid panic
	}

	// Handle Successful Send
	bulkLog.Info().Str("txhash", tx.Hash().String()).Msg("bulk ipfs claim transaction sent, waiting for confirmation...")
	var success bool
	var receipt *types.Receipt

	if *tx.To() == (common.Address{}) {
		bulkLog.Info().Msg("mock router contract, skipping confirmation check for bulk")
		success = true
	} else {
		receipt, success, _, err = senderAccount.WaitForConfirmedTx(tx)
		if err != nil {
			bulkLog.Error().Err(err).Str("txhash", tx.Hash().String()).Msg("failed waiting for bulk ipfs claim confirmation, will retry batch later")
			// Apply backoff to all tasks in the batch on wait error
			tm.ipfsClaimBackoffMutex.Lock()
			retryTime := time.Now().Unix() + 60 // Generic 60s backoff
			for _, taskId := range taskIds {
				tm.ipfsClaimBackoff[taskId] = retryTime
				bulkLog.Debug().Str("taskid", task.TaskId(taskId).String()).Int64("backoff_until", retryTime).Msg("setting backoff for task in failed bulk batch (wait error)")
			}
			tm.ipfsClaimBackoffMutex.Unlock()
			return
		}
	}

	if success {
		bulkLog.Info().Str("txhash", tx.Hash().String()).Msg("bulk ipfs claim transaction confirmed")

		claimedTasks := make(map[task.TaskId]bool)
		if receipt != nil { // Only parse logs if we have a real receipt
			for _, vlog := range receipt.Logs {
				if len(vlog.Topics) > 0 && vlog.Topics[0] == tm.incentiveClaimedEvent {
					parsed, parseErr := tm.services.ArbiusRouter.ParseIncentiveClaimed(*vlog)
					if parseErr != nil {
						bulkLog.Error().Err(parseErr).Msg("failed to parse IncentiveClaimed event")
						continue
					}
					claimedTasks[task.TaskId(parsed.Taskid)] = true
				}
			}
		} else if *tx.To() == (common.Address{}) {
			// Mock success - assume all tasks in batch were claimed for mock scenario
			bulkLog.Warn().Msg("mock router success, assuming all tasks in batch claimed")
			for _, taskId := range taskIds {
				claimedTasks[taskId] = true
			}
		}

		bulkLog.Info().Int("claimed_tasks", len(claimedTasks)).Msg("bulk ipfs claim transaction confirmed")

		// Process results: delete successful, backoff unsuccessful
		var tasksToDelete []task.TaskId
		tm.ipfsClaimBackoffMutex.Lock()
		currentTime := time.Now().Unix()
		for _, taskId := range taskIds {
			if claimedTasks[taskId] {
				tasksToDelete = append(tasksToDelete, taskId)
				delete(tm.ipfsClaimBackoff, taskId) // Clear backoff on success
			} else {
				// Task was in batch but not in success logs (or logs couldn't be parsed)
				bulkLog.Warn().Str("taskid", task.TaskId(taskId).String()).Msg("task included in bulk claim but not confirmed claimed via logs, applying backoff")
				blocktime := taskIDToBlocktimeMap[taskId]
				// Use a generic backoff or blocktime-based if available
				retryTime := int64(blocktime) + 61
				if blocktime == 0 {
					retryTime = currentTime + 60 // Fallback if blocktime unknown
				}
				tm.ipfsClaimBackoff[taskId] = retryTime
			}
		}
		tm.ipfsClaimBackoffMutex.Unlock()

		if len(tasksToDelete) > 0 {
			deletedCount := 0
			for _, tid := range tasksToDelete {
				if err := tm.services.TaskStorage.DeleteIpfsCid(tid); err != nil {
					bulkLog.Error().Err(err).Str("taskid", task.TaskId(tid).String()).Msg("failed to delete successfully claimed entry from storage")
				} else {
					deletedCount++
				}
			}
			if deletedCount > 0 {
				bulkLog.Info().Int("deleted_count", deletedCount).Msg("deleted successfully claimed entries from storage")
			}
		}

	} else { // Bulk transaction reverted
		bulkLog.Warn().Str("txhash", tx.Hash().String()).Msg("bulk ipfs claim transaction reverted on-chain after sending, applying backoff to all in batch")
		// Apply backoff to all tasks in the batch
		tm.ipfsClaimBackoffMutex.Lock()
		currentTime := time.Now().Unix()
		for _, taskId := range taskIds {
			blocktime := taskIDToBlocktimeMap[taskId]
			// Assume revert might be TimeNotPassed for at least one task
			retryTime := int64(blocktime) + 61
			if blocktime == 0 {
				retryTime = currentTime + 60 // Fallback
			}
			tm.ipfsClaimBackoff[taskId] = retryTime
			bulkLog.Warn().Str("taskid", task.TaskId(taskId).String()).Int64("backoff_until", retryTime).Msg("setting backoff for task in failed bulk batch (reverted)")
		}
		tm.ipfsClaimBackoffMutex.Unlock()
	}
}

// handleIpfsClaimError analyzes the error from a claim attempt and logs appropriately.
// Returns the calculated retry time (unix timestamp) for backoff, or 0 if no specific backoff applies.
func (tm *BatchTransactionManager) handleIpfsClaimError(err error, taskId task.TaskId, onChainBlocktime uint64, claimLog zerolog.Logger) int64 {
	errStr := err.Error()
	retryTime := time.Now().Unix() + 10 // Default short backoff

	if strings.Contains(errStr, "TimeNotPassed") {
		claimLog.Info().Msg("claim reverted during estimation/send (TimeNotPassed), will retry later")
		retryTime = int64(onChainBlocktime) + 61 // 60s + 1s buffer
	} else if strings.Contains(errStr, "InvalidSignature") {
		claimLog.Warn().Msg("claim reverted during estimation/send (InvalidSignature), check oracle/CID, will retry later")
		// Longer backoff might be suitable
		retryTime = time.Now().Unix() + 120
	} else if strings.Contains(errStr, "InsufficientSignatures") {
		claimLog.Warn().Msg("claim reverted during estimation/send (InsufficientSignatures), check oracle, will retry later")
		// Longer backoff might be suitable
		retryTime = time.Now().Unix() + 120
	} else if strings.Contains(errStr, "InvalidValidator") {
		claimLog.Error().Msg("claim reverted during estimation/send (InvalidValidator), check config")
		// Maybe no retry? Or very long backoff. For now, default.
	} else if strings.Contains(errStr, "SignersNotSorted") {
		claimLog.Error().Msg("claim reverted during estimation/send (SignersNotSorted), check oracle logic")
		// Might be a persistent oracle issue. Long backoff.
		retryTime = time.Now().Unix() + 300
	} else if strings.Contains(errStr, "nonce too low") || strings.Contains(errStr, "replacement transaction underpriced") {
		claimLog.Warn().Err(err).Msg("Nonce or gas price issue during IPFS claim, relying on NonceManagerWrapper retry")
		// Return 0 as NonceManager handles this internally
		return 0
	} else if strings.Contains(errStr, "insufficient funds") {
		claimLog.Error().Err(err).Msg("Insufficient funds for IPFS claim transaction")
		// Long backoff, requires user intervention
		retryTime = time.Now().Unix() + 3600
	} else {
		// Includes network errors, other reverts
		claimLog.Error().Err(err).Msg("failed to estimate/prepare/send ipfs claim transaction")
		// Generic longer backoff for unknown errors
		retryTime = time.Now().Unix() + 60
	}
	return retryTime
}

// calculateBackoffForError determines the appropriate backoff time based on an error string.
// Used primarily for bulk claim failures where individual task context might be lost.
func (tm *BatchTransactionManager) calculateBackoffForError(err error, blocktime uint64, currentTime int64) int64 {
	errStr := err.Error()
	retryTime := currentTime + 60 // Default backoff

	if strings.Contains(errStr, "TimeNotPassed") && blocktime > 0 {
		retryTime = int64(blocktime) + 61
	} else if strings.Contains(errStr, "InvalidSignature") || strings.Contains(errStr, "InsufficientSignatures") {
		retryTime = currentTime + 120
	} else if strings.Contains(errStr, "SignersNotSorted") {
		retryTime = currentTime + 300
	} else if strings.Contains(errStr, "InvalidValidator") {
		// Maybe longer?
		retryTime = currentTime + 3600
	} else if strings.Contains(errStr, "nonce") || strings.Contains(errStr, "replacement transaction") {
		// Nonce manager handles this, short backoff just in case
		retryTime = currentTime + 10
	} else if strings.Contains(errStr, "insufficient funds") {
		retryTime = currentTime + 3600
	}

	// Ensure backoff is always in the future
	if retryTime <= currentTime {
		retryTime = currentTime + 60
	}
	return retryTime
}
