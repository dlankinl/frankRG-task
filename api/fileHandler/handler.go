package fileHandler

import "FrankRGTask/internal/repository/file"

var Repo *file.PostgresDB

func SetRepository(repo *file.PostgresDB) {
	Repo = repo
}
