package models

import (
	"database/sql"
	"time"
)

var DB *sql.DB
var DBTimeout = time.Second * 3

func SetDatabase(db *sql.DB) {
	DB = db
}
