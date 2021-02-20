package main

import (
	"log"
	"net/http"

	businessConnService "github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/delivery/http"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/delivery/taskManager"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/repository"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/usecase"
	"github.com/Toringol/avito-mx-backend-test-task/app/models"
	"github.com/Toringol/avito-mx-backend-test-task/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	logger := logrus.New()

	us := usecase.NewUsecase(repository.NewRepository())
	taskQueue := make(chan models.Task)
	statsQueue := make(chan models.TaskStats)
	stopCh := make(chan struct{})

	taskManager := taskManager.NewTaskManager(us, taskQueue, statsQueue, stopCh, logger)

	go taskManager.TaskManager()

	router := businessConnService.NewHandlers(us, taskQueue, logger)

	log.Fatal(http.ListenAndServe(viper.GetString("portListen"), router))
}
