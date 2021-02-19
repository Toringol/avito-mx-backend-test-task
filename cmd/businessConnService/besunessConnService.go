package main

import (
	"log"
	"net/http"

	businessConnService "github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/delivery/http"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/repository"
	"github.com/Toringol/avito-mx-backend-test-task/app/businessConnService/usecase"
	"github.com/Toringol/avito-mx-backend-test-task/config"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	router := businessConnService.NewHandlers(usecase.NewUsecase(repository.NewRepository()))

	log.Fatal(http.ListenAndServe(viper.GetString("portListen"), router))
}
