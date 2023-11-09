package handlers

import (
	errs "FrankRGTask/internal/errors"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s Handlers) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	err = s.fileService.Download(r.Context(), w, service.DownloadParams{
		FileId: intID,
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
