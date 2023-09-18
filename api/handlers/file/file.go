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
	Content  []byte `json:"content"`
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	//fn := "api.handlers.file.FileHandler"

	//dirName := chi.URLParam(r, "dir")
	filename := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()
	//
	//queryID := `SELECT id FROM Files WHERE name = $1 AND size = 0`
	//
	//var id int
	//
	//_ = models.DB.QueryRowContext(ctx, queryID, dirName).Scan(&id)

	//queryContent := `SELECT Content FROM Files WHERE parentid = $1 AND name = $2`
	queryContent := `SELECT Content FROM Files WHERE name = $1`

	var content []byte

	_ = models.DB.QueryRowContext(ctx, queryContent, filename).Scan(&content)

	var resp FileResponse

	resp.Content = content
	resp.Filename = filename

	util.WriteJSON(w, http.StatusOK, resp)

	//if err != nil {
	//	logrus.Warnf("%s: %s\n", fn, err)
	//	return
	//}

}
