package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gobius/account"
	"gobius/bindings/arbiusrouterv1"
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
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const profitEstimateBatchSize = 200

type CacheItem struct {
	Value      any
	LastUpdate time.Time
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
	batchMode                     int // 0 - single commitment/soution cycle with no batching, 1 - normal batch system, 2 - do not use this mode
	commitments                   [][32]byte
	solutions                     []BatchSolution
	encodedTaskData               []byte
	taskAccounts                  []*account.Account
	validatorIndex                int
	validators                    Validators
	minClaimSolutionTime          uint64
	minContestationVotePeriodTime uint64
	cache                         *Cache
	//	appQuit                    context.Context
	wg *sync.WaitGroup
	sync.Mutex
}

type Validators []*Validator

func (vl Validators) GetValidatorByAddress(addr common.Address) *Validator {
	for _, v := range vl {
		if v.ValidatorAddress() == addr {
			return v
		}
	}
	return nil
}

/*
implementation for the metrics interface
*/
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
	for _, v := range tm.validators {
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

	basefeeinEth := tm.services.Eth.ToFloat(basefee)
	basefeeinGwei := basefeeinEth * 1000000000

	// Use the PriceOracle interface to get prices
	basePrice, ethPrice, err = tm.services.OracleProvider.GetPrices()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("could not get prices from oracle!")
		// Fallback logic remains for now, but consider making it configurable or removing it
		// if the oracle should be the single source of truth.
		if tm.services.Config.BaseConfig.TestnetType > 0 {
			tm.services.Logger.Warn().Msg("Oracle failed, using default testnet prices (30, 2000)")
			basePrice, ethPrice = 30, 2000
		} else {
			// Consider if a fallback is appropriate on mainnet or if we should return error
			tm.services.Logger.Warn().Msg("Oracle failed, using default mainnet prices (30, 2000)")
			basePrice, ethPrice = 30, 2000
		}
		// If using fallbacks, reset the error so the function can continue
		err = nil
	}

	tm.cache.Set("base_price", basePrice)
	tm.cache.Set("eth_price", ethPrice)

	submitTasksBatchUSD := 0.0
	submitTasksBatch := (129_000 * basefeeinEth * profitEstimateBatchSize)
	submitTasksBatchUSD = submitTasksBatch * ethPrice

	signalCommitmentBatch := (27_900 * basefeeinEth * profitEstimateBatchSize)
	signalCommitmentBatchUSD := signalCommitmentBatch * ethPrice

	submitSolutionBatch := (128_500 * basefeeinEth * profitEstimateBatchSize)
	submitSolutionBatchUSD := submitSolutionBatch * ethPrice

	claimTasksUSD := 0.0
	claimTasks := (47_300 * basefeeinEth * profitEstimateBatchSize)
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
func (tm *BatchTransactionManager) batchClaimPoller(appQuit context.Context, wg *sync.WaitGroup, pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()
	defer wg.Done()

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
					tm.services.Logger.Error().Err(err).Msg("could not get keys from storage")
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
						basefeeinEth := tm.services.Eth.ToFloat(basefeeBig)
						basefeeinGwei := basefeeinEth * 1000000000

						if basefeeinGwei > tm.services.Config.Claim.MaxGas {
							tm.services.Logger.Warn().Msgf("** base gas is too high to claim **")
							continue
						}
					}

					tm.processBulkClaim(tm.services.SenderOwnerAccount, claims, tm.services.Config.Claim.MinClaims, claimBatchSize)
				}
			}
		}
	}
}

