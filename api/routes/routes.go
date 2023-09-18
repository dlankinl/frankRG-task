package routes

import (
	"FrankRGTask/api/handlers/create"
	"FrankRGTask/api/handlers/directory"
	"FrankRGTask/api/handlers/upload"
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dir/root", http.StatusPermanentRedirect)
	})
	//router.Get("/root", directory.DirHandler)
	router.Post("/api/createfile/{name}", create.CreateFileHandler)
	router.Post("/api/uploadfile/{name}", upload.UploadFileHandler)
	router.Get("/dir/{name}", directory.DirHandler)
	//router.Get("/file/{dir}/{name}")

	return router
}
