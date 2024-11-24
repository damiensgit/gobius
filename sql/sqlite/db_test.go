package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	_ "embed"
	task "gobius/common"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tj/assert"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func setup(t *testing.T) *sql.DB {
	os.Remove("dick.db")
	sqlite, err := sql.Open("sqlite3", "dick.db")

	assert.NoError(t, err, "err in Marshal")

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlite, "migrations"); err != nil {
		log.Fatal(err)
	}

	// create tables
	//_, err = sqlite.ExecContext(ctx, ddl)
	//assert.NoError(t, err, "err in Marshal")

	_, err = sqlite.Exec("VACUUM")
	if err != nil {
		log.Fatal(err)
	}

	return sqlite
}

func TestQueue(t *testing.T) {
	ctx := context.Background()

	sqlite := setup(t)

	queries := New(sqlite)

	var jobid int64 = 0
	var wg sync.WaitGroup

	numWorkers := 30

	newTask := make(chan task.TaskId, 2001)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int, wg *sync.WaitGroup) {
			defer wg.Done()

			for t := range newTask {

				log.Printf("%d: got new task so will work on that for 4secs: %s", jobid, t.String())

				time.Sleep(4 * time.Second)

				atomic.AddInt64(&jobid, 1)

			}

		}(i, &wg)
	}

	go func() {

		poppedCount := 0

		for {

			taskId, err := popTaskFromQueue(ctx, sqlite, queries)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("queue is empty, nothing to do, will sleep. popped so far: %d", poppedCount)
					time.Sleep(2 * time.Second)
					// The list is empty, continue to the next iteration
					continue
				}
				log.Printf("err: %s, popped so far: %d", err, poppedCount)
				time.Sleep(5 * time.Second)
				continue
			}

			newTask <- taskId

			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {

			totalTasks, err := queries.TotalQueuedTasks(ctx)
			log.Println("Tasks queued:", totalTasks, err)

			totalTasksToMake := 200

			addedTasks := make([]task.TaskId, 0)
			for i := 0; i < totalTasksToMake; i++ {
				taskid := generateRandom32Bytes()
				addedTasks = append(addedTasks, taskid)
			}

			tx, err := sqlite.Begin()
			assert.NoError(t, err, "err in Begin")
			qtx := queries.WithTx(tx)
			start := time.Now()
			for _, v := range addedTasks {
				err = qtx.AddTask(ctx, AddTaskParams{
					Taskid:        v,
					Cumulativegas: 0.0000000001,
				})
				assert.NoError(t, err, "err in AddTasks")

			}
			tx.Commit()
			tx.Rollback()
			log.Println("AddTask time:", time.Since(start))
		}
	}()

	select {}

}

func popTaskFromQueue(ctx context.Context, sqlite *sql.DB, queries *Queries) (task.TaskId, error) {
	tx, err := sqlite.Begin()
	if err != nil {
		return task.TaskId{}, err
	}
	qtx := queries.WithTx(tx)
	start := time.Now()

	defer tx.Rollback()

	poppedTask, err := qtx.PopTask(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows case
			return task.TaskId{}, nil
		} else {
			// Handle other errors
			return task.TaskId{}, err
		}
	}

	qtx.UpdateTaskStatus(ctx, UpdateTaskStatusParams{
		Taskid: poppedTask,
		Status: 1,
	})

	if err := tx.Commit(); err != nil {
		return task.TaskId{}, err
	}
	log.Println("popTaskFromQueue time:", time.Since(start))

	return poppedTask, nil
}

func deleteTask(ctx context.Context, sqlite *sql.DB, queries *Queries, taskId task.TaskId) error {
	tx, err := sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := queries.WithTx(tx)
	start := time.Now()

	defer tx.Rollback()

	qtx.DeleteCommitment(ctx, taskId)
	qtx.DeleteSolution(ctx, taskId)
	qtx.DeletedClaimedTask(ctx, taskId)

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Println("deleteTask time:", time.Since(start))

	return nil
}

