package storage

import (
	task "gobius/common"
)

type TaskData struct {
	TaskId     task.TaskId
	Commitment [32]byte
	Solution   []byte
}

type TaskDataSlice []TaskData

func (t TaskDataSlice) GetCommitments() ([][32]byte, map[[32]byte]task.TaskId) {
	var commitments [][32]byte
	var commitmentsToTasks map[[32]byte]task.TaskId = make(map[[32]byte]task.TaskId)
	for _, taskData := range t {
		if taskData.Commitment != [32]byte{} {
			commitments = append(commitments, taskData.Commitment)
			commitmentsToTasks[taskData.Commitment] = taskData.TaskId
		}
	}
	return commitments, commitmentsToTasks
}

func (t TaskDataSlice) GetSolutions() ([][]byte, []task.TaskId) {
	var solutions [][]byte
	var tasks []task.TaskId

	for _, taskData := range t {
		if taskData.Solution != nil {
			solutions = append(solutions, taskData.Solution)
			tasks = append(tasks, taskData.TaskId)
		}
	}
	return solutions, tasks
}

type ClaimTask struct {
	ID        task.TaskId
	Time      int64
	TotalCost float64
}

type ClaimTaskSlice []ClaimTask

func (t ClaimTaskSlice) SplitIntoChunks(chunkSize int) []ClaimTaskSlice {
	var chunks []ClaimTaskSlice
	for i := 0; i < len(t); i += chunkSize {
		end := i + chunkSize
		if end > len(t) {
			end = len(t)
		}
		chunks = append(chunks, t[i:end])
	}
	return chunks
}
