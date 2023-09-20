package upload

import (
	"FrankRGTask/api/fileHandler"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2 MB

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), fileHandler.DBTimeout)
	defer cancel()

	parentID, err := fileHandler.Repo.GetParent(ctx, name)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(MAX_UPLOAD_SIZE)
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
		fmt.Println(err)
	}

	fileNew := models.File{
		Name:        handler.Filename,
		Size:        handler.Size,
		ModTime:     time.Now(),
		IsDirectory: false,
		Content:     fileBytes,
		ParentID:    parentID,
	}

	err = fileHandler.Repo.Create(ctx, &fileNew)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	util.WriteJSON(w, http.StatusOK, fileNew)

	logrus.Infof("file '%s' was successfully uploaded\n", handler.Filename)
}
