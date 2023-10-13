package database

import (
	_ "FrankRGTask/internal/logger"
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func ConnectDB(connStr string) *pgx.Conn {
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logrus.Fatalf("%s\n", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		logrus.Fatalf("%s\n", err)
	}

	logrus.Info("postgres started successfully")

	return db
}