func (tm *BatchTransactionManager) processBatchBlockTrigger(appQuit context.Context, wg *sync.WaitGroup) {
	var batchWG sync.WaitGroup
	defer wg.Done()

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

func (tm *BatchTransactionManager) processValidatorStakePoller(appQuit context.Context, wg *sync.WaitGroup, pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()
	defer wg.Done()

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
func (tm *BatchTransactionManager) processBatchPoller(appQuit context.Context, wg *sync.WaitGroup, pollingtime time.Duration) {
	var batchWG sync.WaitGroup
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()
	defer wg.Done()

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

	//if minProfit > 0 {
	if usegwei {
		isProfitable = profitLevel <= minProfit
	} else {
		isProfitable = profitLevel >= minProfit
	}
	//}

	makeTasks := tm.services.Config.BatchTasks.Enabled && isProfitable && totalTasks < int64(tm.services.Config.BatchTasks.MinTasksInQueue) && tm.services.Config.BatchTasks.BatchSize > 0

	// log above for debugging
	tm.services.Logger.Info().Bool("make_tasks", makeTasks).Int64("total_tasks", totalTasks).Int("min_tasks_in_queue", tm.services.Config.BatchTasks.MinTasksInQueue).Int("batch_size", tm.services.Config.BatchTasks.BatchSize).Msg("batch conditions")

	taskBatchCount := tm.services.Config.BatchTasks.NumberOfBatches
	taskBatchSize := tm.services.Config.BatchTasks.BatchSize
	hoardMode := tm.services.Config.BatchTasks.Enabled && tm.services.Config.BatchTasks.HoardMode && baseFee <= tm.services.Config.BatchTasks.HoardMinGasPrice && totalTasks < int64(tm.services.Config.BatchTasks.HoardMaxQueueSize)
	if !makeTasks && hoardMode {
		tm.services.Logger.Warn().Msgf("** task hoard mode on, and basefee is at or below wanted gas price & queue len is below threshold**")
		makeTasks = true
		taskBatchCount = tm.services.Config.BatchTasks.HoardModeNumberOfBatches
		taskBatchSize = tm.services.Config.BatchTasks.HoardModeBatchSize
	}
	//		tm.services.SenderOwnerAccount.UpdateNonce()

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

	totalSolutions, totalClaims, err := tm.services.TaskStorage.TotalSolutionsAndClaims()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get total solutions")
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
		// 47_300 is the average gas per task
		// 1_000_000_000.0 is to adjust for gas price in gwei
		claimTasks := (47_300.0 / 1_000_000_000.0 * baseFee * float64(len(claims)))

		tm.services.Logger.Warn().Msgf("** CHEAPEST BATCH WE CAN SEND OUT**")
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
			// debug for now
			tm.services.Logger.Warn().Msgf("** minclaim lever: %.8g **", claimMinReward)

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
		tm.Lock()
		defer tm.Unlock()
		var tasks storage.TaskDataSlice
		var err error

		tasks, err = tm.services.TaskStorage.GetPendingCommitments(batchSize)
		if err != nil {
			tm.services.Logger.Err(err).Msg("failed to get commitments from storage")
			return nil, err
		}

		if noChecks {
			return tasks, nil
		}

		var commitmentsToDelete []task.TaskId
		var validTasksForBatch storage.TaskDataSlice
		var tasksToUpdateToStatus2 []task.TaskId // New list to track tasks needing status update

		for _, t := range tasks {
			// Ensure task has a commitment hash locally before checking on-chain
			if t.Commitment == [32]byte{} {
				tm.services.Logger.Warn().Str("taskid", t.TaskId.String()).Msg("Task in pending commitments has zero commitment hash, skipping.")
				continue // Skip tasks with invalid local state
			}

			block, err := tm.services.Engine.Engine.Commitments(nil, t.Commitment)
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("taskid", t.TaskId.String()).Msg("error checking on-chain commitment status, skipping task")
				continue
			}

			blockNo := block.Uint64()
			if blockNo > 0 {
				// Commitment already exists on-chain
				commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
				tasksToUpdateToStatus2 = append(tasksToUpdateToStatus2, t.TaskId) // Mark for status update
			} else {
				// Commitment not on-chain, add to the batch to be sent
				validTasksForBatch = append(validTasksForBatch, t)
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
		tm.Lock()
		defer tm.Unlock()

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

		blockTime := time.Unix(int64(blockInfo.Time()), 0)

		var validator *Validator = nil
		validatorHighestMin := int64(-1)
		for _, v := range tm.validators {
			lastSubmission, maxSols, err := v.MaxSubmissions(blockTime)
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

		if noChecks {
			return validator, tasks, nil
		}

		var commitmentsToDelete []task.TaskId
		var solutionsToDelete []task.TaskId
		var validTasksForBatch storage.TaskDataSlice
		var claimsToDelete []task.TaskId

		for _, t := range tasks {
			res, err := tm.services.Engine.Engine.Solutions(nil, t.TaskId)
			if err != nil {
				tm.services.Logger.Err(err).Str("taskid", t.TaskId.String()).Msg("error getting on-chain solution information, skipping task in batch prep")
				continue
			}

			if res.Blocktime > 0 {
				if res.Claimed {
					// Solution is already claimed on-chain
					tm.services.Logger.Info().Str("taskid", t.TaskId.String()).Str("validator", res.Validator.String()).Msg("task already claimed, ensuring local cleanup")
					claimsToDelete = append(claimsToDelete, t.TaskId)
				} else {
					// Solution exists on-chain delete the solution and commitment
					solutionsToDelete = append(solutionsToDelete, t.TaskId)
					commitmentsToDelete = append(commitmentsToDelete, t.TaskId)

					// Solution is on-chain but NOT claimed. Update local status to claimable.
					tm.services.Logger.Info().Str("taskid", t.TaskId.String()).Str("validator", res.Validator.String()).Uint64("blocktime", res.Blocktime).Msg("task solution found and not claimed, updating status to claimable")
					claimTime := time.Unix(int64(res.Blocktime), 0)
					err = tm.services.TaskStorage.UpsertTaskToClaimable(t.TaskId, common.Hash{}, claimTime)
					if err != nil {
						tm.services.Logger.Error().Err(err).Str("taskid", t.TaskId.String()).Msg("failed to update task status to claimable")
						continue
					}
				}
				if res.Validator.String() != validator.ValidatorAddress().String() {
					tm.services.Logger.Warn().Msgf("solution solved by another validator! solver: %s task: %s", res.Validator.String(), t.TaskId.String())
				}
			} else {
				validTasksForBatch = append(validTasksForBatch, t)
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
			return nil // Nothing to process
		}

		if batchCommitmentLen < minBatchSize {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("commitments", batchCommitmentLen).Msg("available commitments below min batch size")
			return nil
		}

		tm.services.Logger.Info().Str("account", account.Address.String()).Msgf("bulk submitting %d commitment(s)", len(batchCommitments))

		isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode

		// Calculate hardcoded limit first
		hardcodedGasLimit := uint64(28_500*batchCommitmentLen + 3_500_000)
		// NonceManagerWrapper call
		tx, err := account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			if isEstimationMode {
				opts.GasLimit = 0 // Force estimation
			} else {
				opts.GasLimit = hardcodedGasLimit // Use hardcoded limit
			}
			opts.NoSend = true // don't send the transaction yet
			tx, err := tm.services.BulkTasks.BulkSignalCommitment(opts, batchCommitments)
			if err != nil {
				return nil, err
			}
			// send the transaction from correct account
			return account.SendSignedTransaction(tx)
		})

		// Log gas limit info
		if tx != nil {
			actualGasLimit := tx.Gas()
			logFields := map[string]interface{}{
				"tx_gas_limit": actualGasLimit,
				"batch_size":   batchCommitmentLen,
				"function":     "sendCommitments",
			}

			if isEstimationMode {
				logFields["mode"] = "estimate"
				logFields["hardcoded_gas_limit"] = hardcodedGasLimit
				if hardcodedGasLimit > 0 {
					diff := (float64(actualGasLimit) - float64(hardcodedGasLimit)) / float64(hardcodedGasLimit) * 100.0
					logFields["diff_percent"] = fmt.Sprintf("%.2f%%", diff)
				} else {
					logFields["diff_percent"] = "N/A (hardcoded is 0)"
				}
				tm.services.Logger.Info().Fields(logFields).Msg("gas limit estimation complete")
			} else {
				logFields["hardcoded_gas_limit"] = hardcodedGasLimit
				logFields["mode"] = "normal"
				tm.services.Logger.Info().Fields(logFields).Msg("gas limit set for transaction")
			}
		}

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error sending batch commitment")
			return err
		}
		if tx == nil {
			return errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		}

		// --- Normal logic continues below (WaitForConfirmedTx, etc.) ---
		receipt, success, _, waitErr := account.WaitForConfirmedTx(tx)
		// Capture waitErr separately from the wrapper err

		txCost := new(big.Int)

		// Process receipt even if WaitMined failed or tx reverted, if receipt exists
		if receipt != nil {
			gasCostPerTask := 0.0
			if batchCommitmentLen > 0 {
				gasCostPerTask = float64(receipt.GasUsed) / float64(batchCommitmentLen)
			}
			txCost.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			gp := tm.services.Config.BaseConfig.BaseToken.ToFloat(receipt.EffectiveGasPrice)
			tm.services.Logger.Info().Uint64("gas_used", receipt.GasUsed).Float64("gas_per_commit", gasCostPerTask).Float64("gas_price", gp).Msg("**** bulk commitment gas used *****")
			tm.cumulativeGasUsed.AddCommitment(txCost) // Add cost regardless of success
		}

		// Handle error from waiting for confirmation
		if waitErr != nil {
			tm.services.Logger.Error().Err(waitErr).Str("txhash", tx.Hash().String()).Msg("Error waiting for commitment confirmation")
			return waitErr // Return the wait error
		}

		if !success {
			// Transaction was mined but reverted
			tm.services.Logger.Error().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments tx reverted")
			return errors.New("batch commitments tx reverted")
		}

		tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments tx accepted!")

		var signalledCommitments [][32]byte
		for _, log := range receipt.Logs {
			if len(log.Topics) > 0 && log.Topics[0] == tm.signalCommitmentEvent {
				parsed, parseErr := tm.services.Engine.Engine.ParseSignalCommitment(*log)
				if parseErr != nil {
					tm.services.Logger.Error().Err(parseErr).Msg("could not parse signal commitment event")
					continue
				}
				signalledCommitments = append(signalledCommitments, parsed.Commitment)
			}
		}

		var commitmentsToUpdateAndDelete []task.TaskId

		for _, commitment := range batchCommitments {
			commitmentStr := task.TaskId(commitment).String()
			if slices.Contains(signalledCommitments, commitment) {
				if taskId, found := commitmentsToTaskMap[commitment]; found {
					commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
				}
			} else {
				blockNo, chainErr := tm.services.Engine.Engine.Commitments(nil, commitment)
				if chainErr != nil {
					tm.services.Logger.Error().Err(chainErr).Str("commitment", commitmentStr).Msg("could not get commitment block number")
					continue
				}
				// if we have a non zero block no, there is already a commitment for this task, so delete it
				if blockNo.Cmp(utils.Zero) != 0 {
					if taskId, found := commitmentsToTaskMap[commitment]; found {
						// Found on chain, but not in *this* batch's logs. Still needs cleanup/status update.
						commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
					}
					tm.services.Logger.Warn().Str("commitment", commitmentStr).Msg("commitment was already accepted (found on-chain)")
				} else {
					tm.services.Logger.Warn().Str("commitment", commitmentStr).Msg("commitment was not accepted (not in logs, not on-chain)")
				}
			}
		}

		// Delete local commitments now that they are confirmed on-chain (or found to be already confirmed)
		if err := deleteCommitments(commitmentsToUpdateAndDelete); err != nil {
			return err
		}

		unacceptedCommitment := len(batchCommitments) - len(signalledCommitments)
		if unacceptedCommitment > 0 {
			tm.services.Logger.Warn().Int("unaccepted", unacceptedCommitment).Msg("⚠️ commitments not accepted ⚠️")
		}

		if len(commitmentsToUpdateAndDelete) > 0 { // Check if there are tasks to update status for
			costPerCommitment := 0.0
			if len(commitmentsToUpdateAndDelete) > 0 && txCost.Cmp(big.NewInt(0)) > 0 {
				costPerCommitment = tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(commitmentsToUpdateAndDelete))
			}
			updateErr := tm.services.TaskStorage.UpdateTaskStatusAndCost(commitmentsToUpdateAndDelete, 2, costPerCommitment)
			if updateErr != nil {
				tm.services.Logger.Error().Err(updateErr).Msg("error updating task data in storage")
				return updateErr
			}
			tm.services.Logger.Info().Int("accepted", len(commitmentsToUpdateAndDelete)).Float64("cost_per_commit", costPerCommitment).Msg("✅ submitted commitments")
		}
		return nil // Successful completion
	}

	sendSolutions := func(validatorToSendSubmits *Validator, batchTasks storage.TaskDataSlice, wg *sync.WaitGroup, noChecks bool, minBatchSize int) error {

		batchSolutions, batchTaskIds := batchTasks.GetSolutions()
		batchSolutionsLen := len(batchSolutions)

		if batchSolutionsLen == 0 {
			return nil // Nothing to process
		}

		if batchSolutionsLen < minBatchSize {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("solutions", batchSolutionsLen).Msg("available solutions below min batch size")
			// No error, just not enough tasks to meet the minimum batch size
			return nil
		}

		var solutionsToSubmit [][]byte
		var tasksToSubmit [][32]byte
		for i, cid := range batchSolutions {
			taskid := batchTaskIds[i]

			solutionsToSubmit = append(solutionsToSubmit, cid)
			tasksToSubmit = append(tasksToSubmit, taskid)
		}

		tm.services.Logger.Info().Str("account", validatorToSendSubmits.ValidatorAddress().String()).Msgf("bulk submitting %d solution(s)", len(solutionsToSubmit))

		if len(solutionsToSubmit) != len(tasksToSubmit) {
			err := errors.New("ASSERT: MISMATCHED NUMBER OF TASKS AND SOLUTIONS BEING SUBMITTED!")
			tm.services.Logger.Error().Err(err).Msg("Internal error")
			return err
		}

		isEstimationMode := tm.services.Config.Miner.EnableGasEstimationMode

		// Calculate hardcoded limit first
		hardcodedGasLimit := uint64(139_500*batchSolutionsLen + 3_500_000)
		// NonceManagerWrapper call
		tx, err := validatorToSendSubmits.Account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			if isEstimationMode {
				opts.GasLimit = 0 // Force estimation
			} else {
				opts.GasLimit = hardcodedGasLimit // Use hardcoded limit
			}
			opts.NoSend = true // don't send the transaction yet
			tx, err := tm.services.Engine.Engine.BulkSubmitSolution(opts, tasksToSubmit, solutionsToSubmit)
			if err != nil {
				return nil, err
			}
			// send the transaction from correct account
			return validatorToSendSubmits.Account.SendSignedTransaction(tx)
		})

		if tx != nil {
			actualGasLimit := tx.Gas()
			logFields := map[string]interface{}{
				"tx_gas_limit": actualGasLimit,
				"batch_size":   batchSolutionsLen,
				"function":     "sendSolutions",
				"validator":    validatorToSendSubmits.ValidatorAddress().String(),
			}

			if isEstimationMode {
				logFields["mode"] = "estimate"
				logFields["hardcoded_gas_limit"] = hardcodedGasLimit
				if hardcodedGasLimit > 0 {
					diff := (float64(actualGasLimit) - float64(hardcodedGasLimit)) / float64(hardcodedGasLimit) * 100.0
					logFields["diff_percent"] = fmt.Sprintf("%.2f%%", diff)
				} else {
					logFields["diff_percent"] = "N/A (hardcoded is 0)"
				}
				tm.services.Logger.Info().Fields(logFields).Msg("gas limit estimation complete")
			} else {
				logFields["mode"] = "forced"
				tm.services.Logger.Info().Fields(logFields).Msg("gas limit set for transaction")
			}
		}

		if err != nil {
			tm.services.Logger.Error().Err(err).Int("batch_size", batchSolutionsLen).Msg("error sending batch solutions transaction")
			return err
		}

		if tx == nil {
			return errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		}

		receipt, success, _, waitErr := validatorToSendSubmits.Account.WaitForConfirmedTx(tx)

		txCost := new(big.Int)

		// Process receipt even if WaitMined failed or tx reverted, if receipt exists
		if receipt != nil {
			txCost.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			gp := tm.services.Config.BaseConfig.BaseToken.ToFloat(receipt.EffectiveGasPrice)
			gasPerSol := 0.0
			if batchSolutionsLen > 0 {
				gasPerSol = float64(receipt.GasUsed) / float64(batchSolutionsLen)
			}
			tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Float64("gas_per_sol", gasPerSol).Float64("gas_price", gp).Msg("**** bulk solution gas usage *****")
			tm.cumulativeGasUsed.AddSolution(txCost) // Add cost regardless of success
		}

		if waitErr != nil {
			tm.services.Logger.Error().Err(waitErr).Str("txhash", tx.Hash().String()).Msg("error waiting for solution confirmation")
			return waitErr
		}

		if !success {
			tm.services.Logger.Error().Str("txHash", tx.Hash().String()).Msg("batch solution transaction reverted on-chain")
			return errors.New("batch solution transaction reverted on-chain")
		} else {
			tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch solutions transaction successfully confirmed!")
		}

		var solutionsSubmitted []task.TaskId
		for _, log := range receipt.Logs {
			if len(log.Topics) > 0 && log.Topics[0] == tm.solutionSubmittedEvent {
				parsed, parseErr := tm.services.Engine.Engine.ParseSolutionSubmitted(*log)
				if parseErr != nil {
					tm.services.Logger.Error().Err(parseErr).Msg("could not parse solution submitted event")
					continue
				}
				solutionsSubmitted = append(solutionsSubmitted, parsed.Task)
			}
		}

		solutionsToDelete := make([]task.TaskId, 0)

		for _, taskid := range tasksToSubmit {
			taskIdStr := task.TaskId(taskid).String()
			if slices.Contains(solutionsSubmitted, taskid) {
				tm.services.TaskTracker.TaskSucceeded()
				solutionsToDelete = append(solutionsToDelete, taskid)
			} else {
				tm.services.TaskTracker.TaskFailed()

				res, chainErr := tm.services.Engine.Engine.Solutions(nil, taskid)
				if chainErr != nil {
					tm.services.Logger.Err(chainErr).Str("taskid", taskIdStr).Msg("error getting solution information for task onchain")
					continue
				}
				if res.Blocktime > 0 {
					solutionsToDelete = append(solutionsToDelete, taskid)
					tm.services.Logger.Warn().Str("taskid", taskIdStr).Str("validator", res.Validator.String()).Msg("solution was already accepted (found on-chain)")
				} else {
					tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("solution was not accepted (not in logs, not on-chain)")
				}
			}
		}

		// do not handle error here as we want to continue processing the next batch
		deleteErr := deleteSolutions(solutionsToDelete)
		if deleteErr != nil {
			tm.services.Logger.Error().Err(deleteErr).Msg("error deleting solutions from storage")
		}

		unacceptedSolutions := len(tasksToSubmit) - len(solutionsSubmitted)
		if unacceptedSolutions > 0 {
			tm.services.Logger.Warn().Int("unaccepted", unacceptedSolutions).Int("total_in_batch", len(tasksToSubmit)).Msg("⚠️ some solutions submitted in batch were not confirmed via logs")
		}

		if len(solutionsSubmitted) < 0 {
			return nil
		}

		gasPerSolution := 0.0
		if len(solutionsSubmitted) > 0 && txCost.Cmp(big.NewInt(0)) > 0 {
			gasPerSolution = tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(solutionsSubmitted))
		}

		_, addClaimErr := tm.services.TaskStorage.AddTasksToClaim(solutionsSubmitted, gasPerSolution)
		if addClaimErr != nil {
			tm.services.Logger.Error().Err(addClaimErr).Msg("error adding tasks to claim in storage")
			return addClaimErr
		}
		tm.services.Logger.Info().Str("validator", validatorToSendSubmits.ValidatorAddress().String()).Int("accepted", len(solutionsSubmitted)).Float64("cost_per_sol", gasPerSolution).Msg("✅ submitted solutions added to claim storage")

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

