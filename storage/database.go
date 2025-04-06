package storage

import (
	"context"
	"database/sql"
	task "gobius/common"
	db "gobius/sql/sqlite"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	"time"
)

// TODO: make task storage an interface which can be implemented by different storage backends

type TaskStorageDB struct {
	sqlite       *sql.DB
	queries      *db.Queries
	minclaimtime time.Duration
	ctx          context.Context
	logger       zerolog.Logger
}

func NewTaskStorageDB(ctx context.Context, sql *sql.DB, minclaimtime time.Duration, logger zerolog.Logger) *TaskStorageDB {

	queries := db.New(sql)

	ts := &TaskStorageDB{
		ctx:          ctx,
		sqlite:       sql,
		queries:      queries,
		minclaimtime: minclaimtime,
		logger:       logger,
	}

	return ts
}

func (ts *TaskStorageDB) MinClaimTime() time.Duration {
	return ts.minclaimtime
}

func (ts *TaskStorageDB) GetPendingSolutionsCountPerValidator() (map[common.Address]int64, error) {
	valData, err := ts.queries.GetPendingSolutionsCountPerValidator(ts.ctx)
	if err != nil {
		return nil, err
	}

	mapValToCounts := map[common.Address]int64{}

	for _, v := range valData {
		//ts.logger.Println("Validator:", v.Validator, "Solution Pending Count:", v.SolutionCount)
		mapValToCounts[v.Validator] = v.SolutionCount
	}

	return mapValToCounts, err
}

// TryAddCommitment attempts to add a commitment to the database
// Returns true if the commitment was added, false if it already exists
// Returns an error if there is an error adding the commitment
func (ts *TaskStorageDB) TryAddCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) (bool, error) {
	exists, err := ts.queries.CheckCommitmentExists(ts.ctx, taskId)
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return false, nil // Commitment already exists, skip adding
	}

	err = ts.queries.CreateCommitment(ts.ctx, db.CreateCommitmentParams{
		Taskid:     taskId,
		Commitment: commitment,
		Validator:  validator,
	})

	return true, err
}

func (ts *TaskStorageDB) AddSolution(validator common.Address, taskId task.TaskId, cid []byte) error {

	err := ts.queries.CreateSolution(ts.ctx, db.CreateSolutionParams{
		Taskid:    taskId,
		Cid:       cid,
		Validator: validator,
	})

	return err
}

func (ts *TaskStorageDB) GetPendingSolutions(validator common.Address, batchSize int) (TaskDataSlice, error) {
	sols, err := ts.queries.GetSolutionBatch(ts.ctx, db.GetSolutionBatchParams{
		Validator: validator,
		Limit:     int64(batchSize),
	})

	if err != nil {
		return nil, err
	}

	tasks := make(TaskDataSlice, len(sols))

	for i, s := range sols {
		tasks[i] = TaskData{
			TaskId:     s.Taskid,
			Commitment: [32]byte{},
			Solution:   s.Cid,
		}
	}

	return tasks, err
}

func (ts *TaskStorageDB) GetAllSolutions() (TaskDataSlice, error) {
	sols, err := ts.queries.GetSolutions(ts.ctx)

	if err != nil {
		return nil, err
	}

	tasks := make(TaskDataSlice, len(sols))

	for i, s := range sols {
		tasks[i] = TaskData{
			TaskId:   s.Taskid,
			Solution: s.Cid,
		}
	}

	return tasks, err
}

func (ts *TaskStorageDB) GetQueuedTasks() (TaskDataSlice, error) {
	alltasks, err := ts.queries.GetQueuedTasks(ts.ctx)

	if err != nil {
		return nil, err
	}

	tasks := make(TaskDataSlice, len(alltasks))

	for i, s := range alltasks {
		tasks[i] = TaskData{
			TaskId: s.Taskid,
		}
	}

	return tasks, err
}

func (ts *TaskStorageDB) GetPendingCommitments(batchSize int) (TaskDataSlice, error) {
	commies, err := ts.queries.GetCommitmentBatch(ts.ctx, int64(batchSize))

	if err != nil {
		return nil, err
	}

	tasks := make(TaskDataSlice, len(commies))
	for i, c := range commies {
		tasks[i] = TaskData{
			TaskId:     c.Taskid,
			Commitment: c.Commitment,
			Solution:   []byte{},
		}
	}

	return tasks, nil
}

