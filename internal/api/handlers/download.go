package handlers

import (
	errs "FrankRGTask/internal/errors"
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
)

func (s Service) Download(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
}
