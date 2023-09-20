package main

import (
	"FrankRGTask/api/fileHandler"
	"FrankRGTask/api/routes"
	config2 "FrankRGTask/config"
	"FrankRGTask/database"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/repository/file"
	"FrankRGTask/pkg/transactor"
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

	txMngr := transactor.NewTransactor()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode)
	db := database.ConnectDBAndMigrate(config, connStr)

	filesRepo := file.NewDBConnection(db, txMngr, "postgres")
	fileHandler.SetRepository(filesRepo)

	address := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)

	err = Serve(address)
	if err != nil {
		logrus.Fatal("error while serving http server: ", err)
	}
}
