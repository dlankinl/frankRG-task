package create

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FileRequest struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Size      int64  `json:"size"`
	IsDir     bool   `json:"is_dir"`
	ParentDir string `json:"parent_dir"`
}

var file models.File

func CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	var fileResp FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileResp)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, errors.New("bad json request"), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	var id int

	query := `SELECT id FROM Files WHERE name = $1 AND size = 0`

	err = models.DB.QueryRowContext(ctx, query, fileResp.ParentDir).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, errors.New("no parentDir found"), http.StatusBadRequest)
		return
	}

	fileCreated := models.NewFile(fileResp.Name, fileResp.Size, time.Now(), fileResp.IsDir, []byte(fileResp.Content), id)
	util.WriteJSON(w, http.StatusOK, fileCreated)

	logrus.Infof("file '%s' was successfully added\n", fileResp.Name)
}
