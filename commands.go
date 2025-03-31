package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/client"
	task "gobius/common"
	"gobius/metrics"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

func importUnsolvedTasks(filename string, logger *zerolog.Logger, ctx context.Context) error {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	logger.Info().Str("file", filename).Msg("importing tasks into task queue")

	file, err := os.Open(filename)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not open file")
	}
	// send file to json decoder

	decoder := json.NewDecoder(file)
	taskIdsToTxes := make(map[string]string)
	err = decoder.Decode(&taskIdsToTxes)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not decode file")
	}
	file.Close()

	totalItemstoProcess := len(taskIdsToTxes)
	logger.Info().Int("tasks", totalItemstoProcess).Msg("processing tasks from file")

	mapOfHashsByTasks := make(map[string][]task.TaskId)
	mapOfPendingSolTasks := make(map[task.TaskId]struct{})

	pendingSols, err := services.TaskStorage.GetAllSolutions()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get all pending solutions")
	}

	for _, v := range pendingSols {
		mapOfPendingSolTasks[v.TaskId] = struct{}{}
	}

	allTasks, err := services.TaskStorage.GetQueuedTasks()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get all tasks")
	}

	mapOfTaskOnwers := make(map[common.Address]int)

	uniqueTasksMap := map[task.TaskId]struct{}{}
	for _, v := range allTasks {
		uniqueTasksMap[v.TaskId] = struct{}{}
	}

	whitelistedCount := 0
	solvedCount := 0
	pendingCount := 0
	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond)
	s.Suffix = " processing tasks..."
	s.FinalMSG = "completed!\n"
	s.Start()
	index := 0
	for taskId, txHash := range taskIdsToTxes {
		index++

		id, err := task.ConvertTaskIdString2Bytes(taskId)
		if err != nil {
			logger.Error().Err(err).Str("task", taskId).Msg("could convert task")
			return err
		}

		if _, ok := mapOfPendingSolTasks[id]; ok {
			pendingCount++
			continue
		}

		if _, ok := uniqueTasksMap[id]; ok {
			logger.Warn().Str("task", taskId).Msg("task already exists in the storage tasks list")
			continue
		}

		taskInfo, err := services.Engine.Engine.Tasks(nil, id)
		if err != nil {
			logger.Error().Err(err).Str("task", taskId).Msg("error getting task information")
			return err
		}

		mapOfTaskOnwers[taskInfo.Owner] = mapOfTaskOnwers[taskInfo.Owner] + 1

		res, err := services.Engine.GetSolution(id)
		if err != nil {
			logger.Err(err).Msg("error getting solution information")
			return err
		}
		logger.Debug().Uint64("blocktime", res.Blocktime).Bool("claimed", res.Claimed).Str("validator", res.Validator.String()).Str("Cid", common.Bytes2Hex(res.Cid[:])).Msg("old tasks being added to queue")

		if res.Blocktime == 0 {
			// isCommitment, err := miner.services.TaskStorage.IsCommitment(id)
			// if err != nil {
			// 	logger.Err(err).Msg("failed to check if task has a commitment")
			// 	return err
			// }

			// if isCommitment {
			// 	continue
			// }

			mapOfHashsByTasks[txHash] = append(mapOfHashsByTasks[txHash], id)

		} else {
			solvedCount++
		}

		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d]\n", index, totalItemstoProcess)

		// if index >= (totalItemstoProcess*10)/100 {
		// 	break
		// }
	}

	addedTasksCount := 0
	for key, value := range mapOfHashsByTasks {
		key := common.HexToHash(key)
		services.TaskStorage.AddTasks(value, key, 0)
		addedTasksCount += len(value)
	}

	for owner, v := range mapOfTaskOnwers {
		logger.Debug().Int("tasks", v).Str("owner", owner.String()).Msg("tasks per owner")
	}

	logger.Info().Int("pending_sol", pendingCount).Int("added_tasks", addedTasksCount).Int("solved", solvedCount).Int("whitelisted", whitelistedCount).Msg("finished adding unsolved tasks to queue")

	return nil
}

