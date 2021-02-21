package repository

import (
	"fmt"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
)

func TestSelectProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"seller_id", "offer_id", "name", "price", "quantity", "available"})

	preparedData := []*models.ProductInfo{
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

	testSellerID := int64(1)
	testOfferID := int64(1)

	for _, item := range preparedData {
		rows = rows.AddRow(item.SellerID, item.OfferID, item.Name, item.Price, item.Quantity, item.Available)
	}

	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnRows(rows)

	repo := &repository{
		DB: db,
	}

	item, err := repo.SelectProduct(testSellerID, testOfferID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, preparedData[0]) {
		t.Errorf("results not match, want %v, have %v", preparedData[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectProduct(testSellerID, testOfferID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"seller_id", "offer_id"}).
		AddRow(1, 2)

	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnRows(rows)

	_, err = repo.SelectProduct(testSellerID, testOfferID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestSelectProductsBySpecificProductInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"seller_id", "offer_id", "name", "price", "quantity", "available"})

	preparedData := []*models.ProductInfo{
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

	testSellerID := int64(1)
	testOfferID := int64(1)

	userListRequest := &models.UserListRequest{
		SellerID: testSellerID,
		OfferID:  testOfferID,
		Name:     "теле",
	}

	for _, item := range preparedData {
		rows = rows.AddRow(item.SellerID, item.OfferID, item.Name, item.Price, item.Quantity, item.Available)
	}

	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnRows(rows)

	repo := &repository{
		DB: db,
	}

	item, err := repo.SelectProductsBySpecificProductInfo(userListRequest)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item[0], preparedData[0]) && !reflect.DeepEqual(item[1], preparedData[1]) {
		t.Errorf("results not match, want %v, have %v", preparedData[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectProductsBySpecificProductInfo(userListRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"seller_id", "offer_id"}).
		AddRow(1, 2)

	mock.
		ExpectQuery("SELECT (.+) FROM productsinfo WHERE").
		WithArgs(testSellerID, testOfferID).
		WillReturnRows(rows)

	_, err = repo.SelectProductsBySpecificProductInfo(userListRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestSelectTaskState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"task_id", "state"})

	testTaskID := int64(1)
	testState := "CREATED"

	rows = rows.AddRow(testTaskID, testState)

	mock.
		ExpectQuery("SELECT (.+) FROM productUploadsTask WHERE task_id").
		WithArgs(testTaskID).
		WillReturnRows(rows)

	repo := &repository{
		DB: db,
	}

	item, err := repo.SelectTaskState(testTaskID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item.State, testState) {
		t.Errorf("results not match, want %v, have %v", testState, item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT (.+) FROM productUploadsTask WHERE task_id").
		WithArgs(testTaskID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectTaskState(testTaskID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"task_id"}).
		AddRow(1)

	mock.
		ExpectQuery("SELECT (.+) FROM productUploadsTask WHERE task_id").
		WithArgs(testTaskID).
		WillReturnRows(rows)

	_, err = repo.SelectTaskState(testTaskID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestSelectTaskStatsByTaskID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"task_id", "products_created", "products_updated",
			"products_deleted", "rows_with_errors"})

	preparedData := []*models.TaskStats{
		{
			TaskID:          1,
			ProductsCreated: 5,
			ProductsUpdated: 2,
			ProductsDeleted: 1,
			RowsWithErrors:  0,
		},
		{
			TaskID:          2,
			ProductsCreated: 8,
			ProductsUpdated: 1,
			ProductsDeleted: 2,
			RowsWithErrors:  2,
		},
	}

	for _, item := range preparedData {
		rows = rows.AddRow(item.TaskID, item.ProductsCreated, item.ProductsUpdated,
			item.ProductsDeleted, item.RowsWithErrors)
	}

	testTaskStats := &models.TaskStats{
		TaskID:          1,
		ProductsCreated: 5,
		ProductsUpdated: 2,
		ProductsDeleted: 1,
		RowsWithErrors:  0,
	}

	mock.
		ExpectQuery("SELECT (.+) FROM productTaskStats WHERE task_id").
		WithArgs(testTaskStats.TaskID).
		WillReturnRows(rows)

	repo := &repository{
		DB: db,
	}

	item, err := repo.SelectTaskStatsByTaskID(testTaskStats.TaskID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, testTaskStats) {
		t.Errorf("results not match, want %v, have %v", testTaskStats, item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT (.+) FROM productTaskStats WHERE task_id").
		WithArgs(testTaskStats.TaskID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectTaskStatsByTaskID(testTaskStats.TaskID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"task_id"}).
		AddRow(1)

	mock.
		ExpectQuery("SELECT (.+) FROM productTaskStats WHERE task_id").
		WithArgs(testTaskStats.TaskID).
		WillReturnRows(rows)

	_, err = repo.SelectTaskStatsByTaskID(testTaskStats.TaskID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &repository{
		DB: db,
	}

	preparedProductInfo := &models.ProductInfo{
		SellerID:  1,
		OfferID:   1,
		Name:      "tele",
		Price:     37.5,
		Quantity:  10,
		Available: true,
	}

	mock.
		ExpectExec("INSERT INTO productsinfo").
		WithArgs(preparedProductInfo.SellerID, preparedProductInfo.OfferID, preparedProductInfo.Name,
			preparedProductInfo.Price, preparedProductInfo.Quantity, preparedProductInfo.Available).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected, err := repo.CreateProduct(preparedProductInfo)
	if rowsAffected != 1 {
		t.Errorf("bad rowsAffected: want %v, have %v", 1, rowsAffected)
		return
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO productsinfo`).
		WithArgs(preparedProductInfo.SellerID, preparedProductInfo.OfferID).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.CreateProduct(preparedProductInfo)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO productsinfo`).
		WithArgs(preparedProductInfo.SellerID, preparedProductInfo.OfferID, preparedProductInfo.Name,
			preparedProductInfo.Price, preparedProductInfo.Quantity, preparedProductInfo.Available).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.CreateProduct(preparedProductInfo)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &repository{
		DB: db,
	}

	testTaskID := int64(0)
	testState := "CREATED"

	mock.
		ExpectExec("INSERT INTO productUploadsTask").
		WithArgs(testState).
		WillReturnResult(sqlmock.NewResult(1, 1))

	taskID, err := repo.CreateTask()
	if taskID != testTaskID {
		t.Errorf("bad taskID: want %v, have %v", testTaskID, taskID)
		return
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO productUploadsTask`).
		WithArgs().
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.CreateTask()
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO productUploadsTask`).
		WithArgs(testState).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.CreateTask()
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestCreateTaskStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &repository{
		DB: db,
	}

	preparedProductInfo := &models.TaskStats{
		TaskID:          1,
		ProductsCreated: 5,
		ProductsUpdated: 2,
		ProductsDeleted: 1,
		RowsWithErrors:  0,
	}

	mock.
		ExpectExec("INSERT INTO productTaskStats").
		WithArgs(preparedProductInfo.TaskID, preparedProductInfo.ProductsCreated, preparedProductInfo.ProductsUpdated,
			preparedProductInfo.ProductsDeleted, preparedProductInfo.RowsWithErrors).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected, err := repo.CreateTaskStats(preparedProductInfo)
	if rowsAffected != 1 {
		t.Errorf("bad rowsAffected: want %v, have %v", 1, rowsAffected)
		return
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO productTaskStats`).
		WithArgs(preparedProductInfo.ProductsCreated, preparedProductInfo.ProductsUpdated).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.CreateTaskStats(preparedProductInfo)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO productTaskStats`).
		WithArgs(preparedProductInfo.TaskID, preparedProductInfo.ProductsCreated, preparedProductInfo.ProductsUpdated,
			preparedProductInfo.ProductsDeleted, preparedProductInfo.RowsWithErrors).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.CreateTaskStats(preparedProductInfo)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"seller_id", "offer_id", "name", "price", "quantity", "available"})

	preparedData := []*models.ProductInfo{
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

	expectData := &models.ProductInfo{
		SellerID:  1,
		OfferID:   1,
		Name:      "телефон",
		Price:     150.5,
		Quantity:  15,
		Available: true,
	}

	for _, item := range preparedData {
		rows = rows.AddRow(item.SellerID, item.OfferID, item.Name, item.Price, item.Quantity, item.Available)
	}

	mock.
		ExpectExec(`UPDATE productsinfo SET`).
		WithArgs(expectData.Name, expectData.Price, expectData.Quantity, expectData.Available,
			expectData.SellerID, expectData.OfferID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &repository{
		DB: db,
	}

	rowsAffected, err := repo.UpdateProduct(expectData)
	if rowsAffected != 0 {
		t.Errorf("bad rowsAffected: want %v, have %v", rowsAffected, 1)
		return
	}

	// query error
	mock.
		ExpectExec(`UPDATE productsinfo SET`).
		WithArgs(expectData.Name, expectData.Price, expectData.Quantity, expectData.Available).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.UpdateProduct(expectData)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.
		ExpectExec(`UPDATE productsinfo SET`).
		WithArgs(expectData.Name, expectData.Price, expectData.Quantity, expectData.Available,
			expectData.SellerID, expectData.OfferID).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.UpdateProduct(expectData)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestUpdateTaskState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"task_id", "state"})

	testTaskID := int64(1)
	testState := "CREATED"

	rows = rows.AddRow(testTaskID, testState)

	expectTaskID := int64(1)
	expectTestState := "IN PROGRESS"

	mock.
		ExpectExec(`UPDATE productUploadsTask SET`).
		WithArgs(expectTaskID, expectTestState).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &repository{
		DB: db,
	}

	rowsAffected, err := repo.UpdateTaskState(expectTaskID, expectTestState)
	if rowsAffected != 0 {
		t.Errorf("bad rowsAffected: want %v, have %v", rowsAffected, 1)
		return
	}

	// query error
	mock.
		ExpectExec(`UPDATE productUploadsTask SET`).
		WithArgs(expectTaskID).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.UpdateTaskState(expectTaskID, expectTestState)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.
		ExpectExec(`UPDATE productUploadsTask SET`).
		WithArgs(expectTaskID, expectTestState).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.UpdateTaskState(expectTaskID, expectTestState)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can`t create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"seller_id", "offer_id", "name", "price", "quantity", "available"})

	preparedData := []*models.ProductInfo{
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

	testData := &models.ProductInfo{
		SellerID: 1,
		OfferID:  1,
	}

	for _, item := range preparedData {
		rows = rows.AddRow(item.SellerID, item.OfferID, item.Name, item.Price, item.Quantity, item.Available)
	}

	mock.
		ExpectExec(`DELETE FROM productsinfo WHERE`).
		WithArgs(testData.SellerID, testData.OfferID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &repository{
		DB: db,
	}

	rowsAffected, err := repo.DeleteProduct(testData.SellerID, testData.OfferID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if rowsAffected != 1 {
		t.Errorf("bad rowsAffected: want %v, have %v", rowsAffected, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`DELETE FROM productsinfo WHERE`).
		WithArgs().
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.DeleteProduct(testData.SellerID, testData.OfferID)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.
		ExpectExec(`DELETE FROM productsinfo WHERE`).
		WithArgs(testData.SellerID, testData.OfferID).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.DeleteProduct(testData.SellerID, testData.OfferID)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
