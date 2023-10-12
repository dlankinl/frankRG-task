package handlers

import (
	fileRepo "FrankRGTask/internal/repository/file"
	fileService "FrankRGTask/internal/service"
	"database/sql"
)

type service struct {
	fileService fileService.Service
}

func NewHandler(db *sql.DB) Handler {
	return service{
		fileService: fileService.NewService(fileRepo.NewDBConnection(db)),
	}
}