func (tm *BatchTransactionManager) processBulkClaimFast(account *account.Account, tasks []storage.ClaimTask, batchSize int) {

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
		tm.services.Logger.Warn().Str("txhash", receipt.TxHash.String()).Msg("⚠️ bulk claim transaction failed/reverted")
		return
	}

	tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Int("tasks", len(taskIds)).Msg("✅ bulk claim completed! checking logs for claimed tasks...")

	tasksClaimed := make([]task.TaskId, 0)

	// TODO: move this to load time
	event, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error getting engine abi")
		return
	}

	// Get the event SolutionClaimed topic ID
	eventTopic := event.Events["SolutionClaimed"].ID
	rewardsPaidTopic := event.Events["RewardsPaid"].ID

	var totalValidatorReward *big.Int = big.NewInt(0)

	for _, log := range receipt.Logs {
		if len(log.Topics) > 0 {
			// Check for SolutionClaimed event
			if log.Topics[0] == eventTopic {
				parsed, err := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("could not parse solution claimed event")
					continue
				}
				tasksClaimed = append(tasksClaimed, task.TaskId(parsed.Task))
			}

			// Check for RewardsPaid event to sum validatorReward
			if log.Topics[0] == rewardsPaidTopic {
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

	/*deleteTaskKeys := utils.Map(tasks, func(s storage.ClaimTask) string {
		return s.Key
	})*/

	/*	err = tm.services.TaskStorage.DeleteClaims(deleteTaskKeys)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not batch delete tasks")
			return
		}*/
	/*for _, taskid := range tasks {
		taskIdStr := task.TaskId(taskid.ID).String()
		err := tm.services.TaskStorage.DeleteClaim(taskid.Key)
		if err != nil {
			tm.services.Logger.Error().Err(err).Str("task", taskIdStr).Msg("could not delete key")
			return
		}
	}*/

	if len(tasksClaimed) > 0 {
		tm.services.Logger.Info().Int("claimed", len(tasksClaimed)).Msg("✅ claimed tasks")
	}

	unclaimedTaskCount := len(tasks) - len(tasksClaimed)
	if unclaimedTaskCount > 0 {
		tm.services.Logger.Warn().Int("unclaimed", unclaimedTaskCount).Msg("⚠️ tasks that failed to be claimed")
	}
}

func (tm *BatchTransactionManager) processBulkClaim(account *account.Account, tasks []storage.ClaimTask, minbatchSize, maxbatchSize int) {

	// Get the CooldownTime for each validator we might be claiming for
	var validatorCooldownTimesMap = make(map[common.Address]uint64)
	for _, v := range tm.validators {
		cooldownTime, err := v.CooldownTime(tm.minClaimSolutionTime, tm.minContestationVotePeriodTime)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error calling MinContestationVotePeriodTime")
			return
		}
		validatorCooldownTimesMap[v.ValidatorAddress()] = cooldownTime
	}

	tasksToClaim := make([]storage.ClaimTask, 0)

	claimsToDelete := make([]task.TaskId, 0)

	for _, task := range tasks {
		taskStr := task.ID.String()

		if task.Time >= time.Now().Unix() {
			continue
		}

		if tm.services.Config.Claim.ValidateClaims {

			result, err := tm.services.Engine.CanTaskIdBeClaimed(task, validatorCooldownTimesMap)
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("taskid", taskStr).Msg("CanTaskIdBeClaimed returned error - not claimable")
				continue
			}
			// if result is false then we can continue to the next task as this one is not claimable
			if !result {
				claimsToDelete = append(claimsToDelete, task.ID)

				tm.services.Logger.Debug().Str("taskid", taskStr).Msg("CanTaskIdBeClaimed returned false - not claimable")
				continue
			}
		}
		tasksToClaim = append(tasksToClaim, task)
	}

	if len(claimsToDelete) > 0 {
		err := tm.services.TaskStorage.DeleteClaims(claimsToDelete)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error deleting claims")
			return
		}
	}

	if len(tasksToClaim) < minbatchSize {
		// nothing to claim
		tm.services.Logger.Info().Int("min_batch", minbatchSize).Int("max_batch", maxbatchSize).Int("claimable", len(tasksToClaim)).Msg("claimable less than min batch size")
		return
	}

	// order the tasks by time oldest first
	sort.Slice(tasksToClaim, func(i, j int) bool {
		return tasksToClaim[i].Time < tasksToClaim[j].Time
	})

	if len(tasksToClaim) > maxbatchSize {
		tasksToClaim = tasksToClaim[:maxbatchSize]
	}

	// Use the Map function to extract the TaskId from each struct
	taskIds := utils.Map(tasksToClaim, func(s storage.ClaimTask) [32]byte {
		return s.ID
	})

	receipt, err := tm.BulkClaimWithAccount(account, taskIds)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("❌ error submitting bulk claim")
		return //err
	}

	//tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Msg("bulk claim tx sent")

	// Check if receipt is nil before accessing Status
	if receipt == nil {
		tm.services.Logger.Error().Msg("❌ bulk claim returned nil receipt despite no error")
		return
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		tm.services.Logger.Warn().Str("txhash", receipt.TxHash.String()).Msg("⚠️ bulk claim transaction failed/reverted")
		return
	}

	tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Int("tasks", len(taskIds)).Msg("✅ bulk claim completed! checking logs for claimed tasks...")

	tasksClaimed := make([]task.TaskId, 0)

	// TODO: move this to load time
	event, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error getting engine abi")
		return
	}

	// Get the event SolutionClaimed topic ID
	eventTopic := event.Events["SolutionClaimed"].ID
	rewardsPaidTopic := event.Events["RewardsPaid"].ID

	var totalValidatorReward *big.Int = big.NewInt(0)

	for _, log := range receipt.Logs {
		if len(log.Topics) > 0 {
			// Check for SolutionClaimed event
			if log.Topics[0] == eventTopic {
				parsed, err := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("could not parse solution claimed event")
					continue
				}
				tasksClaimed = append(tasksClaimed, task.TaskId(parsed.Task))
			}

			// Check for RewardsPaid event to sum validatorReward
			if log.Topics[0] == rewardsPaidTopic {
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

	if totalValidatorReward.Cmp(big.NewInt(0)) > 0 {
		totalValidatorRewardInAius := tm.services.Config.BaseConfig.BaseToken.ToFloat(totalValidatorReward)
		tm.services.Logger.Info().Float64("validator_reward", totalValidatorRewardInAius).Msg("total validator reward for bulk claim")
	}

	claimsToDelete = make([]task.TaskId, 0)

	for _, taskid := range tasksToClaim {
		taskIdStr := task.TaskId(taskid.ID).String()
		claimed := slices.Contains(tasksClaimed, taskid.ID)

		if claimed {
			claimsToDelete = append(claimsToDelete, taskid.ID)
		} else {
			result, err := tm.services.Engine.CanTaskIdBeClaimed(taskid, validatorCooldownTimesMap)
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("task", taskIdStr).Msg("error in canTaskIdBeClaimed")
				continue
			}
			// if result is true, we can re-add the task to the claim list
			if !result {
				tm.services.Logger.Debug().Str("taskid", taskIdStr).Msg("task will be deleted as it is not claimable")
				claimsToDelete = append(claimsToDelete, taskid.ID)
			}
		}
	}

	if len(claimsToDelete) > 0 {
		tm.services.Logger.Info().Int("claimed", len(claimsToDelete)).Msg("deleting claimed tasks from storage")
		tm.services.TaskStorage.DeleteClaims(claimsToDelete)
	}

	if len(tasksClaimed) > 0 {
		tm.services.Logger.Info().Int("claimed", len(tasksClaimed)).Msg("✅ claimed tasks")
	}

	unclaimedTaskCount := len(tasksToClaim) - len(tasksClaimed)
	if unclaimedTaskCount > 0 {
		tm.services.Logger.Warn().Int("unclaimed", unclaimedTaskCount).Msg("⚠️ tasks that failed to be claimed")
	}
}

