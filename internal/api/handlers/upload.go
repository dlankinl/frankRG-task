package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2 MB

func (s Service) Upload(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		logrus.Warnf("error retrieving the file: %s\n", err)
		util.ErrorJSON(w, errors.New("bad json data"), http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Warnf("error reading the file: %s\n", err)
		util.ErrorJSON(w, errors.New("bad file content"), http.StatusBadRequest)
		return
	}

	newFile, err := s.fileService.Upload(r.Context(), fileService.UploadParams{
		Name:        handler.Filename,
		Size:        handler.Size,
		IsDirectory: false,
		Content:     fileBytes,
		ParentDir:   name,
	})

	err = util.WriteJSON(w, http.StatusOK, newFile)
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("file '%s' was successfully uploaded\n", handler.Filename)
}
