package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gobius/account"
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
	Value      interface{}
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

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found { //|| time.Since(item.LastUpdate) > c.ttl {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
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
	batchMode                     int
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

	if basefee == nil {
		basefee, err = tm.services.OwnerAccount.Client.GetBaseFee()
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not get basefee!")
			return 0, 0, 0, 0, 0, err
		}
	}

	basefeeinEth := tm.services.Eth.ToFloat(basefee)
	basefeeinGwei := basefeeinEth * 1000000000

	if tm.services.Config.BaseConfig.TestnetType > 0 {
		basePrice, ethPrice = 30, 2000
	} else {
		basePrice, ethPrice, err = tm.services.Paraswap.GetPrices()
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not get prices from oracle!")
			//return 0, 0, 0, 0, 0, err
			// TODO: for local offline testing we need to handle this!
			// if our oracle is offline for extended period of time the miner will be come inactive as it wont be able to process
			// various tasks
			basePrice, ethPrice = 30, 2000
		}
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

	totalReward, err := tm.services.Engine.GetModelReward(modelId)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("could not get model reward!")
		return 0, 0, 0, 0, 0, err
	}

	rewardInAIUS := tm.services.Config.BaseConfig.BaseToken.ToFloat(totalReward)

	tm.cache.Set("reward", rewardInAIUS)

	rewardInAIUSUSD := rewardInAIUS * basePrice
	rewardsPerBatchUSD := rewardInAIUSUSD * profitEstimateBatchSize

	profit := rewardsPerBatchUSD - totalCostPerBatchUSD

	tm.cumulativeGasUsed.profitEMA.Add(profit)

	tm.services.Logger.Info().
		Str("eth_in_usd", fmt.Sprintf("%.4g$", ethPrice)).
		Str("aius_in_usd", fmt.Sprintf("%.4g$", basePrice)).
		Msg("💰 eth/aius price")

	tm.services.Logger.Info().
		Str("costs_in_usd", fmt.Sprintf("%.4g$", totalCostPerBatchUSD)).
		Str("rewards_in_usd", fmt.Sprintf("%.4g$", rewardsPerBatchUSD)).
		Str("profit_per_batch", fmt.Sprintf("%.4g$", profit)).
		Str("base_fee", fmt.Sprintf("%.8g", basefeeinGwei)).
		Str("profit_metrics", tm.cumulativeGasUsed.profitEMA.String()).
		Msg("💰 batch profits")

	return profit, basefeeinGwei, rewardInAIUS, ethPrice, basePrice, nil
}

