package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/bindings/arbiusrouterv1"
	"gobius/bindings/basetoken" // Added for BaseToken ABI
	"gobius/bindings/engine"
	"gobius/client"
	task "gobius/common"
	"gobius/metrics"
	"gobius/storage"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"bytes" // Added for bytes.Equal
	"gobius/bindings/bulktasks"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types" // Added for deriving addresses
	"github.com/olekukonko/tablewriter"          // Added for table output
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
	basefeeinEth := Eth.ToFloat(basefee)

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
					} else {
						tasksUpdated++
						services.Logger.Info().Str("taskid", v.Taskid.String()).Int64("old_status", v.Status).Int64("new_status", 3).Msg("Task status updated to claimable.") // Status 3 is claimable
					}
				} else {
					tasksUpdated++
					services.Logger.Info().Str("taskid", v.Taskid.String()).Int64("old_status", v.Status).Int64("new_status", 3).Msg("Task status updated to claimable (dry run).") // Status 3 is claimable
				}
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

func sendTestPlaygroundTask(ctx context.Context) {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		log.Fatal("Could not get services from context")
	}

	ctr, err := arbiusrouterv1.NewArbiusRouterV1(services.Config.BaseConfig.ArbiusRouterAddress, services.OwnerAccount.Client.Client)
	if err != nil {
		services.Logger.Err(err).Msg("error creating arbius router")
		return
	}

	// get the baseTokenBalance on owner account as balance may change between checks
	baseTokenBalance, err := services.Basetoken.BalanceOf(nil, services.OwnerAccount.Address)
	if err != nil {
		services.Logger.Err(err).Msg("failed to get balance")
		return
	}

	allowanceAddress := services.Config.BaseConfig.ArbiusRouterAddress

	allowance, err := services.Basetoken.Allowance(nil, services.OwnerAccount.Address, allowanceAddress)
	if err != nil {
		services.Logger.Err(err).Msg("failed to get allowance")
		return
	}

	services.Logger.Debug().Msgf("allowance amount: %s", services.Config.BaseConfig.BaseToken.FormatFixed(allowance))

	// check if the allowance is less than the balance
	if allowance.Cmp(baseTokenBalance) < 0 {
		services.Logger.Info().Msgf("will need to increase allowance")

		allowanceAmount := new(big.Int).Sub(abi.MaxUint256, allowance)

		opts := services.OwnerAccount.GetOpts(0, nil, nil, nil)
		// increase the allowance
		tx, err := services.Basetoken.Approve(opts, allowanceAddress, allowanceAmount)
		if err != nil {
			services.Logger.Err(err).Msg("failed to approve allowance")
			return
		}
		// Wait for the transaction to be mined
		_, success, _, _ := services.OwnerAccount.WaitForConfirmedTx(tx)
		if !success {
			return
		}

		services.Logger.Info().Str("txhash", tx.Hash().String()).Msgf("allowance increased")
	}

	services.Logger.Info().Msgf("submitting task")

	opts := services.OwnerAccount.GetOpts(0, nil, nil, nil)

	tx, err := ctr.SubmitTask(opts, services.AutoMineParams.Version, services.AutoMineParams.Owner, services.AutoMineParams.Model, services.AutoMineParams.Fee, services.AutoMineParams.Input, big.NewInt(111), big.NewInt(1_000_000))
	if err != nil {
		services.Logger.Err(err).Msg("error submitting task")
		return
	}

	_, success, _, err := services.OwnerAccount.WaitForConfirmedTx(tx)
	if err != nil {
		services.Logger.Err(err).Msg("error waiting for task submission")
		return
	}

	if success {
		services.Logger.Info().Str("tx", tx.Hash().String()).Msg("submitted task")
	} else {
		services.Logger.Err(err).Msg("task submission failed")
	}

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
	balAsFloat := Eth.ToFloat(ethBalance)

	if balAsFloat < amount*float64(len(services.Config.BatchTasks.PrivateKeys)) {
		services.Logger.Err(err).Str("balance", fmt.Sprintf("%.4g", balAsFloat)).Msg("not enough eth balance to satisfy transfer")
		return
	}

	amountAsBig := Eth.FromFloat(amount)

	// manually update the gas base fee
	/*err = services.SenderOwnerAccount.Client.UpdateCurrentBasefee()
	if err != nil {
		services.Logger.Error().Err(err).Msg("error in new account")
		return
	}*/

	for _, pk := range services.Config.BatchTasks.PrivateKeys {

		account, err := account.NewAccount(pk, services.OwnerAccount.Client, ctx, services.Config.Blockchain.CacheNonce, services.Logger)
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
		balAsFloat := Eth.ToFloat(accountBal)

		if balAsFloat >= minbal {
			services.Logger.Info().Msgf("%s balance of %.4g eth is at or greater than minbal %.4g eth", account.Address.String(), balAsFloat, minbal)
			continue
		}

		services.OwnerAccount.UpdateNonce()
		if err != nil {
			services.Logger.Error().Err(err).Msg("error updating nonce")
			return
		}

		services.Logger.Info().Msgf("💼 transfering %.4g eth to address %s (with balance of %.4g eth)", amount, account.Address.String(), balAsFloat)

		tx, err := services.OwnerAccount.SendEther(nil, account.Address, amountAsBig)
		if err != nil {
			services.Logger.Error().Err(err).Msg("error sending transfer")
			return
		}
		_, success, _, err := services.OwnerAccount.WaitForConfirmedTx(tx)

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

	feeFloat := Eth.ToFloat(autoParams.Fee)       // Use ToFloat
	feeFormatted := fmt.Sprintf("%.6f", feeFloat) // Format to 6 decimal places

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

type GasStat struct {
	TotalGasUsed      uint64
	TotalItems        int
	TransactionCount  int
	GasPerItemSamples []float64 // Store all gas/item values
}

// calculateGasStats analyzes bulk transaction gas usage within a block range using adaptive log filtering.
func calculateGasStats(appQuit context.Context, services *Services, rpcClient *client.Client, fromBlock int64, endBlock int64, logger zerolog.Logger) error {
	logger.Info().Int64("from", fromBlock).Int64("to", endBlock).Msg("Starting gas stats calculation (Adaptive FilterLogs)")

	// Get ABIs
	engineAbi, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get engine ABI: %w", err)
	}
	bulkTasksAbi, err := bulktasks.BulkTasksMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get bulktasks ABI: %w", err)
	}

	// Method IDs
	bulkSubmitTaskID := engineAbi.Methods["bulkSubmitTask"].ID
	bulkSubmitSolutionID := engineAbi.Methods["bulkSubmitSolution"].ID
	bulkSignalCommitmentID := bulkTasksAbi.Methods["bulkSignalCommitment"].ID
	claimSolutionsID := bulkTasksAbi.Methods["claimSolutions"].ID

	// Event Topics (for initial filtering)
	taskSubmittedTopic := engineAbi.Events["TaskSubmitted"].ID
	solutionSubmittedTopic := engineAbi.Events["SolutionSubmitted"].ID
	signalCommitmentTopic := engineAbi.Events["SignalCommitment"].ID
	solutionClaimedTopic := engineAbi.Events["SolutionClaimed"].ID

	// BaseToken Transfer event details
	baseTokenAbi, err := basetoken.BaseTokenMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get BaseToken ABI: %w", err)
	}
	transferTopicHash := baseTokenAbi.Events["Transfer"].ID

	// Data structures to hold stats
	stats := map[string]*GasStat{
		"BulkSubmitTask":       {GasPerItemSamples: make([]float64, 0)},
		"BulkSubmitSolution":   {GasPerItemSamples: make([]float64, 0)},
		"BulkSignalCommitment": {GasPerItemSamples: make([]float64, 0)},
		"ClaimSolutions":       {GasPerItemSamples: make([]float64, 0)},
	}

	// Data structures for recipient earnings
	recipientEarningsTotal := make(map[common.Address]*big.Int)
	hourlyEarnings := make(map[common.Address]map[time.Time]*big.Int)
	// Note: Validator-specific filtering removed for now.
	// We will track all transfers from Engine/BulkTasks contracts.

	// Cache for block timestamps to reduce RPC calls
	blockTimestamps := make(map[uint64]time.Time)

	// Determine the end block if not provided or invalid
	if endBlock <= 0 || endBlock < fromBlock {
		currentBlockUint, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get current block number: %w", err)
		}
		endBlock = int64(currentBlockUint)
		logger.Info().Int64("block", endBlock).Msg("Scanning up to latest block")
	}

	// Fetch block timestamps for context
	var startBlockTime, endBlockTime time.Time
	startBlockBig := big.NewInt(fromBlock)
	endBlockBig := big.NewInt(endBlock)
	var captionRange string // Declare captionRange here

	startBlockInfo, err := rpcClient.Client.BlockByNumber(context.Background(), startBlockBig)
	if err != nil {
		logger.Warn().Err(err).Int64("block", fromBlock).Msg("Could not get start block info for timestamp")
	} else {
		startBlockTime = time.Unix(int64(startBlockInfo.Time()), 0)
	}

	endBlockInfo, err := rpcClient.Client.BlockByNumber(context.Background(), endBlockBig)
	if err != nil {
		logger.Warn().Err(err).Int64("block", endBlock).Msg("Could not get end block info for timestamp")
	} else {
		endBlockTime = time.Unix(int64(endBlockInfo.Time()), 0)
	}

	// Format caption details
	captionRange = fmt.Sprintf("Blocks %d-%d", fromBlock, endBlock)
	if !startBlockTime.IsZero() && !endBlockTime.IsZero() {
		startStr := startBlockTime.UTC().Format("Jan-02-2006 03:04:05 PM MST") // Use desired UTC format
		endStr := endBlockTime.UTC().Format("Jan-02-2006 03:04:05 PM MST")     // Use desired UTC format
		duration := endBlockTime.Sub(startBlockTime).Round(time.Minute)        // Round duration to nearest minute
		captionRange = fmt.Sprintf("%s (%s - %s, ~%s)", captionRange, startStr, endStr, duration.String())
	} else if !startBlockTime.IsZero() {
		startStr := startBlockTime.UTC().Format("Jan-02-2006 03:04:05 PM MST") // Use desired UTC format
		captionRange = fmt.Sprintf("%s (Starts: %s)", captionRange, startStr)
	} else if !endBlockTime.IsZero() {
		endStr := endBlockTime.UTC().Format("Jan-02-2006 03:04:05 PM MST") // Use desired UTC format
		captionRange = fmt.Sprintf("%s (Ends: %s)", captionRange, endStr)
	}

	if fromBlock > endBlock {
		logger.Warn().Int64("from", fromBlock).Int64("to", endBlock).Msg("From block is after end block, nothing to scan.")
		return nil
	}

	engineAddr := services.Config.BaseConfig.EngineAddress
	bulkTasksAddr := services.Config.BaseConfig.BulkTasksAddress
	baseTokenAddr := services.Config.BaseConfig.BaseTokenAddress

	// Adaptive step parameters
	initialBlockSize := int64(10000) // Starting point, adjust as needed
	currentBlockSize := initialBlockSize
	minBlockSize := int64(100)
	maxBlockSize := int64(50000)
	increaseStep := int64(500)
	if currentBlockSize < minBlockSize {
		currentBlockSize = minBlockSize
	}
	if currentBlockSize > maxBlockSize {
		currentBlockSize = maxBlockSize
	}

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " filtering logs..."
	s.FinalMSG = "log filtering completed!\n"
	s.Start()

	var allLogs []types.Log
	loopFromBlock := fromBlock

	// --- Adaptive Log Filtering Loop ---
	for loopFromBlock <= endBlock {
		select {
		case <-appQuit.Done():
			logger.Warn().Int64("block", loopFromBlock).Msg("Gas stats calculation cancelled during log filtering.")
			return fmt.Errorf("scan aborted")
		default:
		}

		loopToBlock := loopFromBlock + currentBlockSize - 1
		if loopToBlock > endBlock {
			loopToBlock = endBlock
		}

		s.Suffix = fmt.Sprintf(" filtering logs %d to %d (size %d)...", loopFromBlock, loopToBlock, currentBlockSize)

		filterQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(loopFromBlock),
			ToBlock:   big.NewInt(loopToBlock),
			Addresses: []common.Address{engineAddr, bulkTasksAddr, baseTokenAddr}, // Include BaseToken address
			Topics: [][]common.Hash{{
				taskSubmittedTopic, solutionSubmittedTopic, signalCommitmentTopic, solutionClaimedTopic,
				transferTopicHash, // Add Transfer topic
			}},
		}

		logs, err := rpcClient.FilterLogs(appQuit, filterQuery)
		if err != nil {
			errStr := strings.ToLower(err.Error())
			isRangeError := strings.Contains(errStr, "block range") ||
				strings.Contains(errStr, "limit") ||
				strings.Contains(errStr, "response size") ||
				strings.Contains(errStr, "timeout") ||
				strings.Contains(errStr, "too many logs")

			if isRangeError && currentBlockSize > minBlockSize {
				logger.Warn().Err(err).Int64("size", currentBlockSize).Msg("FilterLogs range error, reducing size")
				newSize := max(minBlockSize, currentBlockSize/2)
				if newSize == currentBlockSize { // Ensure reduction if already near min
					newSize = max(minBlockSize, newSize-increaseStep)
				}
				currentBlockSize = newSize
				logger.Info().Int64("new_size", currentBlockSize).Msg("Retrying with smaller block size")
				continue // Retry the same range with smaller size
			} else {
				s.Stop()
				return fmt.Errorf("unrecoverable error filtering logs from %d to %d: %w", loopFromBlock, loopToBlock, err)
			}
		} else {
			// Success
			if len(logs) > 0 {
				logger.Debug().Int("count", len(logs)).Int64("from", loopFromBlock).Int64("to", loopToBlock).Msg("Found logs in range")
				allLogs = append(allLogs, logs...)
			}

			// Advance to the next block range
			loopFromBlock = loopToBlock + 1

			// Try increasing block size for the next iteration
			currentBlockSize = min(maxBlockSize, currentBlockSize+increaseStep)
		}
	}
	s.Stop()

	logger.Info().Int("total_logs", len(allLogs)).Msg("Finished log filtering. Processing transactions...")

	// Collect unique transaction hashes from the logs
	uniqueTxHashes := make(map[common.Hash]struct{})
	for _, logEntry := range allLogs {
		uniqueTxHashes[logEntry.TxHash] = struct{}{} // Use map for easy uniqueness
	}

	if len(uniqueTxHashes) == 0 {
		logger.Info().Msg("No relevant events found in the specified block range.")
		// Display empty table (same as before)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Method", "Tx Count", "Total Items", "Avg Gas/Item", "StdDev", "Min", "P90", "P95", "P99", "Max"})
		table.SetCaption(true, fmt.Sprintf("Gas Usage Statistics (%s)", captionRange))
		for method := range stats {
			table.Append([]string{method, "0", "0", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A"})
		}
		// logger.Info().Msg("\n")
		fmt.Println("") // Print newline using fmt
		table.Render()  // Print the table
		// logger.Info().Msg("")   // Add newline after table
		fmt.Println("") // Print newline using fmt
		return nil
	}

	s.Suffix = " processing transactions..."
	s.Start()

	// Process each unique transaction (same logic as previous FilterLogs-first version)
	processedCount := 0
	for txHash := range uniqueTxHashes {
		processedCount++
		s.Suffix = fmt.Sprintf(" processing transaction %d / %d...", processedCount, len(uniqueTxHashes))

		select {
		case <-appQuit.Done():
			logger.Warn().Msg("Gas stats calculation cancelled during transaction processing.")
			return fmt.Errorf("scan aborted")
		default:
		}

		// Fetch the receipt
		receipt, err := rpcClient.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			logger.Warn().Err(err).Str("tx", txHash.Hex()).Msg("Failed to get receipt, skipping tx stats")
			continue
		}

		// --- Process Logs for Earnings (Moved inside tx loop) ---
		blockTimestamp := time.Unix(0, 0) // Placeholder for potential future hourly breakdown
		_ = blockTimestamp                // Avoid unused variable error for now

		for _, logEntry := range receipt.Logs {
			if logEntry.Address == baseTokenAddr && len(logEntry.Topics) > 0 && logEntry.Topics[0] == transferTopicHash {
				// Correctly unpack Transfer event
				var transferEvent basetoken.BaseTokenTransfer // Use generated struct
				if len(logEntry.Topics) == 3 {                // Transfer event has 3 indexed topics (event sig, from, to)
					transferEvent.From = common.BytesToAddress(logEntry.Topics[1].Bytes())
					transferEvent.To = common.BytesToAddress(logEntry.Topics[2].Bytes())
					// Unpack non-indexed amount from data
					err := baseTokenAbi.UnpackIntoInterface(&transferEvent, "Transfer", logEntry.Data)
					if err != nil {
						logger.Warn().Err(err).Str("tx", txHash.Hex()).Msg("Failed to unpack BaseToken Transfer event data")
						continue
					}
				} else {
					logger.Warn().Str("tx", txHash.Hex()).Int("topics", len(logEntry.Topics)).Msg("Unexpected number of topics for BaseToken Transfer event")
					continue
				}

				// Check if sender is Engine or BulkTasks contract
				if transferEvent.From == engineAddr || transferEvent.From == bulkTasksAddr {

					// --- Aggregate Total Earnings ---
					if existingTotal, ok := recipientEarningsTotal[transferEvent.To]; ok {
						recipientEarningsTotal[transferEvent.To] = new(big.Int).Add(existingTotal, transferEvent.Value)
					} else {
						// First time seeing this recipient
						recipientEarningsTotal[transferEvent.To] = new(big.Int).Set(transferEvent.Value)
					}

					// --- Aggregate Hourly Earnings ---
					blockNum := logEntry.BlockNumber
					blockTime, exists := blockTimestamps[blockNum]
					if !exists {
						// Fetch block info if not cached
						blockInfo, err := rpcClient.Client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
						if err != nil {
							logger.Warn().Err(err).Uint64("block", blockNum).Msg("Could not get block info for hourly timestamp")
							// Cannot aggregate hourly if block fetch fails
						} else {
							blockTime = time.Unix(int64(blockInfo.Time()), 0)
							blockTimestamps[blockNum] = blockTime // Cache it
							exists = true                         // Mark as existing now
						}
					}

					// Only aggregate if we have the timestamp
					if exists {
						hourKey := blockTime.Truncate(time.Hour).UTC() // Use UTC hour
						recipientAddr := transferEvent.To

						// Ensure nested map exists
						if _, ok := hourlyEarnings[recipientAddr]; !ok {
							hourlyEarnings[recipientAddr] = make(map[time.Time]*big.Int)
						}

						if existingHourlyTotal, ok := hourlyEarnings[recipientAddr][hourKey]; ok {
							hourlyEarnings[recipientAddr][hourKey] = new(big.Int).Add(existingHourlyTotal, transferEvent.Value)
						} else {
							hourlyEarnings[recipientAddr][hourKey] = new(big.Int).Set(transferEvent.Value)
						}
					}
				}
			}
		}
		// --- End Process Logs for Earnings ---

		// Now process Gas Stats for the *same* transaction hash
		if receipt.Status != types.ReceiptStatusSuccessful {
			continue // Skip failed transactions for gas stats calculation
		}

		// Fetch the transaction data (we already have receipt)
		tx, isPending, err := rpcClient.Client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			logger.Warn().Err(err).Str("tx", txHash.Hex()).Msg("Failed to get transaction data, skipping tx stats")
			continue
		}
		if isPending {
			logger.Warn().Str("tx", txHash.Hex()).Msg("Transaction is still pending, skipping tx stats")
			continue // Should not happen if we have a receipt, but check anyway
		}

		// Basic transaction validation
		if tx.To() == nil || len(tx.Data()) < 4 {
			continue
		}

		methodIDBytes := tx.Data()[:4]
		var targetAbi *abi.ABI
		var method *abi.Method
		var methodType string
		var itemCount int

		// Identify the contract and method based on tx.To() and methodID
		if *tx.To() == engineAddr {
			targetAbi = engineAbi
			if bytes.Equal(methodIDBytes, bulkSubmitTaskID) {
				methodType = "BulkSubmitTask"
				method, _ = targetAbi.MethodById(methodIDBytes)
			} else if bytes.Equal(methodIDBytes, bulkSubmitSolutionID) {
				methodType = "BulkSubmitSolution"
				method, _ = targetAbi.MethodById(methodIDBytes)
			} else {
				continue // Transaction to Engine, but not a target bulk method
			}
		} else if *tx.To() == bulkTasksAddr {
			targetAbi = bulkTasksAbi
			if bytes.Equal(methodIDBytes, bulkSignalCommitmentID) {
				methodType = "BulkSignalCommitment"
				method, _ = targetAbi.MethodById(methodIDBytes)
			} else if bytes.Equal(methodIDBytes, claimSolutionsID) {
				methodType = "ClaimSolutions"
				method, _ = targetAbi.MethodById(methodIDBytes)
			} else {
				continue // Transaction to BulkTasks, but not a target bulk method
			}
		} else {
			continue // Transaction not to either target contract
		}

		if method == nil {
			// This indicates an ABI mismatch or an unexpected method ID for the target contract
			logger.Warn().Str("tx", txHash.Hex()).Str("to", tx.To().Hex()).Str("methodID", hex.EncodeToString(methodIDBytes)).Msg("Could not find method in ABI for transaction")
			continue
		}

		// Unpack arguments to get item count
		args, err := method.Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			logger.Warn().Err(err).Str("tx", txHash.Hex()).Str("method", methodType).Msg("Failed to unpack arguments, skipping tx stats")
			continue
		}

		// Extract item count based on method type
		switch methodType {
		case "BulkSubmitTask":
			if len(args) >= 6 {
				if countBig, ok := args[5].(*big.Int); ok {
					itemCount = int(countBig.Int64())
				}
			}
		case "BulkSubmitSolution":
			if len(args) >= 1 {
				if tasksList, ok := args[0].([][32]byte); ok {
					itemCount = len(tasksList)
				}
			}
		case "BulkSignalCommitment":
			if len(args) >= 1 {
				if commitmentsList, ok := args[0].([][32]byte); ok {
					itemCount = len(commitmentsList)
				}
			}
		case "ClaimSolutions":
			if len(args) >= 1 {
				if tasksList, ok := args[0].([][32]byte); ok {
					itemCount = len(tasksList)
				}
			}
		default:
			logger.Error().Str("tx", txHash.Hex()).Str("methodType", methodType).Msg("Internal error: Unknown method type in count extraction")
			continue
		}

		if itemCount > 0 {
			stat := stats[methodType]
			stat.TotalGasUsed += receipt.GasUsed
			stat.TotalItems += itemCount
			stat.TransactionCount++
			gasPerItem := float64(receipt.GasUsed) / float64(itemCount)
			stat.GasPerItemSamples = append(stat.GasPerItemSamples, gasPerItem)
		} else {
			logger.Warn().Str("tx", txHash.Hex()).Str("method", methodType).Msg("Processed transaction but item count was zero or could not be extracted")
		}
	}
	s.Stop()

	// Calculate and print final stats using tablewriter
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Tx Count", "Total Items", "Avg Gas/Item", "StdDev", "Min", "P90", "P95", "P99", "Max"})
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetCaption(true, fmt.Sprintf("Gas Usage Statistics (%s)", captionRange))

	for method, stat := range stats {
		if stat.TransactionCount > 0 {
			// Calculation logic remains the same...
			sort.Float64s(stat.GasPerItemSamples)
			minGas := stat.GasPerItemSamples[0]
			maxGas := stat.GasPerItemSamples[len(stat.GasPerItemSamples)-1]

			// Average (Mean)
			sum := 0.0
			for _, v := range stat.GasPerItemSamples {
				sum += v
			}
			avgGas := sum / float64(len(stat.GasPerItemSamples))

			// Standard Deviation
			variance := 0.0
			for _, v := range stat.GasPerItemSamples {
				variance += math.Pow(v-avgGas, 2)
			}
			stdDev := math.Sqrt(variance / float64(len(stat.GasPerItemSamples)))

			// Percentiles
			percentile := func(p float64) float64 {
				if len(stat.GasPerItemSamples) == 0 {
					return 0
				}
				// Correct percentile calculation (Nearest Rank method)
				index := int(math.Ceil(p*float64(len(stat.GasPerItemSamples)))) - 1
				if index < 0 {
					index = 0
				}
				if index >= len(stat.GasPerItemSamples) {
					index = len(stat.GasPerItemSamples) - 1
				}
				return stat.GasPerItemSamples[index]
			}

			p90 := percentile(0.90)
			p95 := percentile(0.95)
			p99 := percentile(0.99)

			row := []string{
				method,
				fmt.Sprintf("%d", stat.TransactionCount),
				fmt.Sprintf("%d", stat.TotalItems),
				fmt.Sprintf("%.0f", avgGas),
				fmt.Sprintf("%.0f", stdDev),
				fmt.Sprintf("%.0f", minGas),
				fmt.Sprintf("%.0f", p90),
				fmt.Sprintf("%.0f", p95),
				fmt.Sprintf("%.0f", p99),
				fmt.Sprintf("%.0f", maxGas),
			}
			table.Append(row)
		} else {
			table.Append([]string{method, "0", "0", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A"})
		}
	}

	// logger.Info().Msg("\n") // Add newline before table
	fmt.Println("") // Print newline using fmt
	table.Render()  // Print the table
	// logger.Info().Msg("")   // Add newline after table
	fmt.Println("") // Print newline using fmt

	// --- Display Earnings Table ---
	earningsTable := tablewriter.NewWriter(os.Stdout)
	earningsTable.SetHeader([]string{"Recipient Address", "Total Received"})
	earningsTable.SetCaption(true, fmt.Sprintf("Total BaseToken Received From Engine/BulkTasks (%s)", captionRange))
	earningsTable.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	earningsTable.SetAlignment(tablewriter.ALIGN_RIGHT)

	// Sort recipients for consistent output
	recipientAddrs := make([]common.Address, 0, len(recipientEarningsTotal))
	for addr := range recipientEarningsTotal {
		recipientAddrs = append(recipientAddrs, addr)
	}
	sort.Slice(recipientAddrs, func(i, j int) bool {
		return recipientAddrs[i].Hex() < recipientAddrs[j].Hex()
	})

	tokenSymbol := services.Config.BaseConfig.BaseToken.Symbol
	grandTotalReceived := big.NewInt(0)

	for _, addr := range recipientAddrs {
		total := recipientEarningsTotal[addr]
		grandTotalReceived.Add(grandTotalReceived, total)
		row := []string{
			addr.Hex(),
			fmt.Sprintf("%s %s", services.Config.BaseConfig.BaseToken.FormatFixed(total), tokenSymbol),
		}
		earningsTable.Append(row)
	}

	if len(recipientAddrs) > 0 {
		// logger.Info().Msg("\n") // Add newline before table
		fmt.Println("") // Print newline using fmt
		earningsTable.SetFooter([]string{"GRAND TOTAL", fmt.Sprintf("%s %s", services.Config.BaseConfig.BaseToken.FormatFixed(grandTotalReceived), tokenSymbol)})
		earningsTable.Render() // Print the earnings table
		// logger.Info().Msg("")  // Add newline after table
		fmt.Println("") // Print newline using fmt
	} else {
		// logger.Info().Msg("No relevant BaseToken transfers recorded in this range.")
		fmt.Println("No relevant BaseToken transfers recorded in this range.")
	}

	// --- Display Hourly Earnings Table ---
	if len(hourlyEarnings) > 0 {
		hourlyTable := tablewriter.NewWriter(os.Stdout)
		// Prepare header: "Hour" + sorted recipient addresses
		hourlyHeader := []string{"Hour (UTC)"}
		for _, addr := range recipientAddrs { // Use the already sorted list from total earnings
			hourlyHeader = append(hourlyHeader, addr.Hex()) // Maybe truncate address later if too wide
		}
		hourlyTable.SetHeader(hourlyHeader)
		hourlyTable.SetCaption(true, fmt.Sprintf("Hourly BaseToken Received From Engine/BulkTasks (%s)", captionRange))
		hourlyTable.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
		hourlyTable.SetAlignment(tablewriter.ALIGN_RIGHT)

		// Collect and sort all unique hours
		allHoursMap := make(map[time.Time]bool)
		for _, hourlyMap := range hourlyEarnings {
			for hour := range hourlyMap {
				allHoursMap[hour] = true
			}
		}
		allHours := make([]time.Time, 0, len(allHoursMap))
		for hour := range allHoursMap {
			allHours = append(allHours, hour)
		}
		sort.Slice(allHours, func(i, j int) bool {
			return allHours[i].Before(allHours[j])
		})

		// Build rows
		for _, hour := range allHours {
			row := make([]string, len(recipientAddrs)+1)
			row[0] = hour.Format("Jan-02 3PM") // Format hour
			for i, addr := range recipientAddrs {
				if earningsMap, ok := hourlyEarnings[addr]; ok {
					if amount, ok := earningsMap[hour]; ok {
						row[i+1] = services.Config.BaseConfig.BaseToken.FormatFixed(amount)
					} else {
						row[i+1] = "0" // No earnings for this addr in this hour
					}
				} else {
					row[i+1] = "0" // No earnings for this addr at all
				}
			}
			hourlyTable.Append(row)
		}

		// Add footer with totals (already calculated)
		footerRow := make([]string, len(recipientAddrs)+1)
		footerRow[0] = "TOTAL"
		for i, addr := range recipientAddrs {
			total := recipientEarningsTotal[addr]
			footerRow[i+1] = fmt.Sprintf("%s %s", services.Config.BaseConfig.BaseToken.FormatFixed(total), tokenSymbol)
		}
		hourlyTable.SetFooter(footerRow)

		// logger.Info().Msg("\n--- Hourly Earnings ---") // Add newline and separator
		fmt.Println("\n--- Hourly Earnings ---") // Print separator using fmt
		hourlyTable.Render()
		// logger.Info().Msg("")
		fmt.Println("") // Print newline using fmt
	}

	return nil
}

// analyzeRewardRecovery scans historical blocks to analyze reward level recovery times.
func analyzeRewardRecovery(appQuit context.Context, services *Services, rpcClient *client.Client, fromBlock int64, endBlock int64, threshold float64, sampleRate int64, logger zerolog.Logger) error {
	logger.Info().Int64("from", fromBlock).Int64("to", endBlock).Float64("threshold", threshold).Int64("sampleRate", sampleRate).Msg("Starting reward recovery analysis")

	// Ensure BaseToken is configured for float conversion
	if services.Config.BaseConfig.BaseToken == nil {
		return fmt.Errorf("BaseToken configuration not found in services")
	}

	// Determine the end block if not provided or invalid
	if endBlock <= 0 || endBlock < fromBlock {
		currentBlockUint, err := rpcClient.Client.BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get current block number: %w", err)
		}
		endBlock = int64(currentBlockUint)
		logger.Info().Int64("block", endBlock).Msg("Scanning up to latest block")
	}

	if fromBlock >= endBlock {
		logger.Warn().Int64("from", fromBlock).Int64("to", endBlock).Msg("From block is after end block, nothing to scan.")
		return nil
	}

	// State tracking
	isBelowThreshold := false
	var blockEnteredBelowThreshold int64
	var timestampEnteredBelowThreshold time.Time
	recoveryDurations := []time.Duration{}
	blockTimestamps := make(map[int64]time.Time) // Cache block timestamps

	// Helper to get block timestamp
	getBlockTimestamp := func(blockNum int64) (time.Time, error) {
		if ts, ok := blockTimestamps[blockNum]; ok {
			return ts, nil
		}
		blockInfo, err := rpcClient.Client.BlockByNumber(context.Background(), big.NewInt(blockNum))
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get block info for %d: %w", blockNum, err)
		}
		ts := time.Unix(int64(blockInfo.Time()), 0)
		blockTimestamps[blockNum] = ts // Cache it
		return ts, nil
	}

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " scanning blocks..."
	s.FinalMSG = "block scan completed!\n"
	s.Start()
	defer s.Stop()

	// Main scanning loop
	for currentBlock := fromBlock; currentBlock <= endBlock; currentBlock += sampleRate {
		select {
		case <-appQuit.Done():
			logger.Warn().Int64("block", currentBlock).Msg("Reward analysis cancelled.")
			return fmt.Errorf("scan aborted")
		default:
		}

		s.Suffix = fmt.Sprintf(" scanning block %d / %d...", currentBlock, endBlock)

		callOpts := &bind.CallOpts{
			Context:     context.Background(), // Use background context for calls
			BlockNumber: big.NewInt(currentBlock),
		}

		rewardBigInt, err := services.Engine.Engine.GetReward(callOpts)
		if err != nil {
			// Log warning and continue, maybe RPC node was temporarily unavailable for that block
			logger.Warn().Err(err).Int64("block", currentBlock).Msg("Failed to get reward for block, skipping")
			continue
		}

		rewardFloat := services.Config.BaseConfig.BaseToken.ToFloat(rewardBigInt)

		logger.Debug().Int64("block", currentBlock).Float64("reward", rewardFloat).Msg("Sampled reward")

		if rewardFloat < threshold {
			if !isBelowThreshold {
				// Just dipped below threshold
				isBelowThreshold = true
				blockEnteredBelowThreshold = currentBlock
				timestamp, err := getBlockTimestamp(currentBlock)
				if err != nil {
					logger.Warn().Err(err).Int64("block", currentBlock).Msg("Could not get timestamp when entering below threshold, time calculation might be inaccurate")
					timestampEnteredBelowThreshold = time.Time{} // Mark as invalid
				} else {
					timestampEnteredBelowThreshold = timestamp
					logger.Info().Int64("block", currentBlock).Float64("reward", rewardFloat).Time("time", timestamp).Msg("Reward dropped below threshold")
				}
			}
			// Still below threshold, do nothing else
		} else {
			if isBelowThreshold {
				// Just recovered above or equal to threshold
				logger.Info().Int64("block", currentBlock).Float64("reward", rewardFloat).Msg("Reward recovered to threshold")
				timestampRecovered, err := getBlockTimestamp(currentBlock)
				if err != nil {
					logger.Warn().Err(err).Int64("block", currentBlock).Msg("Could not get timestamp when recovering, cannot calculate duration")
				} else if !timestampEnteredBelowThreshold.IsZero() {
					duration := timestampRecovered.Sub(timestampEnteredBelowThreshold)
					recoveryDurations = append(recoveryDurations, duration)
					logger.Info().Dur("duration", duration.Round(time.Second)).Int64("from_block", blockEnteredBelowThreshold).Int64("to_block", currentBlock).Msg("Recovery duration recorded")
				} else {
					logger.Warn().Int64("from_block", blockEnteredBelowThreshold).Int64("to_block", currentBlock).Msg("Recovery detected but entry timestamp was missing, cannot calculate duration")
				}
				// Reset state regardless of timestamp success
				isBelowThreshold = false
			}
			// Stayed above threshold, do nothing else
		}
	}

	logger.Info().Int("recoveries", len(recoveryDurations)).Msg("Finished reward analysis scan.")

	if len(recoveryDurations) > 0 {
		var totalDuration time.Duration
		minDuration := recoveryDurations[0]
		maxDuration := recoveryDurations[0]

		for _, d := range recoveryDurations {
			totalDuration += d
			if d < minDuration {
				minDuration = d
			}
			if d > maxDuration {
				maxDuration = d
			}
		}
		avgDuration := totalDuration / time.Duration(len(recoveryDurations))

		// logger.Info().
		// 	Int("count", len(recoveryDurations)).
		// 	Str("avg", avgDuration.Round(time.Second).String()).
		// 	Str("min", minDuration.Round(time.Second).String()).
		// 	Str("max", maxDuration.Round(time.Second).String()).
		// 	Msg("Recovery Duration Summary")
		fmt.Printf("Recovery Duration Summary: Count=%d, Avg=%s, Min=%s, Max=%s\n",
			len(recoveryDurations),
			avgDuration.Round(time.Second).String(),
			minDuration.Round(time.Second).String(),
			maxDuration.Round(time.Second).String(),
		)
	} else {
		// logger.Info().Msg("No complete recovery periods observed within the scanned range.")
		fmt.Println("No complete recovery periods observed within the scanned range.")
	}

	return nil
}