func deleteTasks(ctx context.Context, sqlite *sql.DB, queries *Queries, tasks []Task) error {
	tx, err := sqlite.Begin()
	if err != nil {
		return err
	}
	qtx := queries.WithTx(tx)
	start := time.Now()

	defer tx.Rollback()

	for _, v := range tasks {
		qtx.DeleteCommitment(ctx, v.Taskid)
		qtx.DeleteSolution(ctx, v.Taskid)
		qtx.DeletedClaimedTask(ctx, v.Taskid)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Println("deleteTasks time:", time.Since(start))

	return nil
}

func TestDb(t *testing.T) {
	ctx := context.Background()

	sqlite := setup(t)

	queries := New(sqlite)

	commitments, err := queries.GetCommitments(ctx, 10)
	assert.NoError(t, err, "err in Marshal")

	log.Println(commitments)

	totalTasksToMake := 5000
	totalTasksToSolve := 1000

	validator1 := common.HexToAddress("0x1111111111111111111111111111111111111111")

	addedTasks := make([]task.TaskId, 0)
	for i := 0; i < totalTasksToMake; i++ {
		taskid := generateRandom32Bytes()
		addedTasks = append(addedTasks, taskid)
	}

	/*start := time.Now()
	for _, v := range addedTasks {
		_, err = queries.AddTask(ctx, AddTaskParams{
			Taskid:        v,
			Cumulativegas: 0.0000000001,
		})
		assert.NoError(t, err, "err in AddTasks")
	}
	log.Println("AddTask time:", time.Since(start))*/
	tx, err := sqlite.Begin()
	assert.NoError(t, err, "err in Begin")
	qtx := queries.WithTx(tx)
	start := time.Now()
	for _, v := range addedTasks {
		err = qtx.AddTask(ctx, AddTaskParams{
			Taskid:        v,
			Cumulativegas: 0.0000000001,
		})
		assert.NoError(t, err, "err in AddTasks")

	}
	tx.Commit()
	tx.Rollback()
	log.Println("AddTask time:", time.Since(start))

	/*
		//2024/05/05 00:19:20 AddTask time: 866.9055ms
		err = queries.AddTasksNew(ctx, AddTasksParams2{
			Taskids:       addedTasks,
			Cumulativegas: 0.0000000001,
		})
		assert.NoError(t, err, "err in AddTasks")
		log.Println("AddTask time:", time.Since(start))
	*/

	queued, err := queries.TotalQueuedTasks(ctx)
	assert.NoError(t, err, "err in TotalQueuedTasks")
	log.Println(queued, len(addedTasks))

	solvedTasks := make([]task.TaskId, 0)

	start = time.Now()

	for i := 0; i < totalTasksToSolve; i++ {
		taskId, err := popTaskFromQueue(ctx, sqlite, queries)
		assert.NoError(t, err, "err in popTaskFromQueue")

		solvedTasks = append(solvedTasks, taskId)
	}
	log.Println("popTaskFromQueue time:", time.Since(start))

	tx, err = sqlite.Begin()
	assert.NoError(t, err, "err in Begin")
	qtx = queries.WithTx(tx)
	start = time.Now()

	for _, v := range solvedTasks {
		commitment := generateRandom32Bytes()
		// create an author
		//start := time.Now()

		err := qtx.CreateCommitment(ctx, CreateCommitmentParams{
			Taskid:     v,
			Commitment: commitment,
			Validator:  validator1,
		})
		assert.NoError(t, err, "err in CreateCommitment")

		// err = qtx.UpdateTaskStatus(ctx, UpdateTaskStatusParams{
		// 	Taskid: v,
		// 	Status: 2,
		// })
		// assert.NoError(t, err, "UpdateTaskStatus")

	}
	tx.Commit()
	tx.Rollback()
	log.Println("CreateCommitment time:", time.Since(start))

	tx, err = sqlite.Begin()
	assert.NoError(t, err, "err in Begin")
	qtx = queries.WithTx(tx)
	start = time.Now()
	for _, v := range solvedTasks {
		cid := generateRandomBytes(160)

		// create an author
		err := qtx.CreateSolution(ctx, CreateSolutionParams{
			Taskid:    v,
			Cid:       cid,
			Validator: validator1,
		})
		assert.NoError(t, err, "err in CreateSolution")
		//log.Println(insertedSol.Added)
	}
	tx.Commit()
	tx.Rollback()
	log.Println("CreateSolution time:", time.Since(start))

	commies, err := queries.GetCommitmentBatch(ctx, 500)
	assert.NoError(t, err, "err in GetCommitmentBatch")

	tx, err = sqlite.Begin()
	assert.NoError(t, err, "err in Begin")

	qtx = queries.WithTx(tx)
	start = time.Now()
	for _, v := range commies {
		log.Println(v.Taskid, v.Commitment, v.Validator)
		//err := qtx.UpdateTaskCommitment(ctx, v.Taskid)
		err = qtx.UpdateTaskStatus(ctx, UpdateTaskStatusParams{
			Taskid: v.Taskid,
			Status: 2,
		})
		assert.NoError(t, err, "UpdateTaskStatus")
		assert.NoError(t, err, "err in UpdateTaskCommitment")
	}
	tx.Commit()
	tx.Rollback()

	log.Println("UpdateTaskCommitment time:", time.Since(start))

	sols, err := queries.GetSolutionBatch(ctx, GetSolutionBatchParams{
		Validator: validator1,
		Limit:     247,
	})
	assert.NoError(t, err, "err in GetSolutionBatch")

	tx, err = sqlite.Begin()
	assert.NoError(t, err, "err in Begin")

	qtx = queries.WithTx(tx)
	start = time.Now()

	claimTime := time.Now().Add(2000 * time.Second)

	for _, v := range sols {
		err := qtx.UpdateTaskSolution(ctx, UpdateTaskSolutionParams{
			Taskid:    v.Taskid,
			Claimtime: claimTime.Unix(),
		})
		assert.NoError(t, err, "err in UpdateTaskCommitment")

		log.Println(v.Taskid)
	}

	tx.Commit()
	tx.Rollback()
	log.Println("UpdateTaskSols time:", time.Since(start))

	claims, err := queries.GetTasksByLowestCost(ctx, GetTasksByLowestCostParams{
		Claimtime: time.Now().Add(2100 * time.Second).Unix(),
		Limit:     600,
	})
	assert.NoError(t, err, "err in GetTasksByLowestCost")

	start = time.Now()

	/*	for _, v := range claims {

		err = deleteTasks(ctx, sqlite, queries, v.Taskid)
		assert.NoError(t, err, "err in deleteTasks")

	}*/

	err = deleteTasks(ctx, sqlite, queries, claims)
	assert.NoError(t, err, "err in deleteTasks")

	log.Println("GetTasksByLowestCost time:", time.Since(start))

	log.Println("GetTasksByLowestCost :", len(claims))

}

func generateRandom32Bytes() [32]byte {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return b
}

func generateRandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)

	}
	return b
}

