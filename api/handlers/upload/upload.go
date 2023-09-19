package upload

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"database/sql"
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

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	var id int

	query := `SELECT id FROM Files WHERE name = $1 AND size = 0`

	err := models.DB.QueryRowContext(ctx, query, name).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, errors.New("no files were found"), http.StatusNotFound)
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

	fileCreated := models.NewFile(handler.Filename, handler.Size, time.Now(), false, fileBytes, id)
	util.WriteJSON(w, http.StatusOK, fileCreated)

	logrus.Infof("file '%s' was successfully uploaded\n", handler.Filename)
}