func (ts *TaskStorageDB) DeleteProcessedCommitments(taskIds []task.TaskId) error {
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := ts.queries.WithTx(tx)

	defer tx.Rollback()

	for _, v := range taskIds {
		qtx.DeleteCommitment(ts.ctx, v)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ts *TaskStorageDB) DeleteProcessedSolutions(taskIds []task.TaskId) error {
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := ts.queries.WithTx(tx)

	defer tx.Rollback()

	for _, v := range taskIds {
		qtx.DeleteSolution(ts.ctx, v)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ts *TaskStorageDB) AddTasksToClaim(taskIds []task.TaskId, value float64) (time.Time, error) {

	claimTime := time.Now().Add(ts.minclaimtime)

	// start a transaction
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return claimTime, err
	}
	qtx := ts.queries.WithTx(tx)

	defer tx.Rollback()

	for _, taskId := range taskIds {
		qtx.UpdateTaskSolution(ts.ctx, db.UpdateTaskSolutionParams{
			Taskid:        taskId,
			Claimtime:     claimTime.Unix(),
			Cumulativegas: value,
		})
	}

	if err := tx.Commit(); err != nil {
		return claimTime, err
	}

	return claimTime, nil
}

func (ts *TaskStorageDB) DeleteClaim(taskkey string) error {
	return nil
}

func (ts *TaskStorageDB) DeleteClaims(tasks []task.TaskId) error {
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := ts.queries.WithTx(tx)
	start := time.Now()

	defer tx.Rollback()

	for _, v := range tasks {
		qtx.DeleteCommitment(ts.ctx, v)
		qtx.DeleteSolution(ts.ctx, v)
		qtx.DeletedClaimedTask(ts.ctx, v)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	ts.logger.Println("DeleteClaims time:", time.Since(start))

	return nil

}

func (ts *TaskStorageDB) TotalTasks() (int64, error) {
	count, err := ts.queries.TotalPendingTasks(ts.ctx)

	return count, err
}

func (ts *TaskStorageDB) TotalCommitments() (int64, error) {
	commitments, err := ts.queries.TotalCommitments(ts.ctx)

	return commitments, err
}

func (ts *TaskStorageDB) TotalSolutionsAndClaims() (int64, int64, error) {
	counts, err := ts.queries.TotalSolutionsAndClaims(ts.ctx)

	return counts.TotalSolutions, counts.TotalClaims, err
}

// adds a task to the list
func (ts *TaskStorageDB) AddTask(taskId task.TaskId, txhash common.Hash, gasPerTask float64) error {

	err := ts.queries.AddTask(ts.ctx, db.AddTaskParams{
		Taskid:        taskId,
		Cumulativegas: gasPerTask,
		Txhash:        txhash,
	})

	return err
}

// adds a task to the list with explicit status
func (ts *TaskStorageDB) AddTaskWithStatus(taskId task.TaskId, txhash common.Hash, gasPerTask float64, status int64) error {

	err := ts.queries.AddTaskWithStatus(ts.ctx, db.AddTaskWithStatusParams{
		Taskid:        taskId,
		Cumulativegas: gasPerTask,
		Txhash:        txhash,
		Status:        status,
	})

	return err
}

func (ts *TaskStorageDB) AddOrUpdateTaskWithStatus(taskId task.TaskId, txhash common.Hash, status int64) error {
	err := ts.queries.AddOrUpdateTaskWithStatus(ts.ctx, db.AddOrUpdateTaskWithStatusParams{
		Taskid: taskId,
		Txhash: txhash,
		Status: status,
	})

	return err
}

func (ts *TaskStorageDB) AddTasks(tasks []task.TaskId, txhash common.Hash, gasPerTask float64) error {

	// start a transaction
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := ts.queries.WithTx(tx)

	defer tx.Rollback()

	for _, taskId := range tasks {
		qtx.AddTask(ts.ctx, db.AddTaskParams{
			Taskid:        taskId,
			Cumulativegas: gasPerTask,
			Txhash:        txhash,
		})
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func beginImmediate(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err == nil {
		_, err = tx.Exec("ROLLBACK; BEGIN IMMEDIATE")
	}
	return tx, err
}

func (ts *TaskStorageDB) PopTask() (task.TaskId, common.Hash, error) {
	tx, err := beginImmediate(ts.sqlite)
	if err != nil {
		return task.TaskId{}, common.Hash{}, err
	}

	qtx := ts.queries.WithTx(tx)
	defer tx.Rollback()

	row, err := qtx.PopTask(ts.ctx)
	if err != nil {
		return task.TaskId{}, common.Hash{}, err
	}
	if err := tx.Commit(); err != nil {
		return task.TaskId{}, common.Hash{}, err
	}

	return row.Taskid, row.Txhash, nil
}

func (ts *TaskStorageDB) GetClaims(batchSize int) (ClaimTaskSlice, float64, error) {
	claimsFromDb, err := ts.queries.GetTasksByLowestCost(ts.ctx, db.GetTasksByLowestCostParams{
		Claimtime: time.Now().Unix(),
		Limit:     int64(batchSize),
	})

	if err != nil {
		return nil, 0, err
	}

	claims := make(ClaimTaskSlice, 0, len(claimsFromDb))

	averageGas := 0.0
	averageSamples := 0

	for _, v := range claimsFromDb {
		averageGas += v.Cumulativegas
		averageSamples++

		claim := ClaimTask{
			ID:        v.Taskid,
			Time:      v.Claimtime,
			TotalCost: v.Cumulativegas,
		}
		claims = append(claims, claim)
	}

	return claims, averageGas / float64(averageSamples), nil
}

func (ts *TaskStorageDB) UpdateTaskStatusAndCost(tasks []task.TaskId, status int64, value float64) error {
	// start a transaction
	tx, err := ts.sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := ts.queries.WithTx(tx)

	defer tx.Rollback()

	for _, taskId := range tasks {
		_, err := qtx.UpdateTaskStatusAndGas(ts.ctx, db.UpdateTaskStatusAndGasParams{
			Taskid:        taskId,
			Cumulativegas: value,
			Status:        status,
		})
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ts *TaskStorageDB) AddIpfsCid(taskId task.TaskId, cid []byte) error {
	err := ts.queries.AddIPFSCid(ts.ctx, db.AddIPFSCidParams{
		Taskid: taskId,
		Cid:    cid,
	})

	return err
}

func (ts *TaskStorageDB) GetIpfsCids(batchSize int) (TaskIpfsCidDataSlice, error) {
	ipfsCidsFromDb, err := ts.queries.GetIPFSCids(ts.ctx, int64(batchSize))

	if err != nil {
		return nil, err
	}

	ipfsCids := make(TaskIpfsCidDataSlice, len(ipfsCidsFromDb))

	for i, item := range ipfsCidsFromDb {
		ipfsCids[i] = TaskIpfsCidData{
			TaskId: item.Taskid,
			Cid:    item.Cid,
			Added:  item.Added,
		}
	}

	return ipfsCids, nil
}

func (ts *TaskStorageDB) DeleteIpfsCid(taskId task.TaskId) error {
	_, err := ts.queries.DeletedIPFSCid(ts.ctx, taskId)

	return err
}

func (ts *TaskStorageDB) RecoverStaleTasks() error {
	err := ts.queries.RecoverStaleTasks(ts.ctx)

	return err
}

// get all commiments:
func (ts *TaskStorageDB) GetAllCommitments() (TaskDataSlice, error) {
	commitments, err := ts.queries.GetCommitments(ts.ctx)

	if err != nil {
		return nil, err
	}

	tasks := make(TaskDataSlice, len(commitments))

	for i, c := range commitments {
		tasks[i] = TaskData{
			TaskId:     c.Taskid,
			Commitment: c.Commitment,
		}
	}

	return tasks, nil
}

// RequeueTaskIfNoCommitmentOrSolution re-enqueues a task if it has no commitment or solution
func (ts *TaskStorageDB) RequeueTaskIfNoCommitmentOrSolution(taskId task.TaskId) (requeued bool, err error) {
	rows, err := ts.queries.RequeueTaskIfNoCommitmentOrSolution(ts.ctx, taskId)

	return rows > 0, err
}

func (ts *TaskStorageDB) DeleteTask(taskid task.TaskId) error {
	_, err := ts.queries.DeletedTask(ts.ctx, taskid)
	return err
}

func (ts *TaskStorageDB) UpsertTaskToClaimable(taskid task.TaskId, txhash common.Hash, claimTime time.Time) error {
	claimTime = claimTime.Add(ts.minclaimtime)

	params := db.UpsertTaskToClaimableParams{
		Taskid:    taskid,
		Txhash:    txhash,
		Claimtime: claimTime.Unix(),
	}
	return ts.queries.UpsertTaskToClaimable(ts.ctx, params)
}

func (ts *TaskStorageDB) GetAllTasks() ([]db.Task, error) {
	tasks, err := ts.queries.GetAllTasks(ts.ctx)

	return tasks, err
}
