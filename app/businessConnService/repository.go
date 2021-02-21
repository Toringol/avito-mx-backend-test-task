package businessConnService

import (
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
)

type IRepository interface {
	SelectProduct(int64, int64) (*models.ProductInfo, error)
	SelectProductsBySpecificProductInfo(*models.UserListRequest) ([]*models.ProductInfo, error)
	CreateProduct(*models.ProductInfo) (int64, error)
	UpdateProduct(*models.ProductInfo) (int64, error)
	DeleteProduct(int64, int64) (int64, error)

	SelectTaskState(int64) (*models.TaskState, error)
	CreateTask() (int64, error)
	UpdateTaskState(int64, string) (int64, error)

	SelectTaskStatsByTaskID(int64) (*models.TaskStats, error)
	CreateTaskStats(*models.TaskStats) (int64, error)
}
