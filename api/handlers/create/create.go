package create

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir""`
}

var file models.File

func CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	fn := "api.handlers.create.CreateFileHandler"

	name := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	var id int

	query := `SELECT id FROM Files WHERE name = $1 AND size = 0`

	_ = models.DB.QueryRowContext(ctx, query, name).Scan(&id)

	var fileResp FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileResp)
	if err != nil {
		logrus.Warnf("%s: %s\n", fn, err)
		return
	}

	fileCreated := models.NewFile(fileResp.Name, fileResp.Size, time.Now(), fileResp.IsDir, []byte(fileResp.Content), id)
	util.WriteJSON(w, http.StatusOK, fileCreated)

	logrus.Infof("file '%s' was successfully added\n", fileResp.Name)
}
