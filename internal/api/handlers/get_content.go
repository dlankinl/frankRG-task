package handlers

import (
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
)

func (s Service) GetContent(w http.ResponseWriter, r *http.Request) {
	type FileResponse struct {
		Filename string `json:"filename"`
		ID       int    `json:"id"`
		Content  string `json:"content"`
	}

	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	dataDir := "data"
	filename := chi.URLParam(r, "name")
	dPath := filepath.Join(dataDir, id)

	err = s.fileService.GetContent(r.Context(), fileService.FileViewParams{
		ID:       intID,
		DirPath:  dPath,
		Filename: filename,
	})

	if err != nil {
		logrus.Infof("error while getting content: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var resp FileResponse

	resp.ID = intID
	resp.Filename = filename

	err = util.WriteJSON(w, http.StatusOK, resp)
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("successfully got content of %s file\n", filename)
}
