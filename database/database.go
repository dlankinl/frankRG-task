package database

import (
	"FrankRGTask/config"
	_ "FrankRGTask/internal/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB

func ConnectDB(config config.Config, connStr string) *gorm.DB {
	fn := "database.ConnectDB"

	db, err := gorm.Open(config.DBDriver, connStr)
	if err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	if err = db.DB().Ping(); err != nil {
		logrus.Fatalf("%s: %s\n", fn, err)
	}

	return db
}