func (tm *BatchTransactionManager) BatchCommitments() error {
	tm.Lock()
	copyCommitments := make([][32]byte, len(tm.commitments))
	copy(copyCommitments, tm.commitments)
	tm.commitments = [][32]byte{}
	tm.Unlock()

	tm.services.Logger.Info().Msgf("bulk submitting %d commitment(s)", len(copyCommitments))

	tx, err := tm.services.SenderOwnerAccount.NonceManagerWrapper(5, 425, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = uint64(27_900*len(copyCommitments) + 1_500_000)
		return tm.services.BulkTasks.BulkSignalCommitment(opts, copyCommitments)
	})

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error sending batch commitments")
		return err
	}

	receipt, success, _, _ := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tx)

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

func (tm *BatchTransactionManager) BatchSolutions() error {

	tm.services.Logger.Info().Msgf("bulk submitting %d solution(s)", len(tm.solutions))

	var solutionsToSubmit [][]byte
	var tasksToSubmit [][32]byte

	tm.Lock()
	copySolutions := make([]BatchSolution, len(tm.solutions))
	copy(copySolutions, tm.solutions)
	tm.solutions = []BatchSolution{}
	tm.Unlock()

	for _, solution := range copySolutions {

		soluitionInfo, err := tm.services.Engine.Engine.Solutions(nil, solution.taskID)

		if err != nil {
			tm.services.Logger.Err(err).Msg("error getting solution information")
			return err
		}

		if soluitionInfo.Blocktime > 0 {
			tm.services.Logger.Warn().Msg("solution exists")
		} else {
			solutionsToSubmit = append(solutionsToSubmit, solution.cid)
			tasksToSubmit = append(tasksToSubmit, solution.taskID)
		}
	}

	if len(solutionsToSubmit) != len(tasksToSubmit) {
		tm.services.Logger.Error().Msg("ASSERT: MISMATCHED TASKS AND SOLUTIONS!")
		return errors.New("mismatched tasks and solutions")
	}

	if len(solutionsToSubmit) > 0 {

		tx, err := tm.services.SenderOwnerAccount.NonceManagerWrapper(5, 425, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
			opts.GasLimit = 0
			return tm.services.Engine.Engine.BulkSubmitSolution(opts, tasksToSubmit, solutionsToSubmit)
		})

		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error sending batch solutions")
			return err
		}

		receipt, success, _, _ := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tx)

		if receipt != nil {
			txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Float64("gas_per_sol", float64(receipt.GasUsed)/float64(len(solutionsToSubmit))).Msg("**** bulk Solution gas used *****")
			tm.cumulativeGasUsed.AddSolution(txCost)
		}

		if !success {
			tm.services.Logger.Error().Err(err).Msg("batch solution tx reverted")
			return err
		}

		tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch solutions tx completed!")

	}

	return nil
}

