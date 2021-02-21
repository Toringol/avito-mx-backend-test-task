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

// NewHandlers - create new handlers using gorilla router
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

// swagger:operation POST /loadProduct handleLoadProduct
//
// Get sellerID and xlsx files and return task id
// ---
// produces:
// - multipart/form-data
// parameters:
// - name: seller_id
//   in: formData
//   description: The seller_id needs to match customer id with products.
//   required: true
//   type: text
// - name: products
//   in: formData
//   description: Files with products info.
//   required: true
//   type: file
// responses:
//   200:
//     description: successful operation
//     schema:
//       type: string
//       task_id: string
//       description: Return task id
//   400:
//       description: Invalid seller_id supplied
//   500:
//     description: Sth went wrong
func (h *handlers) handleLoadProduct(w http.ResponseWriter, r *http.Request) {
	sellerIDStr := r.FormValue("seller_id")
	if sellerIDStr == "" {
		h.logger.Info("Empty sellerID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sellerIDInt, err := strconv.ParseInt(sellerIDStr, 10, 64)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		SellerID: sellerIDInt,
		Files:    r.MultipartForm.File,
	}

	h.taskQueue <- task

	w.Header().Set("Content-Type", "application/json")
	w.Write(taskIDJSON)
}

// swagger:operation GET /getProduct handleGetProducts
//
// Get UserListRequest and return xlsx file with all products
// that match with request data
// ---
// produces:
// - application/json
// parameters:
// - name: userListRequest
//   in: json
//   description: userListRequest may contain seller_id, offer_id and name.
//   required: false
//   type: #/definitions/UserListRequest
// responses:
//   200:
//     description: successful operation
//     schema:
//       type: file
//       description: Return xlsx file
//   400:
//     description: Invalid userListRequest supplied
//   500:
//     description: Sth went wrong
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
		http.Error(w, "No such products", http.StatusBadRequest)
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

// swagger:operation GET /getTaskState/{task_id} handleGetTaskState
//
// Get task id and return state
// ---
// summary: Get task state by task id
// operationId: handleGetTaskState
// produces:
// - application/json
// parameters:
// - name: task_id
//   in: path
//   required: true
//   type: string
// responses:
//   200:
//     description: successful operation
//     schema:
//       type: string
//       state: string
//       description: Return state
//   400:
//     description: Invalid taskID supplied
//   500:
//     description: Sth went wrong
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	taskState, err := h.usecase.SelectTaskState(taskID)
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

	taskStateJSON, err := json.Marshal(taskState)
	if err != nil {
		h.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(taskStateJSON)
}

// swagger:operation GET /getTaskStats/{task_id} handleGetTaskStats
//
// Get task id and return stats
// ---
// summary: Get stats by task id
// operationId: handleGetTaskStats
// produces:
// - application/json
// parameters:
// - name: task_id
//   in: path
//   required: true
//   type: string
// responses:
//   200:
//     description: successful operation
//     schema:
//       $ref: '#/definitions/TaskStats'
//   400:
//     description: Invalid taskID supplied
//   500:
//     description: Sth went wrong
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