func taskCheck(logger *zerolog.Logger, ctx context.Context) error {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	allTasks, err := services.TaskStorage.GetQueuedTasks()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get all tasks")
	}

	mapOfTaskOnwers := make(map[common.Address]int)
	mapOfSolutionVals := make(map[common.Address]int)

	totalItemstoProcess := len(allTasks)
	logger.Info().Int("tasks", totalItemstoProcess).Msg("checking tasks")

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond)
	s.Suffix = " processing tasks..."
	s.FinalMSG = "completed!\n"
	s.Start()
	for index, key := range allTasks {

		taskInfo, err := services.Engine.Engine.Tasks(nil, key.TaskId)
		if err != nil {
			logger.Error().Err(err).Str("task", key.TaskId.String()).Msg("error getting task information")
			return err
		}

		mapOfTaskOnwers[taskInfo.Owner] = mapOfTaskOnwers[taskInfo.Owner] + 1

		res, err := services.Engine.GetSolution(key.TaskId)
		if err != nil {
			logger.Err(err).Msg("error getting solution information")
			return err
		}
		logger.Debug().Uint64("blocktime", res.Blocktime).Bool("claimed", res.Claimed).Str("validator", res.Validator.String()).Str("Cid", common.Bytes2Hex(res.Cid[:])).Msg("old tasks being added to queue")

		if res.Blocktime != 0 {
			mapOfSolutionVals[res.Validator] = mapOfSolutionVals[res.Validator] + 1
		}

		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d]\n", index, totalItemstoProcess)

	}

	for owner, v := range mapOfTaskOnwers {
		logger.Info().Int("tasks", v).Str("owner", owner.String()).Msg("tasks per owner")
	}

	for owner, v := range mapOfSolutionVals {
		logger.Info().Int("tasks", v).Str("val", owner.String()).Msg("solved tasks per validator")
	}

	return nil
}

func importUnclaimedTasks(filename string, logger *zerolog.Logger, ctx context.Context) error {

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	logger.Info().Str("file", filename).Msg("importing unclaimed tasks into claims queue")

	file, err := os.Open(filename)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not open file")
	}
	defer file.Close()
	// send file to json decoder

	decoder := json.NewDecoder(file)
	taskIds := []string{}
	err = decoder.Decode(&taskIds)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not decode file")
	}

	logger.Info().Int("tasks", len(taskIds)).Msg("adding claims to storage")

	tasks := make([]task.TaskId, len(taskIds))
	for _, v := range taskIds {
		taskId, err := task.ConvertTaskIdString2Bytes(v)
		if err != nil {
			logger.Error().Err(err).Str("task", v).Msg("could convert task")
			break
		}

		tasks = append(tasks, taskId)

	}
	services.TaskStorage.AddTasksToClaim(tasks, 0)
	//logger.Info().Msg("claims added, now running deduplicate and verify stage")

	// TODO: readd
	// dedupeVerifyClaims(logger, ctx)

	return nil
}