func (tm *BatchTransactionManager) SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error {
	if tm.services.Config.IPFS.IncentiveClaim {
		return tm.services.TaskStorage.AddIpfsCid(taskId, cid)
	}
	return nil
}

func (tm *BatchTransactionManager) SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error {
	switch tm.batchMode {
	case 2:
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task commitment to batch")
		tm.Lock()
		tm.commitments = append(tm.commitments, commitment)
		tm.Unlock()
		return nil
	case 1:
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
	case 0:
		return tm.singleSignalCommitment(taskId, commitment)
	default:
	}
	return nil
}

func (tm *BatchTransactionManager) SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error {
	tm.services.TaskTracker.Solved()
	switch tm.batchMode {
	case 2:
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task solution to batch")

		bs := BatchSolution{
			taskID: taskId,
			cid:    cid,
		}

		tm.Lock()
		tm.solutions = append(tm.solutions, bs)
		tm.Unlock()
	case 1:
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task solution to batch")
		err := tm.services.TaskStorage.AddSolution(validator, taskId, cid)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error adding solution to storage")
		}
		return err
	case 0:
		return tm.singleSubmitSolution(validator, taskId, cid)
	default:

	}
	return nil
}

func (tm *BatchTransactionManager) singleSignalCommitment(taskId task.TaskId, commitment [32]byte) error {

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

	tx, err := tm.services.SenderOwnerAccount.NonceManagerWrapper(5, 425, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
		// nonce := 0
		// if opts.Nonce != nil {
		// 	nonce = int(opts.Nonce.Int64())
		// }
		// m.services.Logger.Info().Int("nonce", nonce).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("NonceManagerWrapper [sending commitment]")
		opts.GasLimit = 200_000

		return tm.services.Engine.Engine.SignalCommitment(opts, commitment)
	})

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error signaling commitment")
		return err
	}

	elapsed := time.Since(start)
	tm.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("signal commitment tx sent")

	go func() {
		receipt, success, _, err := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tx)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error waiting for commitment confirmation")
			return
		}
		if !success {
			tm.services.Logger.Error().Msg("commitment tx reverted")
			return
		}
		if receipt != nil {
			txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Uint64("gas_per_commit", receipt.GasUsed).Msg("single commitment gas used")
			tm.cumulativeGasUsed.AddCommitment(txCost)
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
			tm.services.Logger.Error().Msg("ASSERT: NO COMMITMENTS SIGNALLLED!")
		}
	}()

	// wait a bit to hope commitment is mined before submitting solution
	duration := 2000 + rand.Intn(350)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	return nil
}

