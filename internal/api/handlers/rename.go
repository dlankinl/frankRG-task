package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s Service) Rename(w http.ResponseWriter, r *http.Request) {
	type FileRequest struct {
		ID      int    `json:"id"`
		Newname string `json:"new_name"`
	}

	var fileReq FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = s.fileService.Rename(r.Context(), fileService.RenameParams{
		Newname: fileReq.Newname,
		ID:      fileReq.ID,
	})

	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("file with id=%d was successfully renamed on '%s'\n", fileReq.ID, fileReq.Newname)
}