func verifyClaims(logger *zerolog.Logger, ctx context.Context) {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	logger.Info().Msg("verifying claims in storage")

	_, claims, err := services.TaskStorage.TotalSolutionsAndClaims()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get claim totals from  storage")
	}

	allClaims, _, err := services.TaskStorage.GetClaims(int(claims))
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get claims in storage")
	}

	logger.Info().Int("claims", len(allClaims)).Msg("verifying claims")

	minClaimSolutionTimeBig, err := services.Engine.Engine.MinClaimSolutionTime(nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("error calling MinClaimSolutionTime")
	}

	minContestationVotePeriodTimeBig, err := services.Engine.Engine.MinContestationVotePeriodTime(nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("error calling MinContestationVotePeriodTime")
	}

	cacheValidatorCooldown := map[common.Address]uint64{}

	var claimsToDelete []task.TaskId

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond)
	s.Suffix = " processing tasks..."
	s.FinalMSG = "completed!\n"
	s.Start()
	totalItemstoProcess := len(allClaims)

	for index, value := range allClaims {

		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d]", index, totalItemstoProcess)

		taskStr := value.ID.String()

		contestationDetails, err := services.Engine.GetContestation(value.ID)
		if err != nil {
			logger.Fatal().Err(err).Str("task", taskStr).Msg("could not get contestation details")
		}

		if contestationDetails.Validator.String() != "0x0000000000000000000000000000000000000000" {
			contestor := contestationDetails.Validator.String()

			logger.Warn().Str("task", taskStr).Str("contestor", contestor).Str("validator", contestationDetails.Validator.String()).Str("slashedamount", contestationDetails.SlashAmount.String()).Msg("⚠️ task was contested, deleting ⚠️")

			claimsToDelete = append(claimsToDelete, value.ID)
			continue
		}

		solution, err := services.Engine.GetSolution(value.ID)
		if err != nil {
			logger.Fatal().Err(err).Str("task", taskStr).Msg("cloud not get solution details")
		}

		cooldownTime := uint64(0)
		//cacheValidatorCooldown[solution.Validator]
		if cooldownTime, ok = cacheValidatorCooldown[solution.Validator]; !ok {
			lastContestationLossTimeBig, err := services.Engine.Engine.LastContestationLossTime(nil, solution.Validator)
			if err != nil {
				logger.Fatal().Err(err).Msg("error calling LastContestationLossTime")
			}
			lastContestationLossTime := lastContestationLossTimeBig.Uint64()
			if lastContestationLossTime > 0 {
				minClaimSolutionTime := minClaimSolutionTimeBig.Uint64()
				minContestationVotePeriodTime := minContestationVotePeriodTimeBig.Uint64()
				cooldownTime = lastContestationLossTime + minClaimSolutionTime + minContestationVotePeriodTime
				cacheValidatorCooldown[solution.Validator] = cooldownTime
				logger.Debug().Uint64("lastcontestationlosttime", lastContestationLossTime).Uint64("cooldowntime", cooldownTime).Msg("last contestation time")
			} else {
				cacheValidatorCooldown[solution.Validator] = 0
			}
		}
		if solution.Blocktime <= cooldownTime {
			logger.Warn().Str("taskid", taskStr).Str("validator", solution.Validator.String()).Msg("⚠️ lost due to contestation cooldown ⚠️")
			claimsToDelete = append(claimsToDelete, value.ID)
			continue
		}

		if solution.Claimed {
			logger.Info().Str("task", taskStr).Msgf("task already claimed by %s", solution.Validator.String())
			claimsToDelete = append(claimsToDelete, value.ID)
			continue
		}
	}
	s.Stop()

	if len(claimsToDelete) > 0 {
		logger.Info().Int("claims", len(claimsToDelete)).Msgf("deleting claimed or unclaimable tasks")

		err := services.TaskStorage.DeleteClaims(claimsToDelete)
		if err != nil {
			logger.Fatal().Err(err).Msg("could not delete claims")
		}
	}
	logger.Info().Msgf("verified claims and %d deleted", len(claimsToDelete))
}