func (tm *BatchTransactionManager) singleSubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error {

	val := tm.validators.GetValidatorByAddress(validator)
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
		// avoid gas estimate
		opts.GasLimit = 400_000

		return tm.services.Engine.Engine.SubmitSolution(opts, taskId, cid)
	})
	elapsed := time.Since(start)

	if err != nil {
		tm.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Str("elapsed", elapsed.String()).Msg("❌ error submitting solution")
		return err
	}

	tm.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("solution tx sent")

	go func() {
		// find out who mined the soluition and log it
		defer func() {
			res, err := tm.services.Engine.GetSolution(taskId)
			if err != nil {
				tm.services.Logger.Err(err).Msg("error getting solution information")
				return
			}

			if res.Blocktime > 0 {

				if tm.IsAddressValidator(res.Validator) {
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
				tm.services.Logger.Info().Str("taskid", taskIdStr).Msg("solution not solved")
			}

		}()

		receipt, success, _, _ := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tx)

		if receipt != nil {
			txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Uint64("gas_per_solution", receipt.GasUsed).Msg("**** single solution gas used *****")
			tm.cumulativeGasUsed.AddSolution(txCost)
		}

		if !success {
			return //errors.New("error waiting for solution confirmation")
		}

		tm.services.Logger.Info().Str("taskid", taskIdStr).Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("✅ solution accepted!")

		claims := []task.TaskId{taskId}
		claimTime, err := tm.services.TaskStorage.AddTasksToClaim(claims, 0)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error adding claim in redis")
			return
		}
		tm.services.Logger.Info().Str("taskid", taskIdStr).Str("when", claimTime.String()).Msg("added taskid claim to storage")
	}()

	return nil
}

// func (tm *BatchTransactionManagerV3) ValidatorDeposit(depositAmount *big.Int) (*types.Transaction, error) {
// 	return tm.services.SenderOwnerAccount.NonceManagerWrapper(3, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
// 		return tm.services.Engine.Engine.ValidatorDeposit(opts, tm.ValidatorAddress(), depositAmount)
// 	})
// }

func (tm *BatchTransactionManager) IsAddressValidator(address common.Address) bool {
	tm.Lock()
	defer tm.Unlock()

	for _, v := range tm.validators {
		if v.ValidatorAddress() == address {
			return true
		}
	}

	return false
}

func (tm *BatchTransactionManager) GetNextValidatorAddress() common.Address {
	tm.Lock()
	defer tm.Unlock()

	validator := tm.validators[tm.validatorIndex]
	tm.validatorIndex++
	if tm.validatorIndex >= len(tm.validators) {
		tm.validatorIndex = 0
	}

	return validator.ValidatorAddress()
}

func (tm *BatchTransactionManager) BulkClaim(taskIds [][32]byte) (*types.Receipt, error) {
	return tm.BulkClaimWithAccount(tm.services.SenderOwnerAccount, taskIds)
}

func (tm *BatchTransactionManager) BulkClaimWithAccount(account *account.Account, taskIds [][32]byte) (*types.Receipt, error) {
	tm.cumulativeGasUsed.ClaimValue(len(taskIds))

	tx, err := account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		// we can't rely on the gas estimator for this one
		// 1 in 10 it undershoots the gas limit
		opts.GasLimit = uint64(2_500_000 + len(taskIds)*60_000)
		return tm.services.BulkTasks.ClaimSolutions(opts, taskIds)
	})

	if err != nil {
		return nil, err
	}

	receipt, _, _, err := account.WaitForConfirmedTx(tx)

	if receipt != nil {
		txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)

		tm.cumulativeGasUsed.AddClaim(txCost)
	}

	return receipt, err
}

