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
		ExpectQuery("SELECT state FROM productUploadsTask WHERE task_id").
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
	if !reflect.DeepEqual(item, testState) {
		t.Errorf("results not match, want %v, have %v", testState, item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT state FROM productUploadsTask WHERE task_id").
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
		ExpectQuery("SELECT state FROM productUploadsTask WHERE task_id").
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
