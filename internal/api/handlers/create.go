package handlers

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var file models.File

func (s service) Create() http.HandlerFunc {
	type FileRequest struct {
		Name      string `json:"name"`
		Content   string `json:"content"`
		Size      int64  `json:"size"`
		IsDir     bool   `json:"is_dir"`
		ParentDir string `json:"parent_dir"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var fileReq FileRequest
		err := json.NewDecoder(r.Body).Decode(&fileReq)
		if err != nil {
			logrus.Warnf("%s\n", err)
			util.ErrorJSON(w, errors.New("bad json request"), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
		defer cancel()

		err = s.fileService.Create(ctx, fileService.CreateParams{
			Name:        fileReq.Name,
			Size:        fileReq.Size,
			IsDirectory: fileReq.IsDir,
			Content:     []byte(fileReq.Content),
		})

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

		logrus.Infof("file '%s' was successfully added\n", fileReq.Name)
	}
}
