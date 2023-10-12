package main

import (
	"FrankRGTask/config"
	"FrankRGTask/database"
	"FrankRGTask/internal/api/handlers"
	_ "FrankRGTask/internal/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Serve(addr string, router *chi.Mux) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	logrus.Infof("server is listening on %s address", addr)
	return srv.ListenAndServe()
}

func main() {
	fn := "cmd.main.main"

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)
	db := database.ConnectDBAndMigrate(cfg, connStr)

	router := chi.NewRouter()
	handlers.Register(router, db)

	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	err = Serve(address, router)
	if err != nil {
		logrus.Fatal("error while serving http server: ", err)
	}
}
