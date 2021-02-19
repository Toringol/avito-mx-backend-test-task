package businessConnService

import "github.com/Toringol/avito-mx-backend-test-task/app/models"

type IUsecase interface {
	SelectProduct(int64, int64) (*models.ProductInfo, error)
	SelectProductsBySpecificProductInfo(*models.UserListRequest) ([]*models.ProductInfo, error)
	CreateProduct(*models.ProductInfo) (int64, error)
	UpdateProduct(*models.ProductInfo) (int64, error)
	DeleteProduct(int64, int64) (int64, error)
}
