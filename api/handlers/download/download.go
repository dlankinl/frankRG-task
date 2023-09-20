package download

import (
	"FrankRGTask/api/fileHandler"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	errs "FrankRGTask/pkg/errors"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type FileContent struct {
	IsDirectory bool   `json:"is_directory"`
	Content     []byte `json:"content"`
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	//query := `SELECT isdirectory, content FROM files WHERE id = $1`
	//
	//var fileContent FileContent
	//
	//err = models.DB.QueryRowContext(ctx, query, intID).Scan(&fileContent.IsDirectory, &fileContent.Content)
	//if errors.Is(err, sql.ErrNoRows) {
	//	logrus.Warn("no rows found")
	//	util.ErrorJSON(w, errors.New("no files were found"), http.StatusNotFound)
	//	return
	//}
	//
	//if fileContent.IsDirectory {
	//	logrus.Infof("try to download directory id=%d\n", intID)
	//	util.ErrorJSON(w, errors.New("directories aren't allowed to be downloaded"), http.StatusBadRequest)
	//	return
	//}

	content, err := fileHandler.Repo.GetContent(ctx, intID)
	if errors.Is(err, errs.TypeNotFileErr) {
		logrus.Infof("try to download directory id=%d\n", intID)
		util.ErrorJSON(w, errors.New("directories aren't allowed to be downloaded"), http.StatusBadRequest)
		return
	}
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
}
