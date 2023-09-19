package file

import (
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
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

	err := models.DB.QueryRowContext(ctx, queryContent, filename).Scan(&content)

	if errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, errors.New("no files were found"), http.StatusNotFound)
		return
	}

	var resp FileResponse

	resp.Content = string(content)
	resp.Filename = filename

	util.WriteJSON(w, http.StatusOK, resp)
}
