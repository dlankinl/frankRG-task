package main

import (
	config2 "FrankRGTask/config"
	"FrankRGTask/database"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	logrus.Infof("server is listening on %s address", addr)
	return srv.ListenAndServe()
}

func main() {
	fn := "cmd.main.main"

	config, err := config2.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode)
	db := database.ConnectDB(config, connStr)
	models.DB = db

	address := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2023-09-17 00:17:17")
	newFile := models.NewFile("bro", 324, parsedTime, false, []byte("sge"), 1)
	logrus.Infof("New file is: %v\n", newFile)

	//http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort), nil)
	err = Serve(address)
	if err != nil {
		logrus.Fatal("error while serving http server: ", err)
	}
}
