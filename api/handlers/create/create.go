package create

import (
	"FrankRGTask/api/fileHandler"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), fileHandler.DBTimeout)
	defer cancel()

	parentID, err := fileHandler.Repo.GetParent(ctx, fileResp.ParentDir)

	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	file = models.File{
		Name:        fileResp.Name,
		Size:        fileResp.Size,
		ModTime:     time.Now(),
		IsDirectory: fileResp.IsDir,
		Content:     []byte(fileResp.Content),
		ParentID:    parentID,
	}

	err = fileHandler.Repo.Create(ctx, &file)

	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	util.WriteJSON(w, http.StatusOK, struct {
		Status string
	}{
		Status: "OK",
	})

	logrus.Infof("file '%s' was successfully added\n", fileResp.Name)
}