func getBatchPricingInfo(ctx context.Context) error {
	var err error
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	basePrice, ethPrice, err := services.Paraswap.GetPrices()
	if err != nil {
		services.Logger.Error().Err(err).Msg("could not get prices from oracle api!")
	}

	basefee, err := services.OwnerAccount.Client.GetBaseFee()
	if err != nil {
		services.Logger.Error().Err(err).Msg("could not get basefee!")
	}

	// convert basefee to gwei
	basefeeinEth := services.Eth.ToFloat(basefee)

	//rewardInAIUS := tm.cumulativeGasUsed.rewardEMA.Average()
	reward, err := services.Engine.Engine.GetReward(nil)
	if err != nil {
		services.Logger.Error().Err(err).Msg("could not get reward!")
	}

	rewardInAIUS := services.Config.BaseConfig.BaseToken.ToFloat(reward) * 0.9

	claimMaxBatchSize := services.Config.Claim.MaxClaims

	//claims, err := services.TaskStorage.GetClaims(claimMaxBatchSize)
	claims, averageGas, err := services.TaskStorage.GetClaims(claimMaxBatchSize)
	if err != nil {
		services.Logger.Error().Err(err).Msg("could not get keys from storage")
		return err
	}

	claimMaxBatchSize = len(claims)

	totalCost := 0.0
	for _, task := range claims {
		totalCost += task.TotalCost
	}
	claimTasks := (47_300.0 * basefeeinEth * float64(claimMaxBatchSize))

	services.Logger.Warn().Msgf("** debug. total cost       : %f  **", totalCost)
	services.Logger.Warn().Msgf("**        average gas/task : %f  **", averageGas)

	totalCost += claimTasks

	totalCostInUSD := totalCost * ethPrice //fmt.Sprintf("%0.4f$", totalCost*ethPrice)

	claimValue := rewardInAIUS * float64(claimMaxBatchSize) * basePrice
	services.Logger.Warn().Msgf("**      total cost of mining batch : %0.4g$ (gas spent: %f)**", totalCostInUSD, totalCost)
	services.Logger.Warn().Msgf("**                     batch value : %0.4g$ **", claimValue)
	services.Logger.Warn().Msgf("**                          profit : %0.4g$ **", claimValue-totalCostInUSD)

	return nil
}

func cleanTaskData(ctx context.Context) error {

	// Get the services from the context
	/*
		services, ok := ctx.Value(servicesKey{}).(*Services)
		if !ok {
			log.Fatal("Could not get services from context")
		}


			taskData, err := services.TaskStorage.GetAllTaskData()
			if err != nil {
				services.Logger.Fatal().Err(err).Msg("could not get claims in storage")
			}

			services.Logger.Info().Int("total", len(taskData)).Msg("cleaning task data")

			s := spinner.New(spinner.CharSets[11], 500*time.Millisecond)
			s.Suffix = " processing tasks..."
			s.FinalMSG = "completed!\n"
			s.Start()
			totalItemstoProcess := len(taskData)
			deleted := 0
			for index, v := range taskData {
				res, err := services.Engine.Engine.Solutions(nil, v)

				if err != nil {
					services.Logger.Fatal().Err(err).Msg("error getting solution information")
				}

				if res.Blocktime > 0 {
					err := services.TaskStorage.DeleteTaskaData(v.String())
					if err != nil {
						services.Logger.Fatal().Err(err).Msg("could delete task data key")
					}
					deleted++
				}

				s.Suffix = fmt.Sprintf(" processing tasks [%d/%d] [deleted: %d]\n", index, totalItemstoProcess, deleted)

			}

			s.Stop()


			services.Logger.Info().Int("deleted", deleted).Msg("completed")*/

	return nil
}

