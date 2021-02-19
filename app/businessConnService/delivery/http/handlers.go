package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/middlewares"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
)

type handlers struct {
	usecase businessConnService.IUsecase
}

func NewHandlers(us businessConnService.IUsecase) *mux.Router {
	handlers := handlers{usecase: us}

	r := mux.NewRouter()
	logger := logrus.New()

	r.HandleFunc("/loadProduct", middlewares.LogRequestMiddleware(logger, handlers.handleLoadProduct)).
		Methods("POST")

	r.HandleFunc("/getProduct", middlewares.LogRequestMiddleware(logger, handlers.handleGetProducts)).
		Methods("GET")

	return r
}

func (h *handlers) handleLoadProduct(w http.ResponseWriter, r *http.Request) {
	sellerID := r.FormValue("seller_id")
	if sellerID == "" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			fd, err := hdr.Open()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			f, err := excelize.OpenReader(fd)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			requestStats := new(models.RequestStats)

			sheets := f.GetSheetMap()
			for _, sheet := range sheets {
				rows, err := f.GetRows(sheet)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				for _, row := range rows {
					productInfo, err := convertXlsxToProductInfo(row, sellerID)
					if err != nil {
						requestStats.RowsWithErrors++

						fmt.Println(err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}

					productRecord, err := h.usecase.SelectProduct(productInfo.SellerID, productInfo.OfferID)
					switch {
					case err == sql.ErrNoRows && productInfo.Available:
						rowsAffected, err := h.usecase.CreateProduct(productInfo)
						if err != nil {
							fmt.Println(err)
							http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
							return
						}

						requestStats.ProductsCreated += rowsAffected
						continue
					case err != nil:
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}

					if !productInfo.Available {
						rowsAffected, err := h.usecase.DeleteProduct(productInfo.SellerID, productInfo.OfferID)
						if err != nil {
							fmt.Println(err)
							http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
							return
						}

						requestStats.ProductsDeleted += rowsAffected
					} else {
						productRecord.Name = productInfo.Name
						productRecord.Price = productInfo.Price
						productRecord.Quantity = productInfo.Quantity

						rowsAffected, err := h.usecase.UpdateProduct(productRecord)
						if err != nil {
							fmt.Println(err)
							http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
							return
						}

						requestStats.ProductsUpdated += rowsAffected
					}
				}
			}

			fmt.Println(requestStats)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handlers) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	userListRequest := new(models.UserListRequest)

	if err := json.NewDecoder(r.Body).Decode(userListRequest); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	products, err := h.usecase.SelectProductsBySpecificProductInfo(userListRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	f := excelize.NewFile()

	for i, product := range products {
		counterStr := strconv.Itoa(i + 1)
		f.SetCellValue("Sheet1", "A"+counterStr, product.OfferID)
		f.SetCellValue("Sheet1", "B"+counterStr, product.Name)
		f.SetCellValue("Sheet1", "C"+counterStr, product.Price)
		f.SetCellValue("Sheet1", "D"+counterStr, product.Quantity)
		f.SetCellValue("Sheet1", "D"+counterStr, product.Available)
	}

	f.Write(w)
}

func convertXlsxToProductInfo(row []string, sellerIDStr string) (*models.ProductInfo, error) {
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
