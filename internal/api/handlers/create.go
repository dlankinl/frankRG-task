package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s Service) Create(w http.ResponseWriter, r *http.Request) {
	type FileRequest struct {
		Name      string `json:"name"`
		Content   string `json:"content"`
		Size      int64  `json:"size"`
		IsDir     bool   `json:"is_dir"`
		ParentDir string `json:"parent_dir"`
	}

	var fileReq FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, errors.New("bad json request"), http.StatusBadRequest)
		return
	}

	err = s.fileService.Create(r.Context(), fileService.CreateParams{
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

	err = util.WriteJSON(w, http.StatusOK, struct {
		Status string
	}{
		Status: "OK",
	})
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("file '%s' was successfully added\n", fileReq.Name)
}
