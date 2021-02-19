package usecase

import (
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
)

type usecase struct {
	repo businessConnService.IRepository
}

func NewUsecase(repo businessConnService.IRepository) businessConnService.IUsecase {
	return usecase{repo: repo}
}

func (us usecase) SelectProduct(sellerID, offerID int64) (*models.ProductInfo, error) {
	return us.repo.SelectProduct(sellerID, offerID)
}

func (us usecase) SelectProductsBySpecificProductInfo(userListRequest *models.UserListRequest) ([]*models.ProductInfo, error) {
	return us.repo.SelectProductsBySpecificProductInfo(userListRequest)
}

func (us usecase) CreateProduct(productInfo *models.ProductInfo) (int64, error) {
	return us.repo.CreateProduct(productInfo)
}

func (us usecase) UpdateProduct(productInfo *models.ProductInfo) (int64, error) {
	return us.repo.UpdateProduct(productInfo)
}

func (us usecase) DeleteProduct(sellerID, offerID int64) (int64, error) {
	return us.repo.DeleteProduct(sellerID, offerID)
}