func (tm *BatchTransactionManager) BulkTasks(account *account.Account, count int) (*types.Receipt, error) {

	sendTx := func(ctx context.Context, client *client.Client, opts *bind.TransactOpts) (*types.Receipt, error) {
		tx, err := account.NonceManagerWrapperWithContext(ctx, opts, tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
			// we can't rely on the gas estimator for this one
			// 1 in 10 it undershoots the gas limit
			opts.GasLimit = uint64(1_500_000 + uint64(count)*185_000)
			tasks := big.NewInt(int64(count))
			//if client == nil {
			return tm.services.Engine.Engine.BulkSubmitTask(opts, tm.services.AutoMineParams.Version, tm.services.AutoMineParams.Owner, tm.services.AutoMineParams.Model, tm.services.AutoMineParams.Fee, tm.services.AutoMineParams.Input, tasks)
			/*} else {
				opts.NoSend = true
				tx, err := tm.services.BulkTasks.SubmitMultipleTasksEncoded(opts, tasks, tm.encodedTaskData)
				if err != nil {
					return nil, err
				}

				return client.SendSignedTransaction(tx)
			}*/
		})

		if err != nil {
			return nil, err
		}

		receipt, success, _, err := account.WaitForConfirmedTx(tx)

		if receipt != nil {
			gasPerTask := float64(receipt.GasUsed) / float64(count)
			txCost := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			gp := tm.services.Config.BaseConfig.BaseToken.ToFloat(receipt.EffectiveGasPrice)

			tm.services.Logger.Info().Uint64("gas_used", receipt.GasUsed).Float64("gas_per_task", gasPerTask).Float64("gas_price", gp).Msg("**** bulk tasks gas used *****")
			tm.cumulativeGasUsed.AddTasks(txCost)

			if success {
				var submittedTasks []task.TaskId
				for _, log := range receipt.Logs {

					if len(log.Topics) > 0 && log.Topics[0] == tm.taskSubmittedEvent {
						parsed, err := tm.services.Engine.Engine.ParseTaskSubmitted(*log)
						if err != nil {
							tm.services.Logger.Error().Err(err).Msg("could not parse task submitted event")
							continue
						}
						submittedTasks = append(submittedTasks, task.TaskId(parsed.Id))
					}
				}
				missingTasks := int(count) - len(submittedTasks)
				if missingTasks > 0 {
					tm.services.Logger.Warn().Int("missing", missingTasks).Msg("⚠️ tasks created mismatch ⚠️")
				}

				if len(submittedTasks) > 0 {

					costPerTask := tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(submittedTasks))
					tm.services.Logger.Info().Int("accepted", len(submittedTasks)).Float64("cost_per_task", costPerTask).Msg("✅ tasks created")

					err := tm.services.TaskStorage.AddTasks(submittedTasks, tx.Hash(), costPerTask)
					if err != nil {
						tm.services.Logger.Error().Err(err).Msg("error adding tasks in storage")
						return receipt, err
					}
					tm.services.Logger.Info().Int("tasks", len(submittedTasks)).Msgf("added tasks to storage")
				}
			}
		}

		return receipt, err
	}

	type clientSendResult struct {
		Receipt *types.Receipt
		Err     error
	}

	if len(tm.services.Clients) > 0 {
		// Create a context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // Make sure all paths cancel the context to avoid context leak

		// Create a channel to receive the results
		results := make(chan clientSendResult, len(tm.services.Clients))

		// Start a goroutine for each client
		for _, c := range tm.services.Clients {
			gasPrice, gasFeeCap, gasFeeTip, _ := account.Client.GasPriceOracle(false)
			opts := account.GetOptsWithoutNonceInc(0, gasPrice, gasFeeCap, gasFeeTip)

			go func(client *client.Client) {
				receipt, err := sendTx(ctx, client, opts)
				select {
				case results <- clientSendResult{receipt, err}:
				case <-ctx.Done():
				}
			}(c)
		}
		account.IncNonce()

		// Wait for the first successful result
		for i := 0; i < len(tm.services.Clients); i++ {
			result := <-results
			if result.Err != nil {
				tm.services.Logger.Error().Err(result.Err).Msg("error sending transaction")
			} else if result.Receipt != nil {
				cancel() // Stop the remaining goroutines
				return result.Receipt, nil
			}
		}
		return nil, errors.New("all clients failed to send transaction")
	} else {
		receipt, err := sendTx(context.Background(), nil, nil)
		return receipt, err
	}

}

