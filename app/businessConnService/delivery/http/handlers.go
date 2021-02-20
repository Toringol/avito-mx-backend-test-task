package http

import (
	"database/sql"
	"encoding/json"
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
	usecase   businessConnService.IUsecase
	logger    *logrus.Logger
	taskQueue chan models.Task
}

func NewHandlers(us businessConnService.IUsecase, taskQueue chan models.Task, logger *logrus.Logger) *mux.Router {
	handlers := handlers{
		usecase:   us,
		taskQueue: taskQueue,
		logger:    logger,
	}

	r := mux.NewRouter()

	r.HandleFunc("/loadProduct",
		middlewares.LogRequestMiddleware(handlers.logger, handlers.handleLoadProduct)).
		Methods("POST")

	r.HandleFunc("/getProduct",
		middlewares.LogRequestMiddleware(handlers.logger, handlers.handleGetProducts)).
		Methods("GET")

	r.HandleFunc("/getTaskState/{task_id:[0-9]+}",
		middlewares.LogRequestMiddleware(handlers.logger, handlers.handleGetTaskState)).
		Methods("GET")

	r.HandleFunc("/getTaskStats/{task_id:[0-9]+}",
		middlewares.LogRequestMiddleware(handlers.logger, handlers.handleGetTaskStats)).
		Methods("GET")

	return r
}

func (h *handlers) handleLoadProduct(w http.ResponseWriter, r *http.Request) {
	sellerID := r.FormValue("seller_id")
	if sellerID == "" {
		h.logger.Info("Empty sellerID")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	taskID, err := h.usecase.CreateTask()
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	taskIDJSON, err := json.Marshal(taskID)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	task := models.Task{
		TaskID:   taskID,
		SellerID: sellerID,
		Files:    r.MultipartForm.File,
	}

	h.taskQueue <- task

	w.Header().Set("Content-Type", "application/json")
	w.Write(taskIDJSON)
}

func (h *handlers) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	userListRequest := new(models.UserListRequest)

	if err := json.NewDecoder(r.Body).Decode(userListRequest); err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Info("BadRequest")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	products, err := h.usecase.SelectProductsBySpecificProductInfo(userListRequest)
	switch {
	case err == sql.ErrNoRows:
		h.logger.WithField("ErrInfo", err.Error()).Info("No such products")
		w.Write([]byte("No such products"))
		return
	case err != nil:
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	f := excelize.NewFile()

	for i, product := range products {
		counterStr := strconv.Itoa(i + 1)

		err = f.SetCellValue("Sheet1", "A"+counterStr, product.OfferID)
		if err != nil {
			h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f.SetCellValue("Sheet1", "B"+counterStr, product.Name)
		if err != nil {
			h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f.SetCellValue("Sheet1", "C"+counterStr, product.Price)
		if err != nil {
			h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Quantity)
		if err != nil {
			h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = f.SetCellValue("Sheet1", "D"+counterStr, product.Available)
		if err != nil {
			h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	f.Write(w)
}

func (h *handlers) handleGetTaskState(w http.ResponseWriter, r *http.Request) {
	taskIDStr, ok := mux.Vars(r)["task_id"]
	if !ok {
		h.logger.WithField("TaskID", taskIDStr).Info("BadRequest")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	state, err := h.usecase.SelectTaskState(taskID)
	switch {
	case err == sql.ErrNoRows:
		h.logger.WithField("TaskID", taskID).Info("BadRequest no such task")
		http.Error(w, "No such task", http.StatusBadRequest)
		return
	case err != nil:
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	stateJSON, err := json.Marshal(state)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(stateJSON)
}

func (h *handlers) handleGetTaskStats(w http.ResponseWriter, r *http.Request) {
	taskIDStr, ok := mux.Vars(r)["task_id"]
	if !ok {
		h.logger.WithField("TaskID", taskIDStr).Info("BadRequest")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	stats, err := h.usecase.SelectTaskStatsByTaskID(taskID)
	switch {
	case err == sql.ErrNoRows:
		h.logger.WithField("TaskID", taskID).Info("BadRequest no such task")
		http.Error(w, "No such task", http.StatusBadRequest)
		return
	case err != nil:
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsJSON)
}
