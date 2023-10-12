package handlers

import (
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s service) Delete() http.HandlerFunc {
	type FileRequest struct {
		ID int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		intID, err := strconv.Atoi(id)
		if err != nil {
			logrus.Warnf("%s\n", err)
			util.ErrorJSON(w, errors.New("couldn't convert id from 'string' type to 'int'"), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
		defer cancel()

		deletedRows, err := s.fileService.Delete(ctx, fileService.DeleteParams{
			ID: intID,
		})

		if err != nil {
			logrus.Infof("%s\n", err)
			util.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		util.WriteJSON(
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
		logrus.Infof("successfully deleted %d rows\n", deletedRows)
	}
}
