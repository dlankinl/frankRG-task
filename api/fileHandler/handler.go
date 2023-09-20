package fileHandler

import (
	"FrankRGTask/internal/repository/file"
	"time"
)

var Repo *file.PostgresDB

const DBTimeout = time.Second * 3

func SetRepository(repo *file.PostgresDB) {
	Repo = repo
}
