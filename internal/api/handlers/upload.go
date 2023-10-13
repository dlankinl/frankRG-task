package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s Service) Upload(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	type FileRequest struct {
		Name        string `json:"name"`
		Size        int64  `json:"size"`
		IsDirectory bool   `json:"is_directory"`
		Path        string `json:"path"`
		ParentDir   string
	}

	var fileReq FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, errors.New("bad json request"), http.StatusBadRequest)
		return
	}

	fileReq.ParentDir = name
	newFile, err := s.fileService.Upload(r.Context(), fileService.UploadParams{
		Name:        fileReq.Name,
		Size:        fileReq.Size,
		IsDirectory: fileReq.IsDirectory,
		Path:        sql.NullString{String: fileReq.Path, Valid: true},
		ParentDir:   fileReq.ParentDir,
	})

	err = util.WriteJSON(w, http.StatusOK, newFile)
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("file '%s' was successfully uploaded\n", newFile.Name)
}