func TestPopTask(t *testing.T) {
	ctx := context.Background()

	sqlite := setup(t)

	queries := New(sqlite)

	totalTasksToMake := 1
	totalTasksToPop := 15

	addedTasks := make([]task.TaskId, 0)
	for i := 0; i < totalTasksToMake; i++ {
		taskid := generateRandom32Bytes()
		addedTasks = append(addedTasks, taskid)
	}

	tx, err := sqlite.Begin()
	assert.NoError(t, err, "err in Begin")
	qtx := queries.WithTx(tx)
	start := time.Now()
	for _, v := range addedTasks {
		err = qtx.AddTask(ctx, AddTaskParams{
			Taskid:        v,
			Cumulativegas: 0.0000000001,
		})
		assert.NoError(t, err, "err in AddTasks")

	}
	tx.Commit()
	tx.Rollback()
	log.Println("AddTask time:", time.Since(start))

	queued, err := queries.TotalQueuedTasks(ctx)
	assert.NoError(t, err, "err in TotalQueuedTasks")
	log.Println(queued, len(addedTasks))

	for i := 0; i < totalTasksToPop; i++ {
		start = time.Now()
		taskId, err := popTaskFromQueue(ctx, sqlite, queries)
		assert.NoError(t, err, "err in popTaskFromQueue")

		log.Println("popTaskFromQueue time:", time.Since(start), taskId)
	}
}
