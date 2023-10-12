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
	"strconv"
)

func (s Service) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	content, err := s.fileService.GetContent(r.Context(), fileService.FileViewParams{
		ID: intID,
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

	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = w.Write(content)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
}
