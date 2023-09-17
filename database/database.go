package database

import (
	"FrankRGTask/config"
	_ "FrankRGTask/internal/logger"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func ConnectDB(config config.Config, connStr string) *sql.DB {
	fn := "database.ConnectDB"

	db, err := sql.Open(config.DBDriver, connStr)
	if err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	if err = db.Ping(); err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	return db
}
