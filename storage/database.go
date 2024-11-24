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
	logger       *zerolog.Logger
}

func NewTaskStorageDB(ctx context.Context, sql *sql.DB, minclaimtime time.Duration, logger *zerolog.Logger) *TaskStorageDB {

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
		ts.logger.Println("Validator:", v.Validator, "Solution Pending Count:", v.SolutionCount)
		mapValToCounts[v.Validator] = v.SolutionCount
	}

	return mapValToCounts, err
}

func (ts *TaskStorageDB) AddCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error {
	err := ts.queries.CreateCommitment(ts.ctx, db.CreateCommitmentParams{
		Taskid:     taskId,
		Commitment: commitment,
		Validator:  validator,
	})

	return err
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
	start := time.Now()

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
	ts.logger.Println("DeleteProcessedCommitments time:", time.Since(start))

	return nil
}

func (ts *TaskStorageDB) DeleteProcessedSolutions(taskIds []task.TaskId) error {
	start := time.Now()

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
	ts.logger.Println("DeleteProcessedSolutions time:", time.Since(start))

	return nil
}

// func (ts *TaskStorageDB) AddTaskToClaim(taskId task.TaskId) (time.Time, error) {
// 	// as per enginev2 spec. 2000 seconds claim time + a little bit of buffer e.g. ~35mins
// 	// this value can be changed so we read it from the contract at load time
// 	claimTime := time.Now().Add(ts.minclaimtime)
// 	claimTimeAsUnixString := strconv.FormatInt(claimTime.Unix(), 10)

// 	return claimTime, err
// }

func (ts *TaskStorageDB) AddTasksToClaim(taskIds []task.TaskId, value float64) (time.Time, error) {
	start := time.Now()

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
	ts.logger.Println("AddTasksToClaim time:", time.Since(start))

	return claimTime, nil
}

// func (rc *TaskStorageDB) GetClaims(batchSize int) (ClaimTaskSlice, error) {

// 	return claims, nil
// }

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
func (ts *TaskStorageDB) AddTask(task string, txhash string) error {
	return nil
}

func (ts *TaskStorageDB) AddTasks(tasks []task.TaskId, txhash common.Hash, gasPerTask float64) error {

	start := time.Now()

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
			Cumulativegas: 0.0000000001,
			Txhash:        txhash,
		})
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	ts.logger.Println("AddTasks time:", time.Since(start))

	return nil
}

// // removes and returns a task from the list
// func (ts *TaskStorageDB) PopTaskOld() (task.TaskId, common.Hash, error) {
// 	start := time.Now()

// 	tx, err := ts.sqlite.BeginTx(ts.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
// 	if err != nil {
// 		return task.TaskId{}, common.Hash{}, err
// 	}
// 	qtx := ts.queries.WithTx(tx)

// 	defer tx.Rollback()

// 	row, err := qtx.PopTaskOld(ts.ctx)
// 	if err != nil {
// 		return task.TaskId{}, common.Hash{}, err
// 	}

// 	count, err := qtx.SetTaskQueuedStatus(ts.ctx, row.Taskid)
// 	if err != nil {
// 		return task.TaskId{}, common.Hash{}, err
// 	}

// 	if count == 0 {
// 		log.Println("no rows affected by set task queue status")
// 		return task.TaskId{}, common.Hash{}, sql.ErrNoRows
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return task.TaskId{}, common.Hash{}, err
// 	}
// 	log.Println("popTaskFromQueue time:", time.Since(start))

// 	return row.Taskid, row.Txhash, nil
// }

func beginImmediate(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err == nil {
		_, err = tx.Exec("ROLLBACK; BEGIN IMMEDIATE")
	}
	return tx, err
}

func (ts *TaskStorageDB) PopTask() (task.TaskId, common.Hash, error) {
	start := time.Now()

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
	ts.logger.Println("popTaskFromQueue time:", time.Since(start))

	return row.Taskid, row.Txhash, nil
}

func (ts *TaskStorageDB) GetClaims(batchSize int) (ClaimTaskSlice, float64, error) {
	claimsFromDb, err := ts.queries.GetTasksByLowestCost(ts.ctx, db.GetTasksByLowestCostParams{
		Claimtime: time.Now().Unix() * 2,
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
	start := time.Now()

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
	ts.logger.Println("UpdateTaskStatusAndCost time:", time.Since(start))

	return nil
}