func (m *BatchTransactionManager) ProcessValidatorsStakes() {
	m.services.Logger.Info().Msg("🚀 performing validator stake and balance checks")
	// get the Eth balance
	bal, err := m.services.OwnerAccount.GetBalance()
	if err != nil {
		m.services.Logger.Err(err).Str("account", m.services.OwnerAccount.Address.String()).Msg("could not get eth balance on account")
		return
	}

	balAsFloat := m.services.Eth.ToFloat(bal)

	// check if the Ether balance is less than configured threshold
	if balAsFloat < m.services.Config.ValidatorConfig.EthLowThreshold {
		m.services.Logger.Warn().Float64("threshold", m.services.Config.ValidatorConfig.EthLowThreshold).Msg("⚠️ balance is below threshold")
	}

	for _, v := range m.validators {

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

func (m *BatchTransactionManager) InitiateValidatorWithdraw(validator common.Address, amount float64) error {
	val := m.validators.GetValidatorByAddress(validator)
	if val != nil {
		return val.InitiateValidatorWithdraw(amount)
	}
	return errors.New("validator not found")
}

func (m *BatchTransactionManager) ValidatorWithdraw(validator common.Address) error {
	val := m.validators.GetValidatorByAddress(validator)
	if val != nil {
		return val.ValidatorWithdraw()
	}
	return errors.New("validator not found")
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

func (m *BatchTransactionManager) CancelValidatorWithdraw(validator common.Address, count int64) error {
	val := m.validators.GetValidatorByAddress(validator)
	if val != nil {
		return val.CancelValidatorWithdraw(count)
	}
	return errors.New("validator not found")
}

func (m *BatchTransactionManager) VoteOnContestation(validator common.Address, taskId task.TaskId, yeah bool) error {
	val := m.validators.GetValidatorByAddress(validator)
	if val != nil {
		return val.VoteOnContestation(taskId, yeah)
	}
	return errors.New("validator not found")
}

func (m *BatchTransactionManager) SubmitContestation(validator common.Address, taskId task.TaskId) error {
	val := m.validators.GetValidatorByAddress(validator)
	if val != nil {
		return val.SubmitContestation(taskId)
	}
	return errors.New("validator not found")
}

func (tm *BatchTransactionManager) Start(appQuit context.Context) error {

	tm.wg.Add(1)
	go tm.cumulativeGasUsed.Start(appQuit, tm.wg)

	// if the validator stake / balance check is enabled, start the validator stake / balance check processor
	if tm.services.Config.ValidatorConfig.StakeCheck {
		stakeCheckInterval, err := time.ParseDuration(tm.services.Config.ValidatorConfig.StakeCheckInterval)
		if err != nil {
			return err
		}
		tm.ProcessValidatorsStakes()
		tm.wg.Add(1) // Added wg.Add(1) to match the Done in processValidatorStakePoller
		go tm.processValidatorStakePoller(appQuit, tm.wg, stakeCheckInterval)
	} else {
		tm.services.Logger.Warn().Msg("validator stake / balance checks are disabled!")
	}

	// if the ipfs incentive claim is enabled, start the ipfs claim processor
	if tm.services.Config.IPFS.IncentiveClaim {
		claimInterval, err := time.ParseDuration(tm.services.Config.IPFS.ClaimInterval)
		if err != nil {
			return err
		}

		tm.wg.Add(1)
		go tm.startIpfsClaimProcessor(appQuit, tm.wg, claimInterval)
	}

	// start the batch claim processor
	// 1: normal batch system
	switch tm.batchMode {
	case 1:
		if tm.services.Config.Miner.UsePolling {
			d, err := time.ParseDuration(tm.services.Config.Miner.PollingTime)
			if err != nil {
				return err
			}
			tm.wg.Add(1)
			go tm.processBatchPoller(appQuit, tm.wg, d)
		} else {
			// process batch on block trigger (this hits external apis harder and not recommended)
			tm.wg.Add(1)
			go tm.processBatchBlockTrigger(appQuit, tm.wg)
		}
	case 0, 2:
		if tm.services.Config.Claim.Enabled {
			tm.wg.Add(1)
			go tm.batchClaimPoller(appQuit, tm.wg, time.Duration(tm.services.Config.Claim.Delay)*time.Second)
		}
	default:
		// error
		return errors.New("invalid batch mode")
	}

	return nil
}

func NewBatchTransactionManager(services *Services, ctx context.Context, wg *sync.WaitGroup) (*BatchTransactionManager, error) {

	// TODO: move these out of here!
	engineAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	// Get the event SolutionClaimed topic ID
	signalCommitmentEvent := engineAbi.Events["SignalCommitment"].ID
	solutionSubmittedEvent := engineAbi.Events["SolutionSubmitted"].ID
	taskSubmittedEvent := engineAbi.Events["TaskSubmitted"].ID

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

			account, err := account.NewAccount(pk, services.SenderOwnerAccount.Client, ctx, services.Config.Blockchain.CacheNonce, services.Logger)
			if err != nil {
				return nil, err
			}
			account.UpdateNonce()

			accounts = append(accounts, account)
		}
	} else {
		accounts = append(accounts, services.SenderOwnerAccount)
	}

	ratelimitEth, err := services.Engine.Engine.SolutionRateLimit(nil)
	if err != nil {
		return nil, err
	}

	ratelimit := services.Eth.ToFloat(ratelimitEth)

	var validators []*Validator
	for _, pk := range services.Config.ValidatorConfig.PrivateKeys {
		va, err := NewValidator(services, ctx, pk, services.SenderOwnerAccount.Client, ratelimit)
		if err != nil {
			return nil, err
		}
		validators = append(validators, va)
	}

	cache := NewCache(time.Duration(120) * time.Second)

	btm := &BatchTransactionManager{
		services:                      services,
		cache:                         cache,
		cumulativeGasUsed:             cumulativeGasUsed,
		signalCommitmentEvent:         signalCommitmentEvent,
		solutionSubmittedEvent:        solutionSubmittedEvent,
		taskSubmittedEvent:            taskSubmittedEvent,
		batchMode:                     services.Config.Miner.BatchMode,
		encodedTaskData:               encodedData,
		taskAccounts:                  accounts,
		wg:                            wg,
		validatorIndex:                0,
		validators:                    validators,
		minClaimSolutionTime:          minClaimSolutionTimeBig.Uint64(),
		minContestationVotePeriodTime: minContestationVotePeriodTimeBig.Uint64(),
	}

	return btm, nil
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
func (tm *BatchTransactionManager) processIpfsClaimsWithAccount(account *account.Account, ctx context.Context) error {
	tm.Lock()
	defer tm.Unlock()

	// Get all Ipfs entries from the database that haven't been claimed (oldest first)
	// TODO: make batch size configurable
	ipfsEntries, err := tm.services.TaskStorage.GetIpfsCids(10)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("failed to get unclaimed ipfs entries")
		return err
	}

	if len(ipfsEntries) == 0 {
		tm.services.Logger.Info().Msg("no unclaimed ipfs entries to process")
		return nil
	}

	tm.services.Logger.Info().Int("count", len(ipfsEntries)).Msg("processing ipfs claim entries")

	oracleClient := tm.services.IpfsOracle

	gasPrice, gasFeeCap, gasFeeTip, _ := account.Client.GasPriceOracle(false)

	for _, entry := range ipfsEntries {
		// Check if we should stop processing
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		cidHex := fmt.Sprintf("%x", entry.Cid)
		taskIdStr := entry.TaskId.String()

		// Check if incentive is still available on-chain
		incentiveAmount, err := tm.services.ArbiusRouter.Incentives(nil, entry.TaskId)
		if err != nil {
			tm.services.Logger.Error().
				Err(err).
				Str("taskId", taskIdStr).
				Msg("failed to check incentive availability")
			continue
		}

		if incentiveAmount.Cmp(big.NewInt(0)) == 0 {
			tm.services.Logger.Info().
				Str("taskId", taskIdStr).
				Str("cid", "0x"+cidHex).
				Msg("incentive already claimed for this task")

			// Mark as claimed in database to avoid reprocessing
			if err := tm.services.TaskStorage.DeleteIpfsCid(entry.TaskId); err != nil {
				tm.services.Logger.Error().Err(err).Str("taskId", taskIdStr).Msg("failed to mark entry as claimed")
			}
			continue
		}

		// Check if within 60 second priority window
		isWithinPriorityWindow := time.Since(entry.Added) <= 60*time.Second

		// If outside 60 seconds and not our validator, may want to be more strategic
		if !isWithinPriorityWindow {
			tm.services.Logger.Info().
				Str("taskId", taskIdStr).
				Str("cid", "0x"+cidHex).
				Dur("age", time.Since(entry.Added)).
				Msg("outside 60 second priority window, still attempting claim")
		}

		tm.services.Logger.Debug().
			Str("taskId", taskIdStr).
			Str("cid", "0x"+cidHex).
			Msg("getting signatures for cid")

		// Get signatures from oracle
		signatures, err := oracleClient.GetSignaturesForCID(ctx, cidHex)
		if err != nil {
			tm.services.Logger.Error().
				Err(err).
				Str("taskId", taskIdStr).
				Str("cid", "0x"+cidHex).
				Msg("failed to get signatures for cid")
			continue
		}

		if len(signatures) == 0 {
			tm.services.Logger.Warn().
				Str("taskId", taskIdStr).
				Str("cid", "0x"+cidHex).
				Msg("no signatures returned for cid")
			continue
		}

		// Submit on-chain claim with signatures
		tm.services.Logger.Info().
			Str("taskId", taskIdStr).
			Str("cid", "0x"+cidHex).
			Int("signatureCount", len(signatures)).
			Bool("priorityWindow", isWithinPriorityWindow).
			Str("incentiveAmount", incentiveAmount.String()).
			Msg("submitting ipfs claim on-chain")

		// convert signatures to arbiusrouterv1.Signature
		arbiusSignatures := make([]arbiusrouterv1.Signature, len(signatures))
		for i, signature := range signatures {
			arbiusSignatures[i] = arbiusrouterv1.Signature{
				Signer:    signature.Signer,
				Signature: common.FromHex(signature.Signature),
			}
		}

		opts := account.GetOptsWithoutNonceInc(0, gasPrice, gasFeeCap, gasFeeTip)
		tx, err := tm.services.ArbiusRouter.ClaimIncentive(opts, entry.TaskId, arbiusSignatures)
		if err != nil {
			tm.services.Logger.Error().Err(err).Str("taskId", taskIdStr).Msg("failed to submit ipfs claim on-chain")
			continue
		}

		var success bool

		// special case for mock router contract
		if *tx.To() == (common.Address{}) {
			tm.services.Logger.Info().Str("taskId", taskIdStr).Msg("mock router contract, skipping confirmation check")
			success = true
		} else {

			var receipt *types.Receipt

			receipt, success, _, err = account.WaitForConfirmedTx(tx)
			if err != nil {
				tm.services.Logger.Error().Err(err).Str("taskId", taskIdStr).Msg("failed to wait for ipfs claim on-chain")
				continue
			}

			if receipt != nil {
				tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("ipfs claim tx completed!")
			}
		}

		if success {
			// delete from database
			if err := tm.services.TaskStorage.DeleteIpfsCid(entry.TaskId); err != nil {
				tm.services.Logger.Error().Err(err).Str("taskId", taskIdStr).Msg("failed to mark entry as claimed")
			}
		} else {
			tm.services.Logger.Error().
				Str("taskId", taskIdStr).
				Str("cid", "0x"+cidHex).
				Msg("failed to submit ipfs claim on-chain")
		}
	}

	return nil
}

// StartIpfsClaimProcessor starts a background process to periodically check for and process IPFS claims
func (tm *BatchTransactionManager) startIpfsClaimProcessor(appQuit context.Context, wg *sync.WaitGroup, claimInterval time.Duration) {
	defer wg.Done()
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
			if err := tm.processIpfsClaimsWithAccount(tm.services.SenderOwnerAccount, appQuit); err != nil {
				if err != context.Canceled {
					tm.services.Logger.Error().Err(err).Msg("error processing ipfs claims")
				}
			}
		}
	}
}