func verifySolutions(ctx context.Context) error {

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	deleteCommitments := func(_commitmentsToDelete []task.TaskId) error {
		const batchSize = 1000

		if len(_commitmentsToDelete) > 0 {
			/*err := services.TaskStorage.DeleteProcessedCommitments(_commitmentsToDelete)
			if err != nil {
				services.Logger.Error().Err(err).Msg("error deleting commitment(s) from storage")
				return err
			}
			services.Logger.Warn().Msgf("deleted %d commitments from storage", len(_commitmentsToDelete))
			*/
			for i := 0; i < len(_commitmentsToDelete); i += batchSize {
				end := i + batchSize
				if end > len(_commitmentsToDelete) {
					end = len(_commitmentsToDelete)
				}

				batch := _commitmentsToDelete[i:end]

				err := services.TaskStorage.DeleteProcessedCommitments(batch)
				if err != nil {
					services.Logger.Error().Err(err).Msg("error deleting commitment(s) from storage")
					return err
				}
				services.Logger.Warn().Msgf("deleted %d commitments from storage", len(batch))
			}
		}
		return nil
	}

	deleteSolutions := func(_solutionsToDelete []task.TaskId) error {
		const batchSize = 1000

		if len(_solutionsToDelete) > 0 {

			/*err := services.TaskStorage.DeleteProcessedSolutions(_solutionsToDelete)
			if err != nil {
				services.Logger.Error().Err(err).Msg("error deleting solution(s) from storage")
				return err
			}
			services.Logger.Warn().Msgf("deleted %d solutions from storage", len(_solutionsToDelete))
			*/

			for i := 0; i < len(_solutionsToDelete); i += batchSize {
				end := i + batchSize
				if end > len(_solutionsToDelete) {
					end = len(_solutionsToDelete)
				}

				batch := _solutionsToDelete[i:end]

				err := services.TaskStorage.DeleteProcessedSolutions(batch)
				if err != nil {
					services.Logger.Error().Err(err).Msg("error deleting solution(s) from storage")
					return err
				}
				services.Logger.Warn().Msgf("deleted %d solutions from storage", len(batch))
			}
		}
		return nil
	}

	tasks, err := services.TaskStorage.GetAllSolutions()
	if err != nil {
		services.Logger.Err(err).Msg("failed to get tasks from storage")
	}

	services.Logger.Info().Int("solutions", len(tasks)).Msg("verifying solutions")

	var commitmentsToDelete []task.TaskId
	var solutionsToDelete []task.TaskId
	solvedByMap := map[common.Address]int{}

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond)
	s.Suffix = " processing tasks..."
	s.FinalMSG = "completed!\n"
	s.Start()
	totalItemstoProcess := len(tasks)
	for index, t := range tasks {
		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d]", index, totalItemstoProcess)

		if t.Commitment != [32]byte{} {
			// commitStr := task.TaskId(t.Commitment).String()
			// // commitStr := task.

			// tm.services.Logger.Info().Msgf("bulk submitting commitment: %s ", commitStr)

			block, err := services.Engine.Engine.Commitments(nil, t.Commitment)
			if err != nil {
				services.Logger.Error().Err(err).Msg("error getting commitment")
				continue
			}

			blockNo := block.Uint64()
			if blockNo > 0 {
				commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
				t.Commitment = [32]byte{}
			}
		}

		res, err := services.Engine.Engine.Solutions(nil, t.TaskId)

		if err != nil {
			services.Logger.Err(err).Msg("error getting solution information")
			return nil
		}

		if res.Blocktime > 0 {
			solvedByMap[res.Validator] = solvedByMap[res.Validator] + 1
			solutionsToDelete = append(solutionsToDelete, t.TaskId)
			commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
			// Flag we need to delete both the commitment and the solution
			t.Commitment = [32]byte{}
			t.Solution = nil
		}
	}

	s.Suffix = " deleting commitments"
	if err := deleteCommitments(commitmentsToDelete); err != nil {
		return nil
	}

	s.Suffix = " deleting solutions"
	if err := deleteSolutions(solutionsToDelete); err != nil {
		return nil
	}

	s.Stop()

	for owner, v := range solvedByMap {
		services.Logger.Info().Int("tasks", v).Str("val", owner.String()).Msg("solved tasks per validator")
	}
	return nil

}

