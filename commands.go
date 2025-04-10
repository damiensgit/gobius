package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/bindings/engine"
	"gobius/client"
	task "gobius/common"
	"gobius/metrics"
	"gobius/storage"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

func importUnsolvedTasks(appQuit context.Context, filename string, removeMode bool, logger *zerolog.Logger, ctx context.Context) error {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	if removeMode {
		logger.Info().Str("file", filename).Msg("removing tasks listed in file from task queue")
	} else {
		logger.Info().Str("file", filename).Msg("importing tasks into task queue")
	}

	file, err := os.Open(filename)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not open file")
	}
	// send file to json decoder
	decoder := json.NewDecoder(file)
	taskIdsToTxes := make(map[string]string) // Still need TxHash if importing, might be unused if removing
	err = decoder.Decode(&taskIdsToTxes)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not decode file")
	}
	file.Close()

	totalItemstoProcess := len(taskIdsToTxes)
	if removeMode {
		logger.Info().Int("tasks", totalItemstoProcess).Msg("processing tasks from file for removal")
	} else {
		logger.Info().Int("tasks", totalItemstoProcess).Msg("processing tasks from file for import")
	}

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr)) // Ensure spinner writes to stderr if needed
	if removeMode {
		s.Suffix = " removing tasks..."
	} else {
		s.Suffix = " processing tasks..."
	}
	s.FinalMSG = "completed!\n"
	s.Start()
	defer s.Stop() // Ensure spinner stops

	mapOfPendingSolTasks := make(map[task.TaskId]struct{})
	pendingSols, err := services.TaskStorage.GetAllSolutions()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get all pending solutions")
	}
	for _, v := range pendingSols {
		mapOfPendingSolTasks[v.TaskId] = struct{}{}
	}

	if removeMode {
		removedCount := 0
		skippedPendingCount := 0
		index := 0
		for taskIdStr := range taskIdsToTxes {
			index++
			select {
			case <-appQuit.Done():
				logger.Info().Msg("app quit signal received, stopping removal")
				return nil
			default:
			}

			id, err := task.ConvertTaskIdString2Bytes(taskIdStr)
			if err != nil {
				logger.Error().Err(err).Str("task", taskIdStr).Msg("could not convert task ID string")
				continue // Skip this task
			}

			// *** Check if a solution is pending for this task ***
			if _, isPending := mapOfPendingSolTasks[id]; isPending {
				logger.Warn().Str("task", taskIdStr).Msg("skipping removal: task has a pending solution")
				skippedPendingCount++
				continue // Skip deletion attempt
			}

			// Attempt to delete the task
			err = services.TaskStorage.DeleteTask(id)
			if err != nil {
				logger.Error().Err(err).Str("task", taskIdStr).Msg("failed to delete task")
			} else {
				logger.Debug().Str("task", taskIdStr).Msg("task removed successfully")
				removedCount++
			}
			s.Suffix = fmt.Sprintf(" removing tasks [%d/%d] (removed: %d, skipped_pending: %d)\n", index, totalItemstoProcess, removedCount, skippedPendingCount)
		}
		logger.Info().Int("removed", removedCount).Int("skipped_pending", skippedPendingCount).Int("total_processed", totalItemstoProcess).Msg("finished removing tasks from queue")
		return nil // End function here for removal mode
	}

	mapOfHashsByTasks := make(map[string][]task.TaskId)

	allTasks, err := services.TaskStorage.GetQueuedTasks()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get all tasks")
	}

	mapOfTaskOnwers := make(map[common.Address]int)

	uniqueTasksMap := make(map[task.TaskId]struct{})
	for _, v := range allTasks {
		uniqueTasksMap[v.TaskId] = struct{}{} // Corrected syntax
	}

	whitelistedCount := 0
	solvedCount := 0
	pendingCount := 0
	alreadyExists := 0
	index := 0
	for taskId, txHash := range taskIdsToTxes {
		index++

		select {
		case <-appQuit.Done():
			logger.Info().Msg("app quit signal received, stopping import")
			return nil
		default:
			// continue
		}

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
			logger.Debug().Str("task", taskId).Msg("task already exists in the storage tasks list")
			alreadyExists++
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
			mapOfHashsByTasks[txHash] = append(mapOfHashsByTasks[txHash], id)
		} else {
			solvedCount++
		}

		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d] (already imported: %d)\n", index, totalItemstoProcess, alreadyExists)
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

