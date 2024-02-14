package main

import (
	"ewallet/pkg/config"
	"ewallet/pkg/handler"
	"ewallet/pkg/server"
	"ewallet/pkg/service"
	"ewallet/pkg/storage/postgres"
	
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoadCfg()
	gin.SetMode(cfg.GinMode)

	pgStorage, err := postgres.New(cfg.PgCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := pgStorage.InitDB(); err != nil {
		log.Fatal(err.Error())
	}

	services := service.NewService(pgStorage)
	handlers := handler.New(services)
	server := new(server.Server)

	go func() {
		err := server.Run(cfg.Port, handlers.InitRoutes())
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("EWallet started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("EWallet shutting down")
}