// This should only be run when not performing other batched operations and just does batch claims
func (tm *BatchTransactionManager) batchClaimPoller(appQuit context.Context, wg *sync.WaitGroup, pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-appQuit.Done():
			tm.services.Logger.Error().Msg("batch claimer shutting down")
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
			tm.services.Logger.Warn().Msg("delegated batch processor shutting down")
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
			tm.services.Logger.Warn().Msg("validator stake processor shutting down")
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
				return
			}

			start := time.Now()
			tm.processBatch(appQuit, &batchWG, profitLevel, baseFee, rewardInAIUS, ethPrice, basePrice)
			batchWG.Wait()
			tm.services.Logger.Warn().Str("duration", time.Since(start).String()).Msg("batch processed")
		case <-appQuit.Done():
			tm.services.Logger.Warn().Msg("delegated batch processor shutting down")
			return
		}
	}
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
	/*case "profitfeed":
	// use the profit feed to decide if we should process a batch
	basefeeBig, err := tm.services.OwnerAccount.Client.GetBaseFee()
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("could not get basefee!")
	}
	_, _, pl, _, pt, err := tm.services.ProfitFeed.PriceOracle(basefeeBig)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("price oracle error")
	} else if pt > 0 && pl > 0 {
		minProfit = pt
		profitLevel = pl
		tm.services.Logger.Info().Float64("price_target", pt).Float64("profit_level", pl).Msg("onchain oracle values")
	} else {
		tm.services.Logger.Warn().Msg("profit feed is not set, cannot use profit feed mode")
	}*/
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
	tm.services.Logger.Info().Bool("makeTasks", makeTasks).Int64("totalTasks", totalTasks).Int("minTasksInQueue", tm.services.Config.BatchTasks.MinTasksInQueue).Int("batchSize", tm.services.Config.BatchTasks.BatchSize).Msg("batch conditions")

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

			tm.services.Logger.Warn().Int("batch_size", taskBatchSize).Str("account", account.Address.String()).Msgf("** task queue is low - sending batch **")
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

	if !isProfitable {
		tm.services.Logger.Info().Str("profit_mode", profitMode).Str("min_profit", minProfitFmt).Str("max_profit", fmt.Sprintf("%.4g", maxProfit)).Msg("not profitable to process batch")
		return
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

	tm.services.Logger.Info().Int64("tasks", totalTasks).Int64("solutions", totalSolutions).Int64("commitments", totalCommitments).Int64("claims", totalClaims).Msg("pending totals")

	tm.services.Logger.Info().Str("profit_mode", profitMode).Str("min_profit", minProfitFmt).Msg("profit criteria met - processing batch")

	//validatorToSendSubmits := tm.validators[0]

	// validatorBuffer, err := validatorToSendSubmits.GetValidatorStakeBuffer()
	// if err != nil {
	// 	tm.services.Logger.Error().Err(err).Msg("could not get validator stake buffer")
	// 	return
	// }

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
		tm.services.Logger.Warn().Msgf("** **************************************** **")

		//}

		/*claims, err := tm.services.TaskStorage.GetClaims(claimMaxBatchSize * noOfClaimBatches)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("could not get keys from redis")
			return
		}

		totalCost := 0.0
		for _, task := range claims {
			totalCost += task.TotalCost
		}
		claimTasks := (47_300.0 / 1_000_000_000.0 * baseFee * float64(claimMaxBatchSize*noOfClaimBatches))

		tm.services.Logger.Warn().Msgf("** debug. total cost: %f  **", totalCost)

		totalCost += claimTasks

		totalCostInUSD := totalCost * ethPrice //fmt.Sprintf("%0.4f$", totalCost*ethPrice)

		claimValue := rewardInAIUS * float64(claimMaxBatchSize*noOfClaimBatches) * basePrice
		actualProfit := claimValue - totalCostInUSD

		tm.services.Logger.Warn().Msgf("**      total cost of mining batch : %0.4g$ (gas spent: %f)**", totalCostInUSD, totalCost)
		tm.services.Logger.Warn().Msgf("**                     batch value : %0.4g$ **", claimValue)
		tm.services.Logger.Warn().Msgf("**                          profit : %0.4g$ **", actualProfit)*/

		claimLen := len(claims)
		if claimLen > 0 {
			canClaim := true

			// claim on approach overrides everything else
			if rewardInAIUS >= tm.services.Config.Claim.ClaimMinReward {
				tm.services.Logger.Warn().Msgf("** reward is >= claim min reward, will claim **")
				canClaim = true
			} else if tm.services.Config.Claim.HoardMode && int(totalClaims) < tm.services.Config.Claim.HoardMaxQueueSize {
				canClaim = false
				tm.services.Logger.Warn().Msgf("** claim hoard mode on, and queue length below threshold - skipping claim **")
			} else if actualProfit < tm.services.Config.Claim.MinBatchProfit {
				tm.services.Logger.Warn().Msgf("** batch profit below claim threshold, skipping claim **")
				canClaim = false
			}

			// if tm.services.Config.Claim.ClaimOnApproachMinStake && validatorBuffer < tm.services.Config.Claim.MinStakeBufferLevel {
			// 	tm.services.Logger.Warn().Msgf("** claim on approach to min stake enabled and validator buffer is below set min stake level **")
			// 	canClaim = true
			// }

			if canClaim {
				// if tm.services.Config.DelegatedMiner.ConcurrentBatches {
				// 	wg.Add(1)
				// 	go func(wg *sync.WaitGroup) {
				// 		defer wg.Done()
				// 		tm.processBulkClaim(claims, claimMinBatchSize, claimMaxBatchSize)
				// 	}(&wg)
				// } else {
				// 	tm.processBulkClaim(claims, claimMinBatchSize, claimMaxBatchSize)
				// }

				accountIndex := rand.Intn(len(tm.taskAccounts))
				sendBulkClaim := func(chunk storage.ClaimTaskSlice, account *account.Account, wg *sync.WaitGroup, batchno int) {
					defer wg.Done()
					tm.services.Logger.Warn().Int("max_batch_size", claimMaxBatchSize).Int("batch_no", batchno+1).Str("address", account.Address.String()).Msgf("sending claim batch")
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
			tm.services.Logger.Warn().Msgf("deleted %d commitments from storage", len(_commitmentsToDelete))
		}
		return nil
	}

	deleteSolutions := func(_solutionsToDelete []task.TaskId) error {
		if len(_solutionsToDelete) > 0 {
			tm.services.Logger.Warn().Msgf("deleting %d solutions from storage", len(_solutionsToDelete))

			err := tm.services.TaskStorage.DeleteProcessedSolutions(_solutionsToDelete)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error deleting solution(s) from storage")
				return err
			}
		}
		return nil
	}

	getCommitmentBatchdata := func(batchSize int, noChecks bool) storage.TaskDataSlice {
		tm.Lock()
		defer tm.Unlock()
		var tasks storage.TaskDataSlice
		var err error

		tasks, err = tm.services.TaskStorage.GetPendingCommitments(batchSize)
		if err != nil {
			tm.services.Logger.Err(err).Msg("failed to get commitments from storage")
			return nil
		}

		if noChecks {
			return tasks
		}

		var commitmentsToDelete []task.TaskId
		var solutionsToDelete []task.TaskId

		for i, t := range tasks {

			block, err := tm.services.Engine.Engine.Commitments(nil, t.Commitment)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error getting commitment")
				continue
			}

			blockNo := block.Uint64()
			if blockNo > 0 {
				commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
				t.Commitment = [32]byte{}
			}

			res, err := tm.services.Engine.Engine.Solutions(nil, t.TaskId)

			if err != nil {
				tm.services.Logger.Err(err).Msg("error getting solution information")
				return nil
			}

			if res.Blocktime > 0 {

				// if res.Validator.String() != validatorToSendSubmits.ValidatorAddress().String() {
				// 	tm.services.Logger.Warn().Msgf("solution already exists for our task! solver: %s task: %s", res.Validator.String(), t.TaskId.String())
				// }
				solutionsToDelete = append(solutionsToDelete, t.TaskId)
				commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
				// Flag we need to delete both the commitment and the solution
				t.Commitment = [32]byte{}
				t.Solution = nil
			}

			// store changes to t back into slice
			tasks[i] = t
		}

		if err := deleteCommitments(commitmentsToDelete); err != nil {
			return nil
		}

		if err := deleteSolutions(solutionsToDelete); err != nil {
			return nil
		}

		return tasks
	}

	getSolutionBatchdata := func(batchSize int, noChecks bool) (*Validator, storage.TaskDataSlice, error) {
		tm.Lock()
		defer tm.Unlock()

		// map of validator to number of items we can send
		// loop through pending counts for each validator and do min(count, items we can send)
		// send the one with highest min
		solsPerVal, errTS := tm.services.TaskStorage.GetPendingSolutionsCountPerValidator()
		if errTS != nil {
			tm.services.Logger.Err(errTS).Msg("failed to get pending sols count per val from storage")
			return nil, nil, errTS
		}

		blockInfo, errbn := tm.services.OwnerAccount.Client.Client.BlockByNumber(context.Background(), nil)
		if errbn != nil {
			tm.services.Logger.Error().Err(errbn).Msg("Failed to get latest block")
			return nil, nil, errbn
		}

		blockTime := time.Unix(int64(blockInfo.Time()), 0)

		var validator *Validator = nil
		validatorHighestMin := int64(-1)
		for _, v := range tm.validators {

			maxSols := v.MaxSubmissions(blockTime)

			solsPending, found := solsPerVal[v.ValidatorAddress()]

			if found {
				minWeCanSend := min(maxSols, solsPending)

				tm.services.Logger.Info().Msgf("validator %s: max submissions: %d, sols pending: %d", v.ValidatorAddress().String(), maxSols, solsPending)

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
			tm.services.Logger.Info().Msg("no validator found to send solutions for")
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

		for i, t := range tasks {

			res, err := tm.services.Engine.Engine.Solutions(nil, t.TaskId)

			if err != nil {
				tm.services.Logger.Err(err).Msg("error getting solution information")
				return nil, nil, err
			}

			if res.Blocktime > 0 {

				if res.Validator.String() != validator.ValidatorAddress().String() {
					tm.services.Logger.Warn().Msgf("solution already exists for our task! solver: %s task: %s", res.Validator.String(), t.TaskId.String())
				}
				solutionsToDelete = append(solutionsToDelete, t.TaskId)
				commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
				// Flag we need to delete both the commitment and the solution
				t.Commitment = [32]byte{}
				t.Solution = nil
			}

			// store changes to t back into slice
			tasks[i] = t
		}

		if err := deleteCommitments(commitmentsToDelete); err != nil {
			return nil, nil, err
		}

		if err := deleteSolutions(solutionsToDelete); err != nil {
			return nil, nil, err
		}

		return validator, tasks, nil
	}

	sendCommitments := func(batchTasks storage.TaskDataSlice, account *account.Account, wg *sync.WaitGroup, noChecks bool, minBatchSize int) error {
		defer wg.Done()

		batchCommitments, commitmentsToTaskMap := batchTasks.GetCommitments()

		batchCommitmentLen := len(batchCommitments)

		if batchCommitmentLen >= minBatchSize {

			tm.services.Logger.Info().Str("account", account.Address.String()).Msgf("bulk submitting %d commitment(s)", len(batchCommitments))

			tx, err := account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
				opts.GasLimit = uint64(28_500*len(batchCommitments) + 3_500_000)
				return tm.services.BulkTasks.BulkSignalCommitment(opts, batchCommitments)
			})

			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("error sending batch commitment")
				return err
			}

			receipt, success, _, _ := account.WaitForConfirmedTx(tm.services.Logger, tx)

			txCost := new(big.Int)

			if receipt != nil {
				gasCostPerTask := float64(receipt.GasUsed) / float64(len(batchCommitments))
				txCost.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
				gp := tm.services.Config.BaseConfig.BaseToken.ToFloat(receipt.EffectiveGasPrice)
				tm.services.Logger.Info().Uint64("gas_used", receipt.GasUsed).Float64("gas_per_commit", gasCostPerTask).Float64("gas_price", gp).Msg("**** bulk commitment gas used *****")
				tm.cumulativeGasUsed.AddCommitment(txCost)
			}

			if !success {
				tm.services.Logger.Error().Err(err).Msg("batch commitments tx reverted")
				return err
			}

			tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch commitments tx accepted!")

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

			var commitmentsToUpdateAndDelete []task.TaskId

			for _, commitment := range batchCommitments {
				commitmentStr := task.TaskId(commitment).String()
				if slices.Contains(signalledCommitments, commitment) {
					if taskId, found := commitmentsToTaskMap[commitment]; found {
						commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
					}
				} else {
					blockNo, err := tm.services.Engine.Engine.Commitments(nil, commitment)
					if err != nil {
						tm.services.Logger.Error().Err(err).Str("commitment", commitmentStr).Msg("could not get commitment block number")
					}
					// if we have a non zero block no, there is already a commitment for this task, so delete it
					if blockNo.Cmp(utils.Zero) != 0 {
						if taskId, found := commitmentsToTaskMap[commitment]; found {
							commitmentsToUpdateAndDelete = append(commitmentsToUpdateAndDelete, taskId)
						}
						tm.services.Logger.Warn().Str("commitment", commitmentStr).Msg("commitment was already accepted")
					} else {
						tm.services.Logger.Warn().Str("commitment", commitmentStr).Msg("commitment was not accepted")
					}
				}
			}

			if len(commitmentsToUpdateAndDelete) > 0 {
				if err := deleteCommitments(commitmentsToUpdateAndDelete); err != nil {
					return err
				}
			}
			// if len(commitmentsToDelete) > 0 {
			// 	err := tm.services.TaskStorage.DeleteProcessedCommitments(commitmentsToDelete)
			// 	if err != nil {
			// 		tm.services.Logger.Error().Err(err).Msg("error deleting processed commitments from storage")
			// 	}
			// }

			unacceptedCommitment := len(batchCommitments) - len(signalledCommitments)
			if unacceptedCommitment > 0 {
				tm.services.Logger.Warn().Int("unaccepted", unacceptedCommitment).Msg("⚠️ commitments not accepted ⚠️")
			}

			if len(signalledCommitments) > 0 {

				costPerCommitment := tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(signalledCommitments))

				err = tm.services.TaskStorage.UpdateTaskStatusAndCost(commitmentsToUpdateAndDelete, 2, costPerCommitment)
				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("error updating task data in storage")
					return err
				}
				tm.services.Logger.Info().Int("accepted", len(signalledCommitments)).Float64("cost_per_commit", costPerCommitment).Msg("✅ submitted commitments")

			}

		} else if batchCommitmentLen > 0 {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("commitments", batchCommitmentLen).Msg("available commitments below min batch size")
		}

		return nil
	}

	sendSolutions := func(validatorToSendSubmits *Validator, batchTasks storage.TaskDataSlice, wg *sync.WaitGroup, noChecks bool, minBatchSize int) error {

		batchSolutions, batchTaskIds := batchTasks.GetSolutions()
		batchSolutionsLen := len(batchSolutions)

		if batchSolutionsLen >= minBatchSize {

			//time.Sleep(time.Duration(450 * time.Millisecond))

			tm.services.Logger.Info().Str("account", validatorToSendSubmits.ValidatorAddress().String()).Msgf("bulk submitting %d solution(s)", len(batchSolutions))

			// loop through all the solutions and check we have confirmed commitments on chain
			// if we do, then we can submit the solution
			// if we dont, then we can delete the solution from the storage
			// if we have a commitment on chain, but no solution, then we can delete the commitment from the storage
			var solutionsToSubmit [][]byte
			var tasksToSubmit [][32]byte
			var solutionsToDelete []task.TaskId

			for i, cid := range batchSolutions {
				taskid := batchTaskIds[i]

				if !noChecks {
					if !tm.services.Config.Miner.ConcurrentBatches {

						commitment, err := utils.GenerateCommitment(validatorToSendSubmits.ValidatorAddress(), taskid, cid)
						if err != nil {
							tm.services.Logger.Error().Err(err).Msg("error generating commitment")
							return err
						}

						res, err := tm.services.Engine.Engine.Commitments(nil, commitment)
						if err != nil {
							tm.services.Logger.Err(err).Msg("error getting solution information")
							return err
						}

						if res.Cmp(utils.Zero) == 0 {
							tm.services.Logger.Warn().Msg("commitment does not yet exist for solution")
							continue
						}
					}

					soluitionInfo, err := tm.services.Engine.Engine.Solutions(nil, taskid)

					if err != nil {
						tm.services.Logger.Err(err).Msg("error getting solution information")
						return err
					}

					if soluitionInfo.Blocktime > 0 {
						if soluitionInfo.Validator.String() != validatorToSendSubmits.ValidatorAddress().String() {
							tm.services.Logger.Warn().Msgf("solution already exists for our task! solver: %s task: %s", soluitionInfo.Validator.String(), taskid.String())
						}
						solutionsToDelete = append(solutionsToDelete, taskid)
					} else {
						solutionsToSubmit = append(solutionsToSubmit, cid)
						tasksToSubmit = append(tasksToSubmit, taskid)
					}
				} else {
					solutionsToSubmit = append(solutionsToSubmit, cid)
					tasksToSubmit = append(tasksToSubmit, taskid)
				}
			}

			if len(solutionsToSubmit) != len(tasksToSubmit) {
				tm.services.Logger.Error().Msg("ASSERT: MISMATCHED TASKS AND SOLUTIONS!")
				return err
			}

			if err := deleteSolutions(solutionsToDelete); err != nil {
				return err
			}

			if len(solutionsToSubmit) > 0 {

				// ratelimit, err := tm.services.Engine.Engine.SolutionRateLimit(nil)
				// if err == nil && ratelimit != nil {
				// 	tm.services.Logger.Info().Str("ratelimit", ratelimit.String()).Msg("solution ratelimit")
				// }

				// lastSubmissionBig, err := tm.services.Engine.Engine.LastSolutionSubmission(nil, validatorToSendSubmits.ValidatorAddress())
				// if err != nil {
				// 	tm.services.Logger.Error().Err(err).Msg("LastSolutionSubmission error")
				// } else {
				// 	tm.services.Logger.Info().Str("lastSubmission", lastSubmissionBig.String()).Msg("solution lastSubmission")
				// }

				// blockInfo, err := validatorToSendSubmits.Account.Client.Client.BlockByNumber(context.Background(), nil)
				// if err != nil {
				// 	tm.services.Logger.Error().Err(err).Msg("Failed to get latest block")
				// 	return err
				// }

				// lastSubmission := time.Unix(lastSubmissionBig.Int64(), 0)
				// blockTime := time.Unix(int64(blockInfo.Time()), 0)

				// diff := blockTime.Sub(lastSubmission)
				// secondsSinceSubmission := diff.Seconds()

				// //secondsSinceSubmission := time.Since(lastSubmission).Seconds()
				// tm.services.Logger.Info().Float64("secondsSinceSubmission", secondsSinceSubmission).Str("timelastSubmission", lastSubmission.Format(time.DateTime)).Msg("solution secondsSinceSubmission")

				// if float64(len(solutionsToSubmit)) > secondsSinceSubmission {
				// 	tm.services.Logger.Warn().Msg("rate limit")
				// 	return nil
				// }

				// block.timestamp - lastSolutionSubmission[msg.sender] >
				// solutionRateLimit * n_ / 1e18,

				tx, err := validatorToSendSubmits.Account.NonceManagerWrapper(tm.services.Config.Miner.ErrorMaxRetries, tm.services.Config.Miner.ErrorBackoffTime, tm.services.Config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
					// opts.GasLimit = 1_500_000
					// opts.NoSend = true
					// tx, _ := tm.delegatedvalidatorContract.BulkSubmitSolution(opts, solutiontasks, solutions)
					// if tx != nil {
					// 	hexStr := hex.EncodeToString(tx.Data())
					// 	tm.services.Logger.Warn().Msg(hexStr)
					// }
					opts.GasLimit = uint64(139_500*len(solutionsToSubmit) + 3_500_000)
					// opts.NoSend = false
					//return tm.services.Engine.Engine.BulkSubmitSolution(opts, tasksToSubmit, solutionsToSubmit)
					opts.NoSend = true
					tx, err := tm.services.Engine.Engine.BulkSubmitSolution(opts, tasksToSubmit, solutionsToSubmit)
					if err != nil {
						return nil, err
					}
					return validatorToSendSubmits.Account.SendSignedTransaction(tx)
				})

				if err != nil {
					tm.services.Logger.Error().Err(err).Msg("error sending batch solutions")
					return err
				}

				receipt, success, _, _ := validatorToSendSubmits.Account.WaitForConfirmedTx(tm.services.Logger, tx)

				txCost := new(big.Int)

				if receipt != nil {
					txCost.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
					gp := tm.services.Config.BaseConfig.BaseToken.ToFloat(receipt.EffectiveGasPrice)
					tm.services.Logger.Info().Uint64("gas", receipt.GasUsed).Float64("gas_per_sol", float64(receipt.GasUsed)/float64(len(solutionsToSubmit))).Float64("gas_price", gp).Msg("**** bulk solution gas used *****")
					tm.cumulativeGasUsed.AddSolution(txCost)
				}

				if !success {
					tm.services.Logger.Error().Err(err).Msg("batch solution tx reverted")
					return err
				}

				tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("batch solutions tx completed!")

				var solutionsSubmitted []task.TaskId
				for _, log := range receipt.Logs {
					if len(log.Topics) > 0 && log.Topics[0] == tm.solutionSubmittedEvent {
						parsed, err := tm.services.Engine.Engine.ParseSolutionSubmitted(*log)
						if err != nil {
							tm.services.Logger.Error().Err(err).Msg("could not parse solution submitted event")
							continue
						}
						solutionsSubmitted = append(solutionsSubmitted, parsed.Task)
					}
				}

				//	var solutionsToDelete []task.TaskId

				solutionsToDelete = make([]task.TaskId, 0)

				for _, taskid := range tasksToSubmit {
					taskIdStr := task.TaskId(taskid).String()
					if slices.Contains(solutionsSubmitted, taskid) {
						tm.services.TaskTracker.TaskSucceeded()
						solutionsToDelete = append(solutionsToDelete, taskid)
					} else {
						tm.services.TaskTracker.TaskFailed()

						res, err := tm.services.Engine.Engine.Solutions(nil, taskid)
						if err != nil {
							tm.services.Logger.Err(err).Msg("error getting solution information")
							continue
						}

						if res.Blocktime > 0 {
							solutionsToDelete = append(solutionsToDelete, taskid)

							tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("solution was already accepted")
						} else {
							tm.services.Logger.Warn().Str("taskid", taskIdStr).Msg("solution was not accepted")
						}
					}
				}

				gasPerSolution := tm.services.Config.BaseConfig.BaseToken.ToFloat(txCost) / float64(len(solutionsSubmitted))

				// do not handle error here as we want to continue processing the next batch
				deleteSolutions(solutionsToDelete)

				unacceptedSolutions := len(tasksToSubmit) - len(solutionsSubmitted)
				if unacceptedSolutions > 0 {
					tm.services.Logger.Warn().Int("unaccepted", unacceptedSolutions).Msg("⚠️ solutions not accepted ⚠️")
				}

				if len(solutionsSubmitted) > 0 {
					// Add the cost of this batch to all tasks
					// err = tm.services.TaskStorage.UpdateTaskStatusAndCost(solutionsSubmitted, 3, gasPerSolution)
					// if err != nil {
					// 	tm.services.Logger.Error().Err(err).Msg("error updating task data in storage")
					// 	return err
					// }
					_, err = tm.services.TaskStorage.AddTasksToClaim(solutionsSubmitted, gasPerSolution)
					if err != nil {
						tm.services.Logger.Error().Err(err).Msg("error adding tasks to claim in storage")
						return err
					}
					tm.services.Logger.Info().Str("validator", validatorToSendSubmits.ValidatorAddress().String()).Int("accepted", len(solutionsSubmitted)).Float64("cost_per_sol", gasPerSolution).Msg("✅ submitted solutions and added claims to storage")
					//tm.services.Logger.Info().Int("claims", len(solutionsSubmitted)).Msg("added claims to storage")
				}
			}
		} else if batchSolutionsLen > 0 {
			tm.services.Logger.Info().Int("min_batch_size", minBatchSize).Int("solutions", batchSolutionsLen).Msg("available solutions below min batch size")
		}

		return nil
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
			allBatchCommitments := getCommitmentBatchdata(tm.services.Config.Miner.CommitmentBatch.NumberOfBatches*maxBatchSize, tm.services.Config.Miner.NoChecks)

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
				// 2*90=180
				// 3*90=270
				endIndex := min(len(allBatchCommitments), (i+1)*maxBatchSize)
				batchTasks := allBatchCommitments[startIndex:endIndex]
				//panic: runtime error: slice bounds out of range [360:335]

				//tm.services.Logger.Info().Msgf("[   debug info   ] %d: start %d end %d len %d", i, startIndex, endIndex, len(allBatchCommitments))

				//if 0 == 1 {
				//sendCommitmentsAndSolutions(batchTasks, account, wg)
				//}
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
			if err != nil || batchTasks == nil || validator == nil {
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

	if receipt.Status != types.ReceiptStatusSuccessful {
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

	for _, log := range receipt.Logs {

		if len(log.Topics) > 0 && log.Topics[0] == eventTopic {
			parsed, err := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not parse solution claimed event")
				continue
			}
			tasksClaimed = append(tasksClaimed, task.TaskId(parsed.Task))
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
		//value, _ := task.ConvertTaskIdString2Bytes(s.ID)
		return s.ID
	})

	receipt, err := tm.BulkClaimWithAccount(account, taskIds)
	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("❌ error submitting bulk claim")
		return //err
	}

	//tm.services.Logger.Info().Str("txhash", receipt.TxHash.String()).Msg("bulk claim tx sent")

	if receipt.Status != types.ReceiptStatusSuccessful {
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

	for _, log := range receipt.Logs {

		if len(log.Topics) > 0 && log.Topics[0] == eventTopic {
			parsed, err := tm.services.Engine.Engine.ParseSolutionClaimed(*log)
			if err != nil {
				tm.services.Logger.Error().Err(err).Msg("could not parse solution claimed event")
				continue
			}
			tasksClaimed = append(tasksClaimed, task.TaskId(parsed.Task))
		}
	}

	if len(tasksClaimed) == 0 {
		tm.services.Logger.Warn().Msg("⚠️ successful bulk claim transaction but no tasks were claimed: check for contestations/cooldown issues ⚠️")
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
		tm.services.Logger.Info().Int("claimed", len(claimsToDelete)).Msg("deleting claims from storage")
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

	receipt, success, _, _ := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tm.services.Logger, tx)

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

		receipt, success, _, _ := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tm.services.Logger, tx)

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

