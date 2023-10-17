package handlers

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	type FileRequest struct {
		Name      string
		Size      int64
		ParentDir string
	}

	var fileReq FileRequest

	fileReq.Name = r.Header.Get("name")
	fileReq.ParentDir = r.Header.Get("parent_dir")
	fileReq.Size = r.ContentLength

	newFile, err := s.fileService.Upload(r.Context(), r.Body, service.UploadParams{
		Name:      fileReq.Name,
		Size:      fileReq.Size,
		Path:      sql.NullString{},
		ParentDir: fileReq.ParentDir,
	})

	err = util.WriteJSON(w, http.StatusOK, newFile)
	if err != nil {
		logrus.Infof("error while writing json response: %w", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("file '%s' was successfully uploaded\n", newFile.Name)
}
