package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
	"github.com/spf13/viper"
)

type repository struct {
	DB *sql.DB
}

func NewRepository() businessConnService.IRepository {
	host := viper.GetString("DBHost")
	port := viper.GetInt("DBPort")
	user := viper.GetString("DBUser")
	password := viper.GetString("DBPassword")
	dbname := viper.GetString("DBName")

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println(err)
		return nil
	}
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil
	}

	return &repository{
		DB: db,
	}
}

func (repo *repository) SelectProduct(sellerID, offerID int64) (*models.ProductInfo, error) {
	productInfo := new(models.ProductInfo)

	err := repo.DB.
		QueryRow("SELECT * FROM productsinfo WHERE seller_id = $1 AND offer_id = $2", sellerID, offerID).
		Scan(&productInfo.SellerID, &productInfo.OfferID, &productInfo.Name, &productInfo.Price,
			&productInfo.Quantity, &productInfo.Available)
	if err != nil {
		return nil, err
	}

	return productInfo, nil
}

func (repo *repository) SelectProductsBySpecificProductInfo(userListRequest *models.UserListRequest) ([]*models.ProductInfo, error) {
	products := []*models.ProductInfo{}

	rows := new(sql.Rows)
	var err error

	switch {
	case userListRequest.SellerID > 0 && userListRequest.OfferID > 0:
		rows, err = repo.DB.Query("SELECT * FROM productsinfo WHERE seller_id = $1 AND offer_id = $2",
			userListRequest.SellerID, userListRequest.OfferID)
	case userListRequest.SellerID > 0:
		rows, err = repo.DB.Query("SELECT * FROM productsinfo WHERE seller_id = $1", userListRequest.SellerID)
	case userListRequest.OfferID > 0:
		rows, err = repo.DB.Query("SELECT * FROM productsinfo WHERE offer_id = $1", userListRequest.OfferID)
	default:
		rows, err = repo.DB.Query("SELECT * FROM productsinfo")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := new(models.ProductInfo)

		err := rows.Scan(&product.SellerID, &product.OfferID, &product.Name, &product.Price,
			&product.Quantity, &product.Available)
		if err != nil {
			return nil, err
		}

		if strings.Contains(product.Name, userListRequest.Name) {
			products = append(products, product)
		}
	}

	return products, nil
}

func (repo *repository) CreateProduct(productInfo *models.ProductInfo) (int64, error) {
	res, err := repo.DB.Exec(
		"INSERT INTO productsinfo (seller_id, offer_id, name, price, quantity, available) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		productInfo.SellerID,
		productInfo.OfferID,
		productInfo.Name,
		productInfo.Price,
		productInfo.Quantity,
		productInfo.Available,
	)
	if err != nil {
		return 0, err
	}

	affectedRowsCounter, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRowsCounter, nil
}

func (repo *repository) UpdateProduct(productInfo *models.ProductInfo) (int64, error) {
	res, err := repo.DB.Exec(
		"UPDATE productsinfo SET name = $1, price = $2, quantity = $3 "+
			"WHERE seller_id = $4 AND offer_id = $5",
		productInfo.Name,
		productInfo.Price,
		productInfo.Quantity,
		productInfo.SellerID,
		productInfo.OfferID,
	)
	if err != nil {
		return 0, err
	}

	affectedRowsCounter, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRowsCounter, nil
}

func (repo *repository) DeleteProduct(sellerID, offerID int64) (int64, error) {
	res, err := repo.DB.Exec(
		"DELETE FROM productsinfo WHERE seller_id = $1 AND offer_id = $2",
		sellerID,
		offerID,
	)
	if err != nil {
		return 0, err
	}

	affectedRowsCounter, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRowsCounter, nil
}

func (repo *repository) CreateTask() (int64, error) {
	taskID := int64(0)
	stateDefault := "CREATED"

	err := repo.DB.QueryRow(
		"INSERT INTO productUploadsTask (state) VALUES ($1) RETURNING task_id",
		stateDefault,
	).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func (repo *repository) UpdateTaskState(taskID int64, state string) (int64, error) {
	res, err := repo.DB.Exec(
		"UPDATE productUploadsTask SET state = $1 WHERE task_id = $2",
		state,
		taskID,
	)
	if err != nil {
		return 0, err
	}

	affectedRowsCounter, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRowsCounter, nil
}

func (repo *repository) SelectTaskStatsByTaskID(taskID int64) (*models.TaskStats, error) {
	taskStats := new(models.TaskStats)

	err := repo.DB.
		QueryRow("SELECT * FROM productTaskStats WHERE task_id = $1", taskID).
		Scan(&taskStats.ProductsCreated, &taskStats.ProductsUpdated,
			&taskStats.ProductsDeleted, &taskStats.RowsWithErrors)
	if err != nil {
		return nil, err
	}

	return taskStats, nil
}

func (repo *repository) CreateTaskStats(taskID int64, taskStats *models.TaskStats) (int64, error) {
	res, err := repo.DB.Exec(
		"INSERT INTO productTaskStats "+
			"(task_id, products_created, products_updated, products_deleted, rows_with_errors) "+
			"VALUES ($1, $2, $3, $4, $5)",
		taskID,
		taskStats.ProductsCreated,
		taskStats.ProductsUpdated,
		taskStats.ProductsDeleted,
		taskStats.RowsWithErrors,
	)
	if err != nil {
		return 0, nil
	}

	affectedRowsCounter, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRowsCounter, nil
}
