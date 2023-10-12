package handlers

import (
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s Service) Delete(w http.ResponseWriter, r *http.Request) {
	type FileRequest struct {
		ID int
	}

	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, errors.New("couldn't convert id from 'string' type to 'int'"), http.StatusBadRequest)
		return
	}

	deletedRows, err := s.fileService.Delete(r.Context(), fileService.DeleteParams{
		ID: intID,
	})

	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = util.WriteJSON(
		w,
		http.StatusOK,
		struct {
			Status      string `json:"status"`
			DeletedRows int    `json:"deleted_rows"`
		}{
			Status:      "OK",
			DeletedRows: deletedRows,
		},
	)
	if err != nil {
		logrus.Infof("error while writing json response: %s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("successfully deleted %d rows\n", deletedRows)
}
