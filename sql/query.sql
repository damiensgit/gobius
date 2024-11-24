-- name: GetCommitments :many
SELECT * FROM commitments
ORDER BY added ASC LIMIT ?;

-- name: CreateCommitment :exec
INSERT INTO commitments (
  taskid, commitment, validator
) VALUES (
  ?, ?, ?
);

-- name: DeleteCommitment :exec
DELETE FROM commitments
WHERE taskid = ?;

-- name: CreateSolution :exec
INSERT INTO solutions (
  taskid, cid, validator
) VALUES (
  ?, ?, ?
);

-- name: DeleteSolution :exec
DELETE FROM solutions
WHERE taskid = ?;

-- name: GetTasksByLowestCost :many
SELECT * FROM tasks
WHERE status = 3 AND claimtime < ?
ORDER BY cumulativeGas ASC 
LIMIT ?;

-- name: AddTask :exec
INSERT INTO tasks(
  taskid, txhash, cumulativeGas
) VALUES (
  ?,?, ?
);

-- name: AddTasks :exec
INSERT INTO tasks(
  taskid, txhash, cumulativeGas
) VALUES (sqlc.slice('taskids'),?, ?);

-- name: UpdateTaskSolution :exec
UPDATE tasks
SET status = 3, claimtime = ?, cumulativeGas = cumulativeGas + ?
WHERE taskid = ?;

-- name: UpdateTaskGas :exec
UPDATE tasks
SET cumulativeGas = cumulativeGas + ?
WHERE taskid = ?;

-- name: GetCommitmentBatch :many
SELECT 
commitments.taskid, commitments.commitment, commitments.validator
FROM commitments
JOIN tasks ON commitments.taskid = tasks.taskid 
--WHERE tasks.committed = false
WHERE tasks.status = 1
ORDER BY commitments.added ASC 
LIMIT ?;

-- name: GetSolutionBatch :many
SELECT 
solutions.taskid, solutions.cid 
FROM solutions 
JOIN tasks ON solutions.taskid = tasks.taskid 
WHERE tasks.status = 2 AND solutions.validator = ?
ORDER BY solutions.added ASC 
LIMIT ?;

-- name: GetSolutions :many
SELECT 
solutions.taskid, solutions.cid 
FROM solutions 
JOIN tasks ON solutions.taskid = tasks.taskid;

-- name: GetQueuedTasks :many
SELECT 
taskid, txhash
FROM tasks 
WHERE status = 0;

-- name: TotalPendingTasks :one
SELECT 
count(taskid)
FROM tasks 
WHERE status = 0;

-- name: TotalCommitments :one
SELECT 
count(commitments.taskid)
FROM commitments
JOIN tasks ON commitments.taskid = tasks.taskid 
WHERE tasks.status = 1;

-- name: TotalSolutionsAndClaims :one
SELECT 
    count(CASE WHEN tasks.status = 2 AND solutions.taskid IS NOT NULL THEN 1 END) AS total_solutions,
    count(CASE WHEN tasks.status = 3 AND claimtime > 0 THEN 1 END) AS total_claims
FROM tasks 
LEFT JOIN solutions ON solutions.taskid = tasks.taskid 
WHERE tasks.status IN (2, 3);

-- name: UpdateTaskStatus :execrows
UPDATE tasks SET status = ? WHERE taskid = ?;

-- name: DeletedClaimedTask :execrows
DELETE FROM tasks WHERE taskid = ? AND status = 3;

-- name: DeletedCommitment :execrows
DELETE FROM commitments WHERE taskid = ?;

-- name: DeletedSolution :execrows
DELETE FROM solutions WHERE taskid = ?;

-- name: UpdateTaskStatusAndGas :execrows
UPDATE tasks
SET cumulativeGas = cumulativeGas + ?, status = ?
WHERE taskid = ?;

-- name: SetTaskQueuedStatus :execrows
UPDATE tasks SET status = 1 WHERE taskid = ? and status = 0;

-- name: PopTask :one
UPDATE tasks
SET status = 1
WHERE taskid = (SELECT taskid
FROM tasks
WHERE status = 0
LIMIT 1)
RETURNING taskid, txhash;

-- name: GetPendingSolutionsCountPerValidator :many
SELECT 
    solutions.validator,
    COUNT(solutions.taskid) AS solution_count
FROM solutions 
JOIN tasks ON solutions.taskid = tasks.taskid 
WHERE tasks.status = 2
GROUP BY solutions.validator
ORDER BY solution_count DESC;