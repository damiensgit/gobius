// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	common "github.com/ethereum/go-ethereum/common"
	task "gobius/common"
)

type Commitment struct {
	Taskid     task.TaskId
	Commitment task.TaskId
	Validator  common.Address
	Added      time.Time
}

type Solution struct {
	Taskid    task.TaskId
	Cid       []byte
	Validator common.Address
	Added     time.Time
}

type Task struct {
	Taskid        task.TaskId
	Txhash        common.Hash
	Cumulativegas float64
	Status        int64
	Claimtime     int64
}