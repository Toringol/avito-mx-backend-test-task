package taskManager

import (
	"database/sql"
	"mime/multipart"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
	"github.com/Toringol/avito-mx-backend-test-task/tools"
	"github.com/sirupsen/logrus"
)

type taskManager struct {
	usecase    businessConnService.IUsecase
	taskQueue  chan models.Task
	statsQueue chan models.TaskStats
	stopCh     chan struct{}
	logger     *logrus.Logger
}

func NewTaskManager(us businessConnService.IUsecase, taskQueue chan models.Task,
	statsQueue chan models.TaskStats, stopCh chan struct{}, logger *logrus.Logger) *taskManager {
	return &taskManager{
		usecase:    us,
		taskQueue:  taskQueue,
		statsQueue: statsQueue,
		stopCh:     stopCh,
		logger:     logger,
	}
}

func (tm *taskManager) TaskManager() {
	for {
		select {
		case taskInfo := <-tm.taskQueue:
			go tm.uploadUserFilesPackProducer(&taskInfo, tm.statsQueue)
		case stats := <-tm.statsQueue:
			go tm.uploadStatsProducer(stats)
		case <-tm.stopCh:
			tm.logger.Info("Stop TaskManager")
			return
		}
	}
}

func (tm *taskManager) uploadUserFilesPackProducer(taskInfo *models.Task, statsQueue chan models.TaskStats) {
	_, err := tm.usecase.UpdateTaskState(taskInfo.TaskID, "IN PROGRESS")
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}

	taskStats := new(models.TaskStats)
	taskStats.TaskID = taskInfo.TaskID

	fileStatsQueue := make(chan models.TaskStats)
	endFileStats := make(chan struct{})
	var wg sync.WaitGroup

	for _, fheaders := range taskInfo.Files {
		for _, hdr := range fheaders {
			wg.Add(1)
			go tm.uploadFileProducer(hdr, taskInfo, fileStatsQueue, &wg)
		}
	}

	go func() {
		for fileStats := range fileStatsQueue {
			taskStats.ProductsCreated += fileStats.ProductsCreated
			taskStats.ProductsUpdated += fileStats.ProductsUpdated
			taskStats.ProductsDeleted += fileStats.ProductsDeleted
			taskStats.RowsWithErrors += fileStats.RowsWithErrors
		}
		endFileStats <- struct{}{}
	}()

	wg.Wait()

	close(fileStatsQueue)

	<-endFileStats

	statsQueue <- *taskStats
}

func (tm *taskManager) uploadFileProducer(hdr *multipart.FileHeader, taskInfo *models.Task,
	fileStatsQueue chan models.TaskStats, wg *sync.WaitGroup) {

	defer wg.Done()

	var sheetWG sync.WaitGroup

	fd, err := hdr.Open()
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}

	f, err := excelize.OpenReader(fd)
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}

	sheets := f.GetSheetMap()
	for _, sheet := range sheets {
		sheetWG.Add(1)

		go tm.uploadFileSheetProducer(f, taskInfo, sheet, fileStatsQueue, &sheetWG)
	}

	sheetWG.Wait()
}

func (tm *taskManager) uploadFileSheetProducer(f *excelize.File, taskInfo *models.Task, sheet string,
	fileStatsQueue chan models.TaskStats, sheetWG *sync.WaitGroup) {

	defer sheetWG.Done()

	fileStats := new(models.TaskStats)

	rows, err := f.GetRows(sheet)
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}
	for _, row := range rows {
		productInfo, err := tools.ConvertXlsxRowToProductInfo(row, taskInfo.SellerID)
		if err != nil {
			fileStats.RowsWithErrors++

			tm.logger.WithField("ErrInfo", err.Error()).Info("InternalError")
			break
		}

		productRecord, err := tm.usecase.SelectProduct(productInfo.SellerID, productInfo.OfferID)
		switch {
		case err == sql.ErrNoRows && productInfo.Available:
			rowsAffected, err := tm.usecase.CreateProduct(productInfo)
			if err != nil {
				tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
				return
			}

			fileStats.ProductsCreated += rowsAffected
			continue
		case err != nil:
			return
		}

		if !productInfo.Available {
			rowsAffected, err := tm.usecase.DeleteProduct(productInfo.SellerID, productInfo.OfferID)
			if err != nil {
				tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
				return
			}

			fileStats.ProductsDeleted += rowsAffected
		} else {
			productRecord.Name = productInfo.Name
			productRecord.Price = productInfo.Price
			productRecord.Quantity = productInfo.Quantity

			rowsAffected, err := tm.usecase.UpdateProduct(productRecord)
			if err != nil {
				tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
				return
			}

			fileStats.ProductsUpdated += rowsAffected
		}
	}

	fileStatsQueue <- *fileStats
}

func (tm *taskManager) uploadStatsProducer(stats models.TaskStats) {
	_, err := tm.usecase.UpdateTaskState(stats.TaskID, "DONE")
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}

	_, err = tm.usecase.CreateTaskStats(&stats)
	if err != nil {
		tm.logger.WithField("ErrInfo", err.Error()).Error("InternalError")
		return
	}
}
