package routes

import (
	"FrankRGTask/api/handlers/create"
	"FrankRGTask/api/handlers/directory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", directory.DirHandler)
	//router.Get("/api/createfile", create.CreateFileHandler)
	router.Post("/api/createfile", create.CreateFileHandler)
	//router.Get("/api/createdir", create.CreateDirectoryHandler)

	return router
}