func blockMonitor(ctx context.Context, rpcClient *client.Client) error {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	headers := make(chan *types.Header)
	var newHeadSub ethereum.Subscription

	connectToHeaders := func() {
		var err error

		newHeadSub, err = rpcClient.Client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			services.Logger.Fatal().Err(err).Msg("Failed to subscribe to new headers")
		}
	}

	services.TaskTracker.Silence(true)

	connectToHeaders()

	maxBackoffHeader := time.Second * 30
	currentBackoffHeader := time.Second

	bm := metrics.NewBlockMetrics(rpcClient, services.Config, services.Engine.Engine)

	for {

		select {

		case h := <-headers:
			bm.UpdateBlockMetrics(h)

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

func depositMonitor(ctx context.Context, rpcClient *client.Client, startBlock, endBlock int64) {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	if endBlock <= 0 {
		// Get the current block number
		currentblock, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			log.Fatalf("Failed to get current block number: %v", err)
		}

		endBlock = int64(currentblock)
	}

	bm := metrics.NewBlockMetrics(rpcClient, services.Config, services.Engine.Engine)

	bm.ProcessDepositWithdrawLogs(rpcClient, startBlock, endBlock, services.Engine.Engine)

}

func fundTaskWallets(ctx context.Context, amount float64, minbal float64) {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	if len(services.Config.BatchTasks.PrivateKeys) <= 0 {
		log.Fatal("no task wallets defined in config")
	}

	if amount < 0 || amount > 0.5 {
		log.Fatal("transfer amount is out of range (must be > 0 and 0.5 or under)")
	}

	ethBalance, err := services.OwnerAccount.GetBalance()
	if err != nil {
		services.Logger.Err(err).Str("account", services.OwnerAccount.Address.String()).Msg("could not get eth balance on account")
		return
	}

	// convert ETH balance to float
	balAsFloat := services.Eth.ToFloat(ethBalance)

	if balAsFloat < amount*float64(len(services.Config.BatchTasks.PrivateKeys)) {
		services.Logger.Err(err).Str("balance", fmt.Sprintf("%.4g", balAsFloat)).Msg("not enough eth balance to satisfy transfer")
		return
	}

	amountAsBig := services.Eth.FromFloat(amount)

	// manually update the gas base fee
	/*err = services.SenderOwnerAccount.Client.UpdateCurrentBasefee()
	if err != nil {
		services.Logger.Error().Err(err).Msg("error in new account")
		return
	}*/

	for _, pk := range services.Config.BatchTasks.PrivateKeys {

		account, err := account.NewAccount(pk, services.SenderOwnerAccount.Client, ctx, services.Config.Blockchain.CacheNonce, services.Logger)
		if err != nil {
			services.Logger.Error().Err(err).Msg("error in new account")
			return
		}
		accountBal, err := account.GetBalance()
		if err != nil {
			services.Logger.Err(err).Str("account", account.Address.String()).Msg("could not get eth balance on account")
			return
		}

		// convert ETH balance to float
		balAsFloat := services.Eth.ToFloat(accountBal)

		if balAsFloat >= minbal {
			services.Logger.Info().Msgf("%s balance of %.4g eth is at or greater than minbal %.4g eth", account.Address.String(), balAsFloat, minbal)
			continue
		}

		services.SenderOwnerAccount.UpdateNonce()
		if err != nil {
			services.Logger.Error().Err(err).Msg("error updating nonce")
			return
		}

		services.Logger.Info().Msgf("💼 transfering %.4g eth to address %s (with balance of %.4g eth)", amount, account.Address.String(), balAsFloat)

		tx, err := services.SenderOwnerAccount.SendEther(nil, account.Address, amountAsBig)
		if err != nil {
			services.Logger.Error().Err(err).Msg("error sending transfer")
			return
		}
		_, success, _, err := services.SenderOwnerAccount.WaitForConfirmedTx(tx)

		if err != nil {
			services.Logger.Error().Err(err).Msg("error waiting for transfer")
			return
		}

		if success {
			services.Logger.Info().Str("txhash", tx.Hash().String()).Str("amount_transfered", fmt.Sprintf("%.4g", amount)).Msg("✅ eth transfered")
		}
	}

	services.Logger.Info().Msg("transfers complete")
}
