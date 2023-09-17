package main

import (
	config2 "FrankRGTask/config"
	"FrankRGTask/database"
	"FrankRGTask/internal/handlers"
	_ "FrankRGTask/internal/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	fn := "cmd.main.main"

	config, err := config2.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	http.Handle("/", handlers.RootHandler())

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode)
	db := database.ConnectDB(config, connStr)

	http.ListenAndServe(":7070", nil)
}
