package main

import (
	"FrankRGTask/config"
	"FrankRGTask/database"
	"FrankRGTask/internal/api/handlers"
	_ "FrankRGTask/internal/logger"
	fileRepo "FrankRGTask/internal/repository/file"
	fileService "FrankRGTask/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Serve(addr string, router *chi.Mux) error {
	srv := &http.Server{
		Addr:        addr,
		Handler:     router,
		ReadTimeout: time.Second * 4,
		IdleTimeout: time.Second * 60,
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
	db := database.ConnectDB(connStr)

	router := chi.NewRouter()

	repo := fileRepo.NewDBConnection(db)
	service := fileService.NewService(repo)

	handler := handlers.NewHandler(service)

	router.Post("/api/createfile", handler.Create)
	router.Post("/api/uploadfile/{name}", handler.Upload)
	router.Post("/api/file", handler.Rename)
	router.Get("/file/{id}/{name}", handler.GetContent)
	router.Get("/dir/{name}", handler.ListDirFiles)
	router.Get("/api/downloadfile/{id}/{name}", handler.Download)
	router.Delete("/api/file/{id}", handler.Delete)

	address := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	err = Serve(address, router)
	if err != nil {
		logrus.Fatal("error while serving http server: ", err)
	}
}
