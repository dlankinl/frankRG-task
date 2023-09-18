package file

import (
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type FileResponse struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	queryContent := `SELECT Content FROM Files WHERE name = $1`

	var content []byte

	_ = models.DB.QueryRowContext(ctx, queryContent, filename).Scan(&content)

	var resp FileResponse

	resp.Content = string(content)
	resp.Filename = filename

	util.WriteJSON(w, http.StatusOK, resp)
}
