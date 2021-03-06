package models

// TaskStats is model stats of loading files for user
// swagger:model TaskStats
type TaskStats struct {
	TaskID          int64 `json:"task_id"`
	ProductsCreated int64 `json:"products_created"`
	ProductsUpdated int64 `json:"products_updated"`
	ProductsDeleted int64 `json:"products_deleted"`
	RowsWithErrors  int64 `json:"rows_with_errors"`
}
