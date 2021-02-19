package businessConnService

import "github.com/Toringol/avito-mx-backend-test-task/app/models"

type IUsecase interface {
	SelectProduct(int64, int64) (*models.ProductInfo, error)
	SelectProductsBySpecificProductInfo(*models.UserListRequest) ([]*models.ProductInfo, error)
	CreateProduct(*models.ProductInfo) (int64, error)
	UpdateProduct(*models.ProductInfo) (int64, error)
	DeleteProduct(int64, int64) (int64, error)

	CreateTask() (int64, error)
	UpdateTaskState(int64, string) (int64, error)

	SelectTaskStatsByTaskID(int64) (*models.TaskStats, error)
	CreateTaskStats(int64, *models.TaskStats) (int64, error)
}
