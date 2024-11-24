package db

import (
	"context"
	"strings"

	task "gobius/common"
)

const addTasks2 = `-- name: AddTasks :exec
INSERT INTO tasks(
  taskid, cumulativeGas
) VALUES /*replace*/
`

type AddTasksParams2 struct {
	Taskids       []task.TaskId
	Cumulativegas float64
}

func (q *Queries) AddTasksNew(ctx context.Context, arg AddTasksParams2) error {
	query := addTasks2
	var queryParams []interface{}
	var builder strings.Builder
	if len(arg.Taskids) > 0 {
		for i, v := range arg.Taskids {
			queryParams = append(queryParams, v, arg.Cumulativegas)
			builder.WriteString("(?, ?)")
			if i < len(arg.Taskids)-1 {
				builder.WriteString(", ")
			}
		}
		query = strings.Replace(query, "/*replace*/", builder.String(), 1)
	} else {
		query = strings.Replace(query, "/*replace*/", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