func verifyAllTasks(ctx context.Context, dryMode bool) error {

	// Get the services from the context

	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	allTasks, err := services.TaskStorage.GetAllTasks()
	if err != nil {
		services.Logger.Fatal().Err(err).Msg("could not get claims in storage")
	}

	// get all solutions from storage
	solutions, err := services.TaskStorage.GetAllSolutions()
	if err != nil {
		services.Logger.Fatal().Err(err).Msg("could not get solutions from storage")
	}

	// make a map of solutions by taskid
	solutionsMap := make(map[task.TaskId]storage.TaskData)
	for _, solution := range solutions {
		solutionsMap[solution.TaskId] = solution
	}

	// get all commitments from storage
	commitments, err := services.TaskStorage.GetAllCommitments()
	if err != nil {
		services.Logger.Fatal().Err(err).Msg("could not get commitments from storage")
	}

	// make a map of commitments by taskid
	commitmentsMap := make(map[task.TaskId]storage.TaskData)
	for _, commitment := range commitments {
		commitmentsMap[commitment.TaskId] = commitment
	}

	services.Logger.Info().Int("total", len(allTasks)).Bool("dry_mode", dryMode).Msg("verifying all tasks")

	var commitmentsToDelete []task.TaskId
	var solutionsToDelete []task.TaskId

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " processing tasks..."
	s.FinalMSG = "completed!\n"
	s.Start()
	totalItemstoProcess := len(allTasks)
	deleted := 0
	claimable := 0
	tasksUpdated := 0
	for index, v := range allTasks {
		res, err := services.Engine.Engine.Solutions(nil, v.Taskid)

		if err != nil {
			services.Logger.Fatal().Err(err).Msg("error getting solution information")
		}

		if res.Blocktime > 0 {
			// if the task is claimed, delete the task from storage and flag any commitments and solutions to delete
			if res.Claimed {
				// delete the task from storage
				if !dryMode {
					err := services.TaskStorage.DeleteTask(v.Taskid)
					if err != nil {
						services.Logger.Fatal().Err(err).Msg("could delete task data key")
					}
				}
				deleted++
			} else {
				// if the task is not claimed, make sure the task is updated for claims
				if !dryMode {
					claimTime := time.Unix(int64(res.Blocktime), 0)

					err = services.TaskStorage.UpsertTaskToClaimable(v.Taskid, common.Hash{}, claimTime)
					if err != nil {
						services.Logger.Error().Err(err).Msg("error updating task in storage")
					}
				}
				claimable++
			}
			// flag any commitments and solutions to delete
			commitmentsToDelete = append(commitmentsToDelete, v.Taskid)
			solutionsToDelete = append(solutionsToDelete, v.Taskid)
		} else {
			// No solution on-chain: Determine state based on local commitment and on-chain commitment
			// We know we have a local solution because we are iterating through GetAllSolutions()
			log.Printf("DEBUG: Task %s has no solution on-chain. Checking local/on-chain commitment...", v.Taskid.String())

			// Initialize taskStatusToSet with the current status to only update if needed
			taskStatusToSet := v.Status
			shouldDeleteCommitment := false
			shouldDeleteSolution := false

			commitmentData, hasLocalCommitment := commitmentsMap[v.Taskid]
			_, hasLocalSolution := solutionsMap[v.Taskid]

			if !hasLocalCommitment {
				// Case 1: No local commitment. Task needs to start over.
				taskStatusToSet = 0
				if hasLocalSolution {
					// If there's an orphaned local solution without a commitment, delete it.
					shouldDeleteSolution = true
					services.Logger.Debug().Str("taskid", v.Taskid.String()).Msg("No local commitment, deleting orphaned local solution.")
				} else {
					services.Logger.Debug().Str("taskid", v.Taskid.String()).Msg("No local commitment, setting status to 0.")
				}
			} else {
				// Case 2: Local commitment exists. Check its on-chain status and local solution presence.
				isOnChainCommitment := false
				if commitmentData.Commitment != [32]byte{} {
					block, err := services.Engine.Engine.Commitments(nil, commitmentData.Commitment)
					if err != nil {
						services.Logger.Error().Err(err).Str("taskid", v.Taskid.String()).Msg("Error checking on-chain commitment status, skipping task update.")
						continue // Skip update for this task if chain check fails
					}
					isOnChainCommitment = block.Uint64() > 0
				} else {
					// Commitment record exists locally but hash is zero - invalid state.
					services.Logger.Warn().Str("taskid", v.Taskid.String()).Msg("Local commitment record found with zero hash, treating as invalid.")
					// Treat as needing to start over.
					taskStatusToSet = 0
					shouldDeleteCommitment = true // Delete the invalid record
					shouldDeleteSolution = true   // Delete potentially related solution
				}

				if !hasLocalSolution {
					// Case 2a: Local commitment, but no local solution. Need to regenerate both.
					taskStatusToSet = 0
					shouldDeleteCommitment = true // Delete stale commitment
					services.Logger.Debug().Str("taskid", v.Taskid.String()).Bool("has_onchain_commitment", isOnChainCommitment).Msg("Local commitment found, but no local solution. Setting status to 0, deleting commitment.")
				} else {
					// Case 2b: Local commitment AND local solution exist.
					if isOnChainCommitment {
						// Subcase: Commitment is confirmed on-chain. Ready for solution submission.
						taskStatusToSet = 2
						shouldDeleteCommitment = true // Commitment is on-chain, remove local record
						services.Logger.Debug().Str("taskid", v.Taskid.String()).Msg("Local commitment confirmed on-chain, local solution exists. Setting status to 2, deleting commitment record.")
					} else {
						// Subcase: Commitment is NOT yet on-chain. Local solution exists.
						// This is a valid state (e.g., status 1, ready for commitment submission).
						// DO NOT change status or delete local records.
						services.Logger.Debug().Str("taskid", v.Taskid.String()).Msg("Local commitment NOT confirmed on-chain, local solution exists. State is valid, no changes needed.")
						// Ensure flags are false and taskStatusToSet remains v.Status
						shouldDeleteCommitment = false
						shouldDeleteSolution = false
					}
				}
			}

			if !dryMode {
				// Add to delete lists if flagged
				if shouldDeleteCommitment {
					commitmentsToDelete = append(commitmentsToDelete, v.Taskid)
				}
				if shouldDeleteSolution {
					solutionsToDelete = append(solutionsToDelete, v.Taskid)
				}
			} else {
				// log the changes that would be made e.g. we are deleting a commitment or solution
				services.Logger.Debug().Str("taskid", v.Taskid.String()).Bool("delete_commitment", shouldDeleteCommitment).Bool("delete_solution", shouldDeleteSolution).Msg("Would delete commitment and/or solution")
			}

			// Update task status in storage only if it has changed
			if taskStatusToSet != v.Status {
				if !dryMode {
					err = services.TaskStorage.AddOrUpdateTaskWithStatus(v.Taskid, v.Txhash, taskStatusToSet)
					if err != nil {
						services.Logger.Error().Err(err).Str("taskid", v.Taskid.String()).Int64("targetStatus", taskStatusToSet).Msg("Error updating task status in storage")
					} else {
						tasksUpdated++
						services.Logger.Info().Str("taskid", v.Taskid.String()).Int64("old_status", v.Status).Int64("new_status", taskStatusToSet).Msg("Task status updated.")
					}
				} else {
					tasksUpdated++
					services.Logger.Info().Str("taskid", v.Taskid.String()).Int64("old_status", v.Status).Int64("new_status", taskStatusToSet).Msg("Task status updated (dry run).")
				}
			} else {
				services.Logger.Debug().Str("taskid", v.Taskid.String()).Int64("status", v.Status).Msg("Task is already in the correct state.")
			}
		}

		s.Suffix = fmt.Sprintf(" processing tasks [%d/%d] [deleted: %d] [claimable: %d] [updated: %d]\n", index+1, totalItemstoProcess, deleted, claimable, tasksUpdated) // Updated suffix
	}

	if len(commitmentsToDelete) > 0 {
		services.Logger.Info().Int("commitments", len(commitmentsToDelete)).Msg("deleting commitments")
		err := services.TaskStorage.DeleteProcessedCommitments(commitmentsToDelete)
		if err != nil {
			services.Logger.Error().Err(err).Msg("error deleting commitments from storage")
		}
	}

	if len(solutionsToDelete) > 0 {
		services.Logger.Info().Int("solutions", len(solutionsToDelete)).Msg("deleting solutions")
		err := services.TaskStorage.DeleteProcessedSolutions(solutionsToDelete)
		if err != nil {
			services.Logger.Error().Err(err).Msg("error deleting solutions from storage")
		}
	}

	s.Stop()

	services.Logger.Info().Msg("completed verifying qall tasks 	")
	services.Logger.Info().Int("deleted", deleted).Msg("deleted tasks")
	services.Logger.Info().Int("claimable", claimable).Msg("new claimable tasks")
	services.Logger.Info().Int("updated", tasksUpdated).Msg("updated tasks (status)")
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

	commitments, err := services.TaskStorage.GetAllCommitments()
	if err != nil {
		services.Logger.Err(err).Msg("failed to get tasks from storage")
	}

	commitmentsMap := make(map[task.TaskId]storage.TaskData)
	for _, commitment := range commitments {
		commitmentsMap[commitment.TaskId] = commitment
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
	claimedAlready := 0
	toClaim := 0
	tasksUpdated := 0
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
			}
		}

		res, err := services.Engine.Engine.Solutions(nil, t.TaskId)

		if err != nil {
			services.Logger.Err(err).Msg("error getting solution information")
			return nil
		}

		if res.Blocktime > 0 {

			if res.Claimed {
				claimedAlready++
				services.Logger.Warn().Msgf("task %s was claimed by %s", t.TaskId.String(), res.Validator.String())
				// delete the task from storage
				err := services.TaskStorage.DeleteTask(t.TaskId)
				if err != nil {
					services.Logger.Error().Err(err).Msg("error deleting task from storage")
				}
			} else {
				toClaim++
				// update the task in storage with claim information
				// set empty txhash as we know the task exists in and will be updated
				claimTime := time.Unix(int64(res.Blocktime), 0)
				err = services.TaskStorage.UpsertTaskToClaimable(t.TaskId, common.Hash{}, claimTime)
				if err != nil {
					services.Logger.Error().Err(err).Msg("error updating task in storage")
				}
			}
			solvedByMap[res.Validator] = solvedByMap[res.Validator] + 1
			// Flag we need to delete both the commitment and the solution
			solutionsToDelete = append(solutionsToDelete, t.TaskId)
			//commitmentsToDelete = append(commitmentsToDelete, t.TaskId)
		} else {
			// No solution on-chain: Determine state based on local commitment and on-chain commitment
			// We know we have a local solution because we are iterating through GetAllSolutions()
			log.Printf("DEBUG: Task %s has no solution on-chain. Checking local/on-chain commitment...", t.TaskId.String())

			taskStatusToSet := int64(0) // Default: Requeue/Generate Commitment (Status 0)
			shouldDeleteCommitment := false
			shouldDeleteSolution := false

			commitmentData, hasLocalCommitment := commitmentsMap[t.TaskId]

			if !hasLocalCommitment {
				// Case 1: Orphaned local solution (no local commitment). Task needs to start over.
				taskStatusToSet = 0
				shouldDeleteSolution = true
				services.Logger.Debug().Str("taskid", t.TaskId.String()).Msg("Local solution found, but no local commitment. Setting status to 0, deleting solution.")
			} else {
				// Case 2: Local commitment exists. Check its on-chain status and local solution presence.
				isOnChainCommitment := false
				if commitmentData.Commitment != [32]byte{} {
					block, err := services.Engine.Engine.Commitments(nil, commitmentData.Commitment)
					if err != nil {
						services.Logger.Error().Err(err).Str("taskid", t.TaskId.String()).Msg("Error checking on-chain commitment status, skipping task update.")
						continue // Skip update for this task if chain check fails
					}
					isOnChainCommitment = block.Uint64() > 0
				} else {
					// Commitment record exists locally but hash is zero - invalid state.
					services.Logger.Warn().Str("taskid", t.TaskId.String()).Msg("Local commitment record found with zero hash, treating as invalid.")
					// Treat as needing to start over.
					taskStatusToSet = 0
					shouldDeleteCommitment = true // Delete the invalid record
					shouldDeleteSolution = true   // Delete the solution as well
				}

				// Determine status based on on-chain commitment presence
				if isOnChainCommitment {
					// Subcase: Commitment is confirmed on-chain. Ready for solution submission.
					taskStatusToSet = 2
					shouldDeleteCommitment = true // Commitment is on-chain, remove local record
					services.Logger.Debug().Str("taskid", t.TaskId.String()).Msg("Local commitment confirmed on-chain, local solution exists. Setting status to 2, deleting commitment record.")
				} else {
					// Subcase: Commitment is NOT yet on-chain. Local solution exists.
					// This is a valid state (e.g., status 1, ready for commitment submission).
					// DO NOT change status or delete local records.
					services.Logger.Debug().Str("taskid", t.TaskId.String()).Msg("Local commitment NOT confirmed on-chain, local solution exists. State is valid, no changes needed.")
					continue
				}
			}

			// Perform deletions if flagged
			if shouldDeleteCommitment {
				err = services.TaskStorage.DeleteProcessedCommitments([]task.TaskId{t.TaskId})
				if err != nil {
					log.Printf("WARN: Failed to delete local commitment for task %s: %v", t.TaskId.String(), err)
				}
			}
			if shouldDeleteSolution {
				err = services.TaskStorage.DeleteProcessedSolutions([]task.TaskId{t.TaskId})
				if err != nil {
					log.Printf("WARN: Failed to delete local solution for task %s: %v", t.TaskId.String(), err)
				}
			}

			// Update task status in storage
			err = services.TaskStorage.AddOrUpdateTaskWithStatus(t.TaskId, common.Hash{}, taskStatusToSet)
			if err != nil {
				log.Printf("WARN: Failed to add/update task %s to status %d: %v", t.TaskId.String(), taskStatusToSet, err)
			} else {
				log.Printf("Set task %s status to %d", t.TaskId.String(), taskStatusToSet)
				tasksUpdated++ // Count as added/updated
			}
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
	services.Logger.Info().Int("tasks_updated", tasksUpdated).Msg("tasks updated to push solutions")
	services.Logger.Info().Int("tasks_to_claim", toClaim).Msg("tasks to claim")
	services.Logger.Info().Int("tasks_claimed_already", claimedAlready).Msg("tasks claimed already")
	return nil

}

