package http

import (
	"bytes"
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetTaskState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test expect behavior

	usecase := businessConnService.NewMockIUsecase(ctrl)

	expectedData := &models.TaskState{
		TaskID: 1,
		State:  "IN PROGRESS",
	}

	usecase.EXPECT().SelectTaskState(expectedData.TaskID).Return(expectedData, nil)

	outputJSON := `{"task_id":1,"state":"IN PROGRESS"}`

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodGet, "/getTaskState/", nil)

	qParams := map[string]string{
		"task_id": "1",
	}

	request = mux.SetURLVars(request, qParams)

	response := httptest.NewRecorder()

	handlers.handleGetTaskState(response, request)

	if assert.Equal(t, http.StatusOK, response.Code) {
		assert.Equal(t, outputJSON, strings.Trim(response.Body.String(), "\n"))
	}

	// test error bad request (nil query)

	badRequestNilQuery := httptest.NewRequest(http.MethodGet, "/getTaskState/", nil)

	responseBadRequestNilQuery := httptest.NewRecorder()

	handlers.handleGetTaskState(responseBadRequestNilQuery, badRequestNilQuery)

	assert.Equal(t, http.StatusBadRequest, responseBadRequestNilQuery.Code)

	// test error bad request (bad query)

	badRequestIncorrectQuery := httptest.NewRequest(http.MethodGet, "/getTaskState/", nil)

	qParams = map[string]string{
		"task_id": "1.5",
	}

	badRequestIncorrectQuery = mux.SetURLVars(badRequestIncorrectQuery, qParams)

	responseBadRequestIncorrectQuery := httptest.NewRecorder()

	handlers.handleGetTaskState(responseBadRequestIncorrectQuery, badRequestIncorrectQuery)

	assert.Equal(t, http.StatusBadRequest, responseBadRequestIncorrectQuery.Code)

	// test DB return error

	usecase.EXPECT().SelectTaskState(expectedData.TaskID).Return(expectedData, errors.New("DB error"))

	requestDBError := httptest.NewRequest(http.MethodGet, "/getTaskState/", nil)

	qParams = map[string]string{
		"task_id": "1",
	}

	requestDBError = mux.SetURLVars(requestDBError, qParams)

	responseDBError := httptest.NewRecorder()

	handlers.handleGetTaskState(responseDBError, requestDBError)

	assert.Equal(t, http.StatusInternalServerError, responseDBError.Code)

	// test DB return sql.NoRows (Bad Request)

	usecase.EXPECT().SelectTaskState(expectedData.TaskID).Return(expectedData, sql.ErrNoRows)

	requestDBNoRows := httptest.NewRequest(http.MethodGet, "/getTaskState/", nil)

	qParams = map[string]string{
		"task_id": "1",
	}

	requestDBNoRows = mux.SetURLVars(requestDBNoRows, qParams)

	responseDBNoRows := httptest.NewRecorder()

	handlers.handleGetTaskState(responseDBNoRows, requestDBNoRows)

	assert.Equal(t, http.StatusBadRequest, responseDBNoRows.Code)
}

func TestHandleGetTaskStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test expect behavior

	usecase := businessConnService.NewMockIUsecase(ctrl)

	expectedData := &models.TaskStats{
		TaskID:          1,
		ProductsCreated: 5,
		ProductsUpdated: 2,
		ProductsDeleted: 1,
		RowsWithErrors:  0,
	}

	usecase.EXPECT().SelectTaskStatsByTaskID(expectedData.TaskID).Return(expectedData, nil)

	outputJSON := `{"task_id":1,"products_created":5,"products_updated":2,"products_deleted":1,"rows_with_errors":0}`

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodGet, "/getTaskStats/", nil)

	qParams := map[string]string{
		"task_id": "1",
	}

	request = mux.SetURLVars(request, qParams)

	response := httptest.NewRecorder()

	handlers.handleGetTaskStats(response, request)

	if assert.Equal(t, http.StatusOK, response.Code) {
		assert.Equal(t, outputJSON, strings.Trim(response.Body.String(), "\n"))
	}

	/// test error bad request (nil query)

	badRequestNilQuery := httptest.NewRequest(http.MethodGet, "/getTaskStats/", nil)

	responseBadRequestNilQuery := httptest.NewRecorder()

	handlers.handleGetTaskStats(responseBadRequestNilQuery, badRequestNilQuery)

	assert.Equal(t, http.StatusBadRequest, responseBadRequestNilQuery.Code)

	// test error bad request (bad query)

	badRequestIncorrectQuery := httptest.NewRequest(http.MethodGet, "/getTaskStats/", nil)

	qParams = map[string]string{
		"task_id": "1.5",
	}

	badRequestIncorrectQuery = mux.SetURLVars(badRequestIncorrectQuery, qParams)

	responseBadRequestIncorrectQuery := httptest.NewRecorder()

	handlers.handleGetTaskStats(responseBadRequestIncorrectQuery, badRequestIncorrectQuery)

	assert.Equal(t, http.StatusBadRequest, responseBadRequestIncorrectQuery.Code)

	// test DB return error

	usecase.EXPECT().SelectTaskStatsByTaskID(expectedData.TaskID).Return(expectedData, errors.New("DB error"))

	requestDBError := httptest.NewRequest(http.MethodGet, "/getTaskStats/", nil)

	qParams = map[string]string{
		"task_id": "1",
	}

	requestDBError = mux.SetURLVars(requestDBError, qParams)

	responseDBError := httptest.NewRecorder()

	handlers.handleGetTaskStats(responseDBError, requestDBError)

	assert.Equal(t, http.StatusInternalServerError, responseDBError.Code)

	// test DB return sql.NoRows (Bad Request)

	usecase.EXPECT().SelectTaskStatsByTaskID(expectedData.TaskID).Return(expectedData, sql.ErrNoRows)

	requestDBNoRows := httptest.NewRequest(http.MethodGet, "/getTaskStats/", nil)

	qParams = map[string]string{
		"task_id": "1",
	}

	requestDBNoRows = mux.SetURLVars(requestDBNoRows, qParams)

	responseDBNoRows := httptest.NewRecorder()

	handlers.handleGetTaskStats(responseDBNoRows, requestDBNoRows)

	assert.Equal(t, http.StatusBadRequest, responseDBNoRows.Code)

}

func TestHandleGetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test correct behavior

	usecase := businessConnService.NewMockIUsecase(ctrl)

	inputData := &models.UserListRequest{
		SellerID: 1,
		OfferID:  0,
		Name:     "теле",
	}

	inputDataJSON := `{"seller_id":1,"offer_id":0,"name":"теле"}`

	expectedData := []*models.ProductInfo{
		{
			SellerID:  1,
			OfferID:   1,
			Name:      "телефон",
			Price:     100.25,
			Quantity:  10,
			Available: true,
		},
		{
			SellerID:  1,
			OfferID:   2,
			Name:      "телевизор",
			Price:     57.6,
			Quantity:  15,
			Available: true,
		},
	}

	usecase.EXPECT().SelectProductsBySpecificProductInfo(inputData).Return(expectedData, nil)

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodGet, "/getProduct", strings.NewReader(inputDataJSON))

	response := httptest.NewRecorder()

	handlers.handleGetProducts(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	// test incorrect input data

	requestIncorrectInput := httptest.NewRequest(http.MethodGet, "/getProduct", nil)

	responseIncorrectInput := httptest.NewRecorder()

	handlers.handleGetProducts(responseIncorrectInput, requestIncorrectInput)

	assert.Equal(t, http.StatusBadRequest, responseIncorrectInput.Code)

	// test DB return error

	usecase.EXPECT().SelectProductsBySpecificProductInfo(inputData).Return(expectedData, errors.New("DB error"))

	requestDBError := httptest.NewRequest(http.MethodGet, "/getProduct", strings.NewReader(inputDataJSON))

	responseDBError := httptest.NewRecorder()

	handlers.handleGetProducts(responseDBError, requestDBError)

	assert.Equal(t, http.StatusInternalServerError, responseDBError.Code)

	// test DB return sql.NoRows (Bad Request)

	usecase.EXPECT().SelectProductsBySpecificProductInfo(inputData).Return(expectedData, sql.ErrNoRows)

	requestDBNoRows := httptest.NewRequest(http.MethodGet, "/getProduct", strings.NewReader(inputDataJSON))

	responseDBNoRows := httptest.NewRecorder()

	handlers.handleGetProducts(responseDBNoRows, requestDBNoRows)

	assert.Equal(t, http.StatusBadRequest, responseDBNoRows.Code)
}

func TestHandleLoadProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test correct behavior

	usecase := businessConnService.NewMockIUsecase(ctrl)

	testTaskID := int64(1)

	usecase.EXPECT().CreateTask().Return(testTaskID, nil)

	outputJSON := "1"

	loadData := []*models.ProductInfo{
		{
			SellerID:  1,
			OfferID:   1,
			Name:      "телефон",
			Price:     100.25,
			Quantity:  10,
			Available: true,
		},
		{
			SellerID:  1,
			OfferID:   2,
			Name:      "телевизор",
			Price:     57.6,
			Quantity:  15,
			Available: true,
		},
	}

	f := excelize.NewFile()

	for i, product := range loadData {
		counterStr := strconv.Itoa(i + 1)

		err := f.SetCellValue("Sheet1", "A"+counterStr, product.OfferID)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "B"+counterStr, product.Name)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "C"+counterStr, product.Price)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Quantity)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Available)
		assert.NoError(t, err)
	}

	correctBody := &bytes.Buffer{}

	writer := multipart.NewWriter(correctBody)

	part, err := writer.CreateFormFile("products", "testFile.xlsx")
	assert.NoError(t, err)

	err = f.Write(part)
	assert.NoError(t, err)

	err = writer.WriteField("seller_id", "1")
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodPost, "/loadProduct", correctBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()

	go func() {
		<-handlers.taskQueue
	}()

	handlers.handleLoadProduct(response, request)

	if assert.Equal(t, http.StatusOK, response.Code) {
		assert.Equal(t, outputJSON, strings.Trim(response.Body.String(), "\n"))
	}

	// test incorrect input data (empty data)

	requestIncorrectInput := httptest.NewRequest(http.MethodPost, "/loadProduct", nil)

	responseIncorrectInput := httptest.NewRecorder()

	handlers.handleLoadProduct(responseIncorrectInput, requestIncorrectInput)

	assert.Equal(t, http.StatusBadRequest, responseIncorrectInput.Code)
}

func TestHandleLoadProductIncorrectData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := businessConnService.NewMockIUsecase(ctrl)

	loadData := []*models.ProductInfo{
		{
			SellerID:  1,
			OfferID:   1,
			Name:      "телефон",
			Price:     100.25,
			Quantity:  10,
			Available: true,
		},
		{
			SellerID:  1,
			OfferID:   2,
			Name:      "телевизор",
			Price:     57.6,
			Quantity:  15,
			Available: true,
		},
	}

	f := excelize.NewFile()

	for i, product := range loadData {
		counterStr := strconv.Itoa(i + 1)

		err := f.SetCellValue("Sheet1", "A"+counterStr, product.OfferID)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "B"+counterStr, product.Name)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "C"+counterStr, product.Price)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Quantity)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Available)
		assert.NoError(t, err)
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("products", "testFile.xlsx")
	assert.NoError(t, err)

	err = f.Write(part)
	assert.NoError(t, err)

	err = writer.WriteField("seller_id", "1.5")
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodPost, "/loadProduct", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()

	handlers.handleLoadProduct(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandleLoadProductDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := businessConnService.NewMockIUsecase(ctrl)

	testTaskID := int64(1)

	usecase.EXPECT().CreateTask().Return(testTaskID, errors.New("DB error"))

	loadData := []*models.ProductInfo{
		{
			SellerID:  1,
			OfferID:   1,
			Name:      "телефон",
			Price:     100.25,
			Quantity:  10,
			Available: true,
		},
		{
			SellerID:  1,
			OfferID:   2,
			Name:      "телевизор",
			Price:     57.6,
			Quantity:  15,
			Available: true,
		},
	}

	f := excelize.NewFile()

	for i, product := range loadData {
		counterStr := strconv.Itoa(i + 1)

		err := f.SetCellValue("Sheet1", "A"+counterStr, product.OfferID)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "B"+counterStr, product.Name)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "C"+counterStr, product.Price)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Quantity)
		assert.NoError(t, err)

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Available)
		assert.NoError(t, err)
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("products", "testFile.xlsx")
	assert.NoError(t, err)

	err = f.Write(part)
	assert.NoError(t, err)

	err = writer.WriteField("seller_id", "1")
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	handlers := &handlers{
		usecase:   usecase,
		taskQueue: make(chan models.Task),
		logger:    logrus.New(),
	}

	request := httptest.NewRequest(http.MethodPost, "/loadProduct", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()

	handlers.handleLoadProduct(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