func (tm *BatchTransactionManager) SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error {
	switch tm.batchMode {
	case 2:
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task commitment to batch")
		tm.Lock()
		tm.commitments = append(tm.commitments, commitment)
		tm.Unlock()
	case 1:
		tm.services.Logger.Debug().Str("taskid", taskId.String()).Msg("adding task commitment to batch")
		err := tm.services.TaskStorage.AddCommitment(validator, taskId, commitment)
		if err != nil {
			tm.services.Logger.Error().Err(err).Msg("error adding commitment to storage")
		}

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
		return tm.singleSubmitSolution(taskId, cid)
	default:

	}
	return nil
}

func (m *BatchTransactionManager) singleSignalCommitment(taskId task.TaskId, commitment [32]byte) error {

	taskIdStr := taskId.String()

	m.services.Logger.Info().Str("taskid", taskIdStr).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("sending commitment")

	start := time.Now()

	tx, err := m.services.SenderOwnerAccount.NonceManagerWrapper(5, 425, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
		// nonce := 0
		// if opts.Nonce != nil {
		// 	nonce = int(opts.Nonce.Int64())
		// }
		// m.services.Logger.Info().Int("nonce", nonce).Str("commitment", "0x"+hex.EncodeToString(commitment[:])).Msg("NonceManagerWrapper [sending commitment]")
		opts.GasLimit = 200_000

		return m.services.Engine.Engine.SignalCommitment(opts, commitment)
	})

	if err != nil {
		m.services.Logger.Error().Err(err).Msg("error signaling commitment")
		return err
	}

	elapsed := time.Since(start)
	m.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("signal commitment tx sent")

	go func() {
		receipt, _, _, _ := m.services.SenderOwnerAccount.WaitForConfirmedTx(m.services.Logger, tx)
		if receipt != nil {
			txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			m.services.Logger.Info().Uint64("gas", receipt.GasUsed).Uint64("gas_per_commit", receipt.GasUsed).Msg("**** single commitment gas used *****")
			m.cumulativeGasUsed.AddCommitment(txCost)
		}
	}()

	// wait a bit to hope commitment is mined
	duration := 200 + rand.Intn(150)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	return nil
}

