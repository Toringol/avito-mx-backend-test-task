package tools

import (
	"testing"

	"github.com/Toringol/avito-mx-backend-test-task/app/models"
	"github.com/stretchr/testify/assert"
)

func TestConvertXlsxRowToProductInfo(t *testing.T) {
	testRowData := []string{"1", "a", "10.5", "10", "true"}
	testSellerID := int64(1)

	expectProductInfo := &models.ProductInfo{
		SellerID:  1,
		OfferID:   1,
		Name:      "a",
		Price:     10.5,
		Quantity:  10,
		Available: true,
	}

	productInfo, err := ConvertXlsxRowToProductInfo(testRowData, testSellerID)
	if assert.NoError(t, err) {
		assert.Equal(t, expectProductInfo, productInfo)
	}

	testRowDataIncorrectOfferID := []string{"1.5", "a", "10.5", "10", "true"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataIncorrectOfferID, testSellerID)
	assert.Error(t, err)

	testRowDataIncorrectPrice := []string{"1", "a", "abc", "10", "true"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataIncorrectPrice, testSellerID)
	assert.Error(t, err)

	testRowDataIncorrectQuantity := []string{"1", "a", "10.5", "10.5", "true"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataIncorrectQuantity, testSellerID)
	assert.Error(t, err)

	testRowDataIncorrectLength := []string{"1", "a", "-10.5", "10"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataIncorrectLength, testSellerID)
	assert.Error(t, err)

	testRowDataWithEmptyField := []string{"1", "a", "", "-10", "true"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataWithEmptyField, testSellerID)
	assert.Error(t, err)

	testRowDataNegativePrice := []string{"1", "a", "-10", "10", "true"}

	_, err = ConvertXlsxRowToProductInfo(testRowDataNegativePrice, testSellerID)
	assert.Error(t, err)
}