func verifyCommitment(ctx context.Context) error {

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	deleteCommitments := func(_commitmentsToDelete []task.TaskId) error {
		const batchSize = 1000

		if len(_commitmentsToDelete) > 0 {

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

	tasks, err := services.TaskStorage.GetAllCommitments()
	if err != nil {
		services.Logger.Err(err).Msg("failed to get tasks from storage")
	}

	services.Logger.Info().Int("commitments", len(tasks)).Msg("verifying commitments")

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

	s.Suffix = fmt.Sprintf(" deleting %d commitments", len(commitmentsToDelete))
	if err := deleteCommitments(commitmentsToDelete); err != nil {
		return nil
	}

	s.Suffix = fmt.Sprintf(" deleting %d solutions", len(solutionsToDelete))
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

// RunAutoTaskSubmit periodically submits single tasks based on config using Engine bindings.
func RunAutoTaskSubmit(appCtx context.Context, services *Services, interval time.Duration) {
	logger := services.Logger.With().Str("command", "autotasksubmit").Logger()
	logger.Info().Dur("interval", interval).Msg("Starting automatic task submitter...")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ownerAccount := services.OwnerAccount

	autoParams := services.AutoMineParams
	if autoParams == nil {
		logger.Fatal().Msg("AutoMine parameters not initialized in services")
		return
	}

	feeFloat := services.Eth.ToFloat(autoParams.Fee) // Use ToFloat
	feeFormatted := fmt.Sprintf("%.6f", feeFloat)    // Format to 6 decimal places

	logger.Info().
		Str("sender", ownerAccount.Address.Hex()).
		Str("taskOwner", autoParams.Owner.Hex()).
		Str("model", common.Bytes2Hex(autoParams.Model[:])).
		Str("fee", feeFormatted+" aius").
		Uint8("version", autoParams.Version).
		Msg("Task parameters loaded from AutoMineParams")

	// --- Submission Loop ---
	for {
		select {
		case <-appCtx.Done():
			logger.Info().Msg("Shutting down automatic task submitter...")
			return
		case <-ticker.C:
			if appCtx.Err() != nil { // Check context again after ticker
				continue
			}

			logger.Debug().Msg("Preparing to submit new task...")

			submitFunc := func(opts *bind.TransactOpts) (interface{}, error) {
				// Pass parameters directly from autoParams
				return services.Engine.Engine.SubmitTask(opts, autoParams.Version, autoParams.Owner, autoParams.Model, autoParams.Fee, autoParams.Input)
			}

			maxRetries := services.Config.Miner.ErrorMaxRetries
			backoffTime := services.Config.Miner.ErrorBackoffTime
			backoffMultiplier := services.Config.Miner.ErrorBackofMultiplier
			tx, err := ownerAccount.NonceManagerWrapper(maxRetries, backoffTime, backoffMultiplier, false, submitFunc)

			if err != nil {
				logger.Error().Err(err).Msg("Failed to submit task transaction via wrapper")
			} else {
				logger.Info().Str("txHash", tx.Hash().Hex()).Msg("Task submitted successfully using AutoMineParams")
			}
		}
	}
}

func recoverStaleTasks(ctx context.Context) error {
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	err := services.TaskStorage.RecoverStaleTasks()
	if err != nil {
		services.Logger.Error().Err(err).Msg("error recovering stale tasks")
		return err
	}

	services.Logger.Info().Msg("completed recovering stale tasks")

	return nil
}

// getUnsolvedTasks scans blocks for TaskSubmitted events and adds unsolved ones to storage.
// If endBlock is 0 or less, it scans up to the current latest block.
// If senderFilter is not the zero address, only tasks submitted by that sender are processed.
func getUnsolvedTasks(appQuit context.Context, services *Services, rpcClient *client.Client, fromBlock int64, endBlock int64, initialBlockSize int64, senderFilter common.Address) {

	eventAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	taskSubmittedEvent := eventAbi.Events["TaskSubmitted"].ID

	// Determine the end block if not provided or invalid
	if endBlock <= 0 || endBlock < fromBlock {
		log.Printf("End block not specified or invalid (%d), fetching latest block number...", endBlock)
		currentBlockUint, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			log.Fatalf("Failed to get current block number: %v", err)
		}
		endBlock = int64(currentBlockUint)
		log.Printf("Scanning up to latest block: %d", endBlock)
	}

	if fromBlock >= endBlock {
		log.Printf("From block %d is already at or beyond end block %d. Nothing to scan.", fromBlock, endBlock)
		return
	}

	// Initialize maps in the outer scope
	commitmentMap := make(map[task.TaskId]storage.TaskData)
	solutionsMap := make(map[task.TaskId]storage.TaskData)

	// first get all commitments from the db
	commitments, err := services.TaskStorage.GetAllCommitments()
	if err != nil {
		log.Printf("ERROR: Failed to get all commitments from storage: %v", err)
		return
	}

	// loop through all commitments and add them to a map
	for _, commitment := range commitments {
		commitmentMap[commitment.TaskId] = commitment
	}

	// Fetch all local solutions as well
	solutions, err := services.TaskStorage.GetAllSolutions()
	if err != nil {
		log.Printf("ERROR: Failed to get all solutions from storage: %v", err)
		return
	}

	// loop through all solutions and add them to a map
	for _, solution := range solutions {
		solutionsMap[solution.TaskId] = solution
	}

	log.Printf("Scanning for unsolved tasks from block %d to %d", fromBlock, endBlock)

	currentBlockSize := initialBlockSize
	minBlockSize := int64(100)   // Don't go smaller than this
	maxBlockSize := int64(50000) // Limit increase
	increaseStep := int64(500)   // How much to increase size on success

	// Ensure initial size is within bounds
	if currentBlockSize < minBlockSize {
		currentBlockSize = minBlockSize
	}
	if currentBlockSize > maxBlockSize {
		currentBlockSize = maxBlockSize
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{services.Config.BaseConfig.EngineAddress},
		Topics:    [][]common.Hash{{taskSubmittedEvent}},
	}

	loopFromBlock := fromBlock // Use a separate variable for loop iteration
	totalFound := 0
	totalAdded := 0

	isFiltering := senderFilter != (common.Address{}) // Check if filter is active
	if isFiltering {
		log.Printf("Filtering for tasks submitted by sender: %s", senderFilter.Hex())
	}

	for loopFromBlock <= endBlock {
		select {
		case <-appQuit.Done():
			log.Printf("WARN: getUnsolvedTasks received shutdown signal. Aborting scan at block %d.", loopFromBlock)
			return // Exit the function early
		default:
		}

		// Calculate the end block for this specific query batch
		loopToBlock := loopFromBlock + currentBlockSize - 1 // Range is inclusive
		if loopToBlock > endBlock {
			loopToBlock = endBlock
		}

		log.Printf("Querying logs from %d to %d (size %d)...", loopFromBlock, loopToBlock, currentBlockSize)
		query.FromBlock = big.NewInt(loopFromBlock)
		query.ToBlock = big.NewInt(loopToBlock)

		// Use the appQuit context for the potentially long-running FilterLogs call
		logs, err := rpcClient.FilterLogs(appQuit, query)

		if err != nil {
			// Check if error suggests block range is too large (heuristic check)
			errStr := strings.ToLower(err.Error())
			isRangeError := strings.Contains(errStr, "block range") ||
				strings.Contains(errStr, "limit") ||
				strings.Contains(errStr, "response size") ||
				strings.Contains(errStr, "timeout") ||
				strings.Contains(errStr, "too many logs")

			if isRangeError && currentBlockSize > minBlockSize {
				log.Printf("WARN: FilterLogs error (likely range too large) with size %d: %v. Reducing size.", currentBlockSize, err)
				// Reduce block size significantly, e.g., halve it, but not below min
				newSize := max(minBlockSize, currentBlockSize/2)
				if newSize == currentBlockSize { // Ensure reduction if already near min
					newSize = max(minBlockSize, newSize-increaseStep)
				}
				currentBlockSize = newSize
				log.Printf("Retrying same range with smaller block size: %d", currentBlockSize)
				// Loop continues without advancing loopFromBlock, trying the smaller size
			} else {
				// Error is not recognized as a range issue, or we're already at min size
				log.Fatalf("Unrecoverable error filtering logs from %d to %d: %v", loopFromBlock, loopToBlock, err)
			}
		} else {
			if len(logs) > 0 {
				log.Printf("Found %d TaskSubmitted logs in range %d-%d", len(logs), loopFromBlock, loopToBlock)
				totalFound += len(logs)
			}

			for _, currentlog := range logs {
				parsedLog, err := services.Engine.Engine.ParseTaskSubmitted(currentlog)
				if err != nil {
					log.Printf("WARN: Failed to parse TaskSubmitted log at block %d, tx %s: %v", currentlog.BlockNumber, currentlog.TxHash.Hex(), err)
					continue
				}

				// Apply sender filter if active
				if isFiltering && parsedLog.Sender != senderFilter {
					// log.Printf("DEBUG: Skipping task %s from sender %s (filter: %s)", task.TaskId(parsedLog.Id).String(), parsedLog.Sender.Hex(), senderFilter.Hex())
					continue // Skip this task, does not match filter
				}

				taskId := task.TaskId(parsedLog.Id)

				// Fetch on-chain solution info first, as it's the primary determinant
				sol, err := services.Engine.Engine.Solutions(nil, taskId)
				if err != nil {
					log.Printf("WARN: Failed to get solution info for task %s: %v", taskId.String(), err)
					continue // Skip if we can't get solution info
				}

				hasSolutionOnChain := sol.Blocktime != 0

				log.Printf("DEBUG: Task %s - OnChain State: Solution=%t, Claimed=%t",
					taskId.String(), hasSolutionOnChain, sol.Claimed)

				if hasSolutionOnChain {
					if sol.Claimed {
						// Solution exists and is claimed on-chain
						log.Printf("DEBUG: Task %s already claimed on-chain. Ensuring local cleanup.", taskId.String())
						// Ensure task and related records are removed locally
						_ = services.TaskStorage.DeleteTask(taskId)
						_ = services.TaskStorage.DeleteProcessedCommitments([]task.TaskId{taskId})
						_ = services.TaskStorage.DeleteProcessedSolutions([]task.TaskId{taskId})
						continue // Skip to next log entry
					} else {
						// Solution exists but is not claimed on-chain
						log.Printf("DEBUG: Task %s has unclaimed solution on-chain. Upserting to claimable locally.", taskId.String())
						// Ensure local commitment/solution records are removed first
						_ = services.TaskStorage.DeleteProcessedCommitments([]task.TaskId{taskId})
						_ = services.TaskStorage.DeleteProcessedSolutions([]task.TaskId{taskId})
						// Add/Update task to claimable state (status 3)
						claimTime := time.Unix(int64(sol.Blocktime), 0)
						err = services.TaskStorage.UpsertTaskToClaimable(taskId, currentlog.TxHash, claimTime)
						if err != nil {
							log.Printf("WARN: Failed to upsert task %s to claimable: %v", taskId.String(), err)
						} else {
							totalAdded++ // Count as added/updated
						}
						continue
					}
				} else {
					// No solution on-chain: Determine state based on local commitment and on-chain commitment
					// We know we have a local solution because we are iterating through GetAllSolutions()
					log.Printf("DEBUG: Task %s has no solution on-chain. Checking local/on-chain commitment...", taskId.String())

					taskStatusToSet := int64(0) // Default: Requeue/Generate Commitment (Status 0)
					shouldDeleteCommitment := false
					shouldDeleteSolution := false // Will be true if we reset to status 0

					commitmentData, hasLocalCommitment := commitmentMap[taskId]

					if !hasLocalCommitment {
						// Case 1: Orphaned local solution (no local commitment). Task needs to start over.
						taskStatusToSet = 0
						shouldDeleteSolution = true
						services.Logger.Debug().Str("taskid", taskId.String()).Msg("Local solution found, but no local commitment. Setting status to 0, deleting solution.")
					} else {
						// Case 2: Local commitment exists. Check its on-chain status and local solution presence.
						_, hasLocalSolution := solutionsMap[taskId] // Explicitly check for local solution

						if !hasLocalSolution {
							// Case 2a: Local commitment, but no local solution. Need to regenerate both.
							taskStatusToSet = 0
							shouldDeleteCommitment = true // Delete stale commitment
							services.Logger.Debug().Str("taskid", taskId.String()).Msg("Local commitment found, but no local solution. Setting status to 0, deleting commitment.")
							// NOTE: No 'continue' here, let it fall through to update/delete logic below
						} else {
							// Case 2b: Local commitment AND local solution exist. Check commitment's on-chain status.
							isOnChainCommitment := false
							if commitmentData.Commitment != [32]byte{} {
								block, err := services.Engine.Engine.Commitments(nil, commitmentData.Commitment)
								if err != nil {
									services.Logger.Error().Err(err).Str("taskid", taskId.String()).Msg("Error checking on-chain commitment status, skipping task update.")
									continue // Skip update for this task if chain check fails
								}
								isOnChainCommitment = block.Uint64() > 0
							} else {
								// Commitment record exists locally but hash is zero - invalid state.
								services.Logger.Warn().Str("taskid", taskId.String()).Msg("Local commitment record found with zero hash, treating as invalid.")
								// Treat as needing to start over.
								taskStatusToSet = 0
								shouldDeleteCommitment = true // Delete the invalid record
								shouldDeleteSolution = true   // Delete the solution as well
							}

							// Determine status based on on-chain commitment presence
							if isOnChainCommitment {
								// Subcase: Commitment is confirmed on-chain. Ready for solution submission.
								taskStatusToSet = 2
								shouldDeleteCommitment = true // Commitment is on-chain, remove local record
								services.Logger.Debug().Str("taskid", taskId.String()).Msg("Local commitment confirmed on-chain, local solution exists. Setting status to 2, deleting commitment record.")
							} else {
								// Subcase: Commitment is NOT yet on-chain. Local solution exists.
								// This is a valid state (e.g., status 1, ready for commitment submission).
								// DO NOT change status or delete local records.
								services.Logger.Debug().Str("taskid", taskId.String()).Msg("Local commitment NOT confirmed on-chain, local solution exists. State is valid, no changes needed.")
								// Skip to next log entry
								continue
							}
						}
					}

					// Perform deletions if flagged
					if shouldDeleteCommitment {
						err = services.TaskStorage.DeleteProcessedCommitments([]task.TaskId{taskId})
						if err != nil {
							log.Printf("WARN: Failed to delete local commitment for task %s: %v", taskId.String(), err)
						}
					}
					if shouldDeleteSolution {
						err = services.TaskStorage.DeleteProcessedSolutions([]task.TaskId{taskId})
						if err != nil {
							log.Printf("WARN: Failed to delete local solution for task %s: %v", taskId.String(), err)
						}
					}

					// Update task status in storage, using the TxHash from the log event
					err = services.TaskStorage.AddOrUpdateTaskWithStatus(taskId, currentlog.TxHash, taskStatusToSet)
					if err != nil {
						log.Printf("WARN: Failed to add/update task %s to status %d: %v", taskId.String(), taskStatusToSet, err)
					} else {
						log.Printf("Set task %s status to %d (Tx: %s)", taskId.String(), taskStatusToSet, currentlog.TxHash.Hex())
						totalAdded++ // Count as added/updated
					}
				}
			}

			// Successfully processed this range, advance to the next block
			loopFromBlock = loopToBlock + 1

			// Try increasing block size for the next iteration, up to the max
			currentBlockSize = min(maxBlockSize, currentBlockSize+increaseStep)
		}
	}

	log.Printf("Finished scanning. Total TaskSubmitted logs found: %d, Total unsolved tasks added: %d", totalFound, totalAdded)
}

// exportUnsolvedTasks scans blocks for TaskSubmitted events and exports unsolved ones to a JSON file.
// It filters by sender and block range, checking only for on-chain solutions, ignoring local storage.
// The output format matches the input format for importUnsolvedTasks.
func exportUnsolvedTasks(appQuit context.Context, services *Services, rpcClient *client.Client, fromBlock int64, endBlock int64, initialBlockSize int64, senderFilter common.Address, outputFilename string) error {

	eventAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("error getting engine abi: %w", err)
	}

	taskSubmittedEvent := eventAbi.Events["TaskSubmitted"].ID

	// Determine the end block if not provided or invalid
	if endBlock <= 0 || endBlock < fromBlock {
		log.Printf("End block not specified or invalid (%d), fetching latest block number...", endBlock)
		currentBlockUint, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get current block number: %w", err)
		}
		endBlock = int64(currentBlockUint)
		log.Printf("Scanning up to latest block: %d", endBlock)
	}

	if fromBlock >= endBlock {
		log.Printf("From block %d is already at or beyond end block %d. Nothing to scan.", fromBlock, endBlock)
		return nil
	}

	log.Printf("Scanning for unsolved tasks from block %d to %d for export to %s", fromBlock, endBlock, outputFilename)

	currentBlockSize := initialBlockSize
	minBlockSize := int64(100)   // Don't go smaller than this
	maxBlockSize := int64(50000) // Limit increase
	increaseStep := int64(500)   // How much to increase size on success

	// Ensure initial size is within bounds
	if currentBlockSize < minBlockSize {
		currentBlockSize = minBlockSize
	}
	if currentBlockSize > maxBlockSize {
		currentBlockSize = maxBlockSize
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{services.Config.BaseConfig.EngineAddress},
		Topics:    [][]common.Hash{{taskSubmittedEvent}},
	}

	loopFromBlock := fromBlock // Use a separate variable for loop iteration
	totalFound := 0
	totalExported := 0
	unsolvedTasksMap := make(map[string]string) // Map[taskIdString]txHashString

	isFiltering := senderFilter != (common.Address{}) // Check if filter is active
	if isFiltering {
		log.Printf("Filtering for tasks submitted by sender: %s", senderFilter.Hex())
	}

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " scanning blocks..."
	s.FinalMSG = "scan completed!\n"
	s.Start()
	defer s.Stop()

	for loopFromBlock <= endBlock {
		select {
		case <-appQuit.Done():
			log.Printf("WARN: exportUnsolvedTasks received shutdown signal. Aborting scan at block %d.", loopFromBlock)
			return fmt.Errorf("scan aborted") // Exit the function early
		default:
		}

		// Calculate the end block for this specific query batch
		loopToBlock := loopFromBlock + currentBlockSize - 1 // Range is inclusive
		if loopToBlock > endBlock {
			loopToBlock = endBlock
		}

		s.Suffix = fmt.Sprintf(" scanning blocks %d to %d (size %d)... found: %d exported: %d", loopFromBlock, loopToBlock, currentBlockSize, totalFound, totalExported)
		query.FromBlock = big.NewInt(loopFromBlock)
		query.ToBlock = big.NewInt(loopToBlock)

		// Use the appQuit context for the potentially long-running FilterLogs call
		logs, err := rpcClient.FilterLogs(appQuit, query)

		if err != nil {
			// Check if error suggests block range is too large (heuristic check)
			errStr := strings.ToLower(err.Error())
			isRangeError := strings.Contains(errStr, "block range") ||
				strings.Contains(errStr, "limit") ||
				strings.Contains(errStr, "response size") ||
				strings.Contains(errStr, "timeout") ||
				strings.Contains(errStr, "too many logs")

			if isRangeError && currentBlockSize > minBlockSize {
				log.Printf("WARN: FilterLogs error (likely range too large) with size %d: %v. Reducing size.", currentBlockSize, err)
				// Reduce block size significantly, e.g., halve it, but not below min
				newSize := max(minBlockSize, currentBlockSize/2)
				if newSize == currentBlockSize { // Ensure reduction if already near min
					newSize = max(minBlockSize, newSize-increaseStep)
				}
				currentBlockSize = newSize
				log.Printf("Retrying same range with smaller block size: %d", currentBlockSize)
				// Loop continues without advancing loopFromBlock, trying the smaller size
			} else {
				// Error is not recognized as a range issue, or we're already at min size
				return fmt.Errorf("unrecoverable error filtering logs from %d to %d: %w", loopFromBlock, loopToBlock, err)
			}
		} else {
			if len(logs) > 0 {
				log.Printf("Found %d TaskSubmitted logs in range %d-%d", len(logs), loopFromBlock, loopToBlock)
				totalFound += len(logs)
			}

			for _, currentlog := range logs {
				parsedLog, err := services.Engine.Engine.ParseTaskSubmitted(currentlog)
				if err != nil {
					log.Printf("WARN: Failed to parse TaskSubmitted log at block %d, tx %s: %v", currentlog.BlockNumber, currentlog.TxHash.Hex(), err)
					continue
				}

				// Apply sender filter if active
				if isFiltering && parsedLog.Sender != senderFilter {
					continue // Skip this task, does not match filter
				}

				taskId := task.TaskId(parsedLog.Id)

				// Fetch on-chain solution info ONLY
				sol, err := services.Engine.Engine.Solutions(nil, taskId)
				if err != nil {
					log.Printf("WARN: Failed to get solution info for task %s: %v", taskId.String(), err)
					continue // Skip if we can't get solution info
				}

				// If NO solution exists on-chain, add it to our export map
				if sol.Blocktime == 0 {
					taskIdStr := taskId.String()
					txHashStr := currentlog.TxHash.Hex()
					if _, exists := unsolvedTasksMap[taskIdStr]; !exists {
						unsolvedTasksMap[taskIdStr] = txHashStr
						totalExported++
						log.Printf("DEBUG: Found unsolved task %s (Tx: %s) for export.", taskIdStr, txHashStr)
					} else {
						log.Printf("DEBUG: Task %s already found, skipping duplicate.", taskIdStr)
					}
				}
			}

			// Successfully processed this range, advance to the next block
			loopFromBlock = loopToBlock + 1

			// Try increasing block size for the next iteration, up to the max
			currentBlockSize = min(maxBlockSize, currentBlockSize+increaseStep)
		}
	}

	log.Printf("Finished scanning. Total TaskSubmitted logs found: %d, Total unsolved tasks identified for export: %d", totalFound, totalExported)

	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(unsolvedTasksMap, "", "  ") // Indent for readability
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	// Write JSON data to the output file
	err = os.WriteFile(outputFilename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", outputFilename, err)
	}

	log.Printf("Successfully exported %d unsolved tasks to %s", totalExported, outputFilename)
	return nil
}
