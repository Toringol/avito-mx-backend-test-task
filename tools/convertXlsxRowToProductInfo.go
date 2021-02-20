package tools

import (
	"errors"
	"strconv"

	"github.com/Toringol/avito-mx-backend-test-task/app/models"
)

func ConvertXlsxRowToProductInfo(row []string, sellerIDStr string) (*models.ProductInfo, error) {
	if len(row) != 5 {
		return nil, errors.New("Xlsx row has wrong length")
	}

	offerIDStr := row[0]
	nameStr := row[1]
	priceStr := row[2]
	quantityStr := row[3]
	availableStr := row[4]

	if offerIDStr == "" || nameStr == "" || priceStr == "" ||
		quantityStr == "" || availableStr == "" {
		return nil, errors.New("Nil col value")
	}

	productInfo := new(models.ProductInfo)

	sellerID, err := strconv.ParseInt(sellerIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	offerID, err := strconv.ParseInt(offerIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	name := nameStr

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, err
	}

	quantity, err := strconv.ParseInt(quantityStr, 10, 64)
	if err != nil {
		return nil, err
	}

	available, err := strconv.ParseBool(availableStr)
	if err != nil {
		return nil, err
	}

	if offerID <= 0 || sellerID <= 0 || price < 0 || quantity < 0 {
		return nil, errors.New("Misrepresentation of values")
	}

	productInfo.SellerID = sellerID
	productInfo.OfferID = offerID
	productInfo.Name = name
	productInfo.Price = price
	productInfo.Quantity = quantity
	productInfo.Available = available

	return productInfo, nil
}
