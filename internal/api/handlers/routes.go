package handlers

import (
	fileRepo "FrankRGTask/internal/repository/file"
	fileService "FrankRGTask/internal/service"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

const DBTimeout = time.Second * 3

type service struct {
	fileService fileService.Service
	router      http.Handler
}

func newHandler(db *sql.DB) service {
	return service{
		fileService: fileService.NewService(fileRepo.NewDBConnection(db)),
	}
}

func Register(router *chi.Mux, db *sql.DB) {
	handler := newHandler(db)

	router.Post("/api/createfile/", handler.Create())
	router.Post("/api/uploadfile/{name}", handler.Upload())
	router.Post("/api/file", handler.Rename())
	router.Get("/file/{id}/{name}", handler.GetContent())
	router.Get("/dir/{name}", handler.ListDirFiles())
	router.Get("/api/downloadfile/{id}", handler.Download())
	router.Delete("/api/file/{id}", handler.Delete())
}