// cleanQueueLocal removes tasks from the local queue (status 0) if their on-chain owner
// does not match the miner's OwnerAccount address.
func cleanQueueLocal(appQuit context.Context, ctx context.Context) error {
	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		return fmt.Errorf("could not get services from context")
	}

	logger := services.Logger
	ownerAddress := services.OwnerAccount.Address

	logger.Info().Str("owner", ownerAddress.Hex()).Msg("Starting local task queue cleanup based on owner address...")

	queuedTasks, err := services.TaskStorage.GetQueuedTasks()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get queued tasks from storage")
		return err
	}

	if len(queuedTasks) == 0 {
		logger.Info().Msg("Local task queue is empty, no cleanup needed.")
		return nil
	}

	totalTasks := len(queuedTasks)
	tasksToDelete := make([]task.TaskId, 0)
	checkedCount := 0
	mismatchCount := 0

	s := spinner.New(spinner.CharSets[11], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " checking tasks..."
	s.FinalMSG = "task check completed!\n"
	s.Start()
	defer s.Stop()

	for _, queuedTask := range queuedTasks {
		select {
		case <-appQuit.Done():
			logger.Warn().Msg("Cleanup process cancelled by user.")
			return fmt.Errorf("cleanup cancelled")
		default:
		}

		checkedCount++
		s.Suffix = fmt.Sprintf(" checking tasks [%d/%d] (mismatched: %d)...", checkedCount, totalTasks, mismatchCount)

		// Get task info from the blockchain
		taskInfo, err := services.Engine.Engine.Tasks(nil, queuedTask.TaskId)
		if err != nil {
			// Log error but continue; maybe the task doesn't exist on-chain anymore
			logger.Warn().Err(err).Str("task", queuedTask.TaskId.String()).Msg("Failed to get task info from blockchain, skipping task.")
			continue
		}

		// Compare owner address
		if taskInfo.Owner != ownerAddress {
			mismatchCount++
			tasksToDelete = append(tasksToDelete, queuedTask.TaskId)
			logger.Debug().Str("task", queuedTask.TaskId.String()).Str("expected_owner", ownerAddress.Hex()).Str("actual_owner", taskInfo.Owner.Hex()).Msg("Owner mismatch detected, marking task for deletion.")
		}
	}

	s.Stop() // Stop spinner before final logs

	if len(tasksToDelete) > 0 {
		logger.Info().Int("count", len(tasksToDelete)).Msg("Deleting tasks with mismatched owners from local storage...")
		// Delete tasks individually as DeleteTasks is not available
		deletedCount := 0
		failedCount := 0
		for _, taskID := range tasksToDelete {
			err = services.TaskStorage.DeleteTask(taskID)
			if err != nil {
				logger.Error().Err(err).Str("task", taskID.String()).Msg("Failed to delete task from storage")
				failedCount++
			} else {
				deletedCount++
			}
		}

		if failedCount > 0 {
			logger.Warn().Int("deleted", deletedCount).Int("failed", failedCount).Msg("Finished deleting tasks with some errors.")
			// Return an error if any deletion failed
			return fmt.Errorf("failed to delete %d tasks", failedCount)
		} else {
			logger.Info().Int("deleted_count", deletedCount).Msg("Successfully deleted tasks with mismatched owners.")
		}
	} else {
		logger.Info().Msg("No tasks found with mismatched owners. Local queue is clean.")
	}

	logger.Info().Msg("Local task queue cleanup finished.")
	return nil
}
