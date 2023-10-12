package database

import (
	"FrankRGTask/config"
	_ "FrankRGTask/internal/logger"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func migrateUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func ConnectDBAndMigrate(config config.Config, connStr string) *sql.DB {
	//db, err := pgx.Connect(context.Background(), connStr)
	db, err := sql.Open(config.DBDriver, connStr)
	if err != nil {
		logrus.Fatalf("%s\n", err)
	}

	if err = db.Ping(); err != nil {
		logrus.Fatalf("%s\n", err)
	}

	if err = migrateUp(db); err != nil {
		logrus.Fatalf("error while migrating up: %s\n", err)
	} else {
		logrus.Info("successful migrating up")
	}

	logrus.Info("postgres started successfully")

	return db
}
