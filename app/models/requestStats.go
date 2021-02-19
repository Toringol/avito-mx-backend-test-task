package models

type RequestStats struct {
	ProductsCreated int64 `json:"products_created"`
	ProductsUpdated int64 `json:"products_updated"`
	ProductsDeleted int64 `json:"products_deleted"`
	RowsWithErrors  int64 `json:"rows_with_errors"`
}