func (m *BatchTransactionManager) singleSubmitSolution(taskId task.TaskId, cid []byte) error {

	taskIdStr := taskId.String()
	m.services.Logger.Info().Str("taskid", taskIdStr).Msg("sending single solution")

	start := time.Now()
	tx, err := m.services.SenderOwnerAccount.NonceManagerWrapper(8, 250, 1.5, false, func(opts *bind.TransactOpts) (interface{}, error) {
		// nonce := 0
		// if opts.Nonce != nil {
		// 	nonce = int(opts.Nonce.Int64())
		// }
		// m.services.Logger.Info().Int("nonce", nonce).Str("taskId", taskId.String()).Msg("NonceManagerWrapper [sending solution]")

		opts.GasLimit = 400_000

		return m.services.Engine.Engine.SubmitSolution(opts, taskId, cid)
	})
	elapsed := time.Since(start)

	if err != nil {
		m.services.Logger.Error().Err(err).Str("taskid", taskIdStr).Str("elapsed", elapsed.String()).Msg("❌ error submitting solution")
		return err
	}

	m.services.Logger.Info().Str("taskid", taskIdStr).Uint64("nonce", tx.Nonce()).Str("txhash", tx.Hash().String()).Str("elapsed", elapsed.String()).Msg("solution tx sent")

	go func() {
		// find out who mined the soluition and log it
		defer func() {
			res, err := m.services.Engine.GetSolution(taskId)
			if err != nil {
				m.services.Logger.Err(err).Msg("error getting solution information")
				return
			}

			if res.Blocktime > 0 {

				if m.IsAddressValidator(res.Validator) {
					m.services.TaskTracker.TaskSucceeded()
				} else {
					m.services.TaskTracker.TaskFailed()
				}
				solversCid := common.Bytes2Hex(res.Cid[:])
				ourCid := common.Bytes2Hex(cid)
				if ourCid != solversCid {
					m.services.Logger.Warn().Msg("=======================================================================")
					m.services.Logger.Warn().Msg("  WARNING: our solution cid does not match the solvers cid!")
					m.services.Logger.Warn().Msg("  our cid: " + ourCid)
					m.services.Logger.Warn().Msg("  ther cid: " + solversCid)
					m.services.Logger.Warn().Str("validator", res.Validator.String()).Msg("  solvers address")
					m.services.Logger.Warn().Msg("========================================================================")
				}
				m.services.Logger.Info().Str("taskid", taskIdStr).Str("validator", res.Validator.String()).Str("Cid", solversCid).Msg("solution information")
			} else {
				m.services.Logger.Info().Str("taskid", taskIdStr).Msg("solution not solved")
			}

		}()

		receipt, success, _, _ := m.services.SenderOwnerAccount.WaitForConfirmedTx(m.services.Logger, tx)

		if receipt != nil {
			txCost := receipt.EffectiveGasPrice.Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
			m.services.Logger.Info().Uint64("gas", receipt.GasUsed).Uint64("gas_per_solution", receipt.GasUsed).Msg("**** single solution gas used *****")
			m.cumulativeGasUsed.AddSolution(txCost)
		}

		if !success {
			return //errors.New("error waiting for solution confirmation")
		}

		m.services.Logger.Info().Str("taskid", taskIdStr).Str("txhash", tx.Hash().String()).Uint64("block", receipt.BlockNumber.Uint64()).Msg("✅ solution accepted!")

		claims := []task.TaskId{taskId}
		claimTime, err := m.services.TaskStorage.AddTasksToClaim(claims, 0)
		if err != nil {
			m.services.Logger.Error().Err(err).Msg("error adding claim in redis")
			return
		}
		m.services.Logger.Info().Str("taskid", taskIdStr).Str("when", claimTime.String()).Msg("added taskid claim to storage")
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

	receipt, _, _, err := account.WaitForConfirmedTx(tm.services.Logger, tx)

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

		receipt, success, _, err := account.WaitForConfirmedTx(tm.services.Logger, tx)

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
						tm.services.Logger.Error().Err(err).Msg("error adding tasks to claim in storage")
						return receipt, err
					}
					tm.services.Logger.Info().Int("claims", len(submittedTasks)).Msgf("added tasks to storage")
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
	// get the baseTokenBalance
	baseTokenBalance, err := m.services.Basetoken.BalanceOf(nil, m.services.OwnerAccount.Address)
	if err != nil {
		m.services.Logger.Err(err).Msg("failed to get balance")
		return
	}

	ethAsFmt := fmt.Sprintf("%.8g Ξ", balAsFloat)
	baseAsFmt := fmt.Sprintf("%.8g %s", m.services.Config.BaseConfig.BaseToken.ToFloat(baseTokenBalance), m.services.Config.BaseConfig.BaseToken.Symbol)
	//m.services.Config.BaseConfig.BaseToken.Symbol, m.services.Config.BaseConfig.BaseToken.ToFloat(baseTokenBalance)
	m.services.Logger.Info().Str("eth_bal", ethAsFmt).Str("basetoken_bal", baseAsFmt).Msg("wallet balances")

	for _, v := range m.validators {
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
	tm.cumulativeGasUsed.Start()

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
	}

	return nil
}

func NewBatchTransactionManager(ctx context.Context, wg *sync.WaitGroup) (*BatchTransactionManager, error) {

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		panic("could not get services from context")
	}

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

			account, err := account.NewAccount(pk, services.SenderOwnerAccount.Client, ctx, services.Config.Blockchain.CacheNonce)
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
		va, err := NewValidator(ctx, pk, services.SenderOwnerAccount.Client, ratelimit)
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
