package main

import (
	"FrankRGTask/api/routes"
	config2 "FrankRGTask/config"
	"FrankRGTask/database"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: routes.Routes(),
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

	//http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort), nil)

	//http.Handle("/", directory.DirHandler(1))
	//http.Get("/api")
	err = Serve(address)
	if err != nil {
		logrus.Fatal("error while serving http server: ", err)
	}
}
