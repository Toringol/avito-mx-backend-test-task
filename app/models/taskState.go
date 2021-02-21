package models

// TaskState is model for observe state of task
// swagger:model TaskState
type TaskState struct {
	TaskID int64  `json:"task_id"`
	State  string `json:"state"`
}
