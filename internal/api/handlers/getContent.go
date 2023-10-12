package handlers

import (
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s service) GetContent() http.HandlerFunc {
	type FileResponse struct {
		Filename string `json:"filename"`
		ID       int    `json:"id"`
		Content  string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		intID, err := strconv.Atoi(id)
		if err != nil {
			logrus.Warnf("%s\n", err)
			return
		}

		filename := chi.URLParam(r, "name")

		ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
		defer cancel()

		content, err := s.fileService.GetContent(ctx, fileService.FileViewParams{
			ID: intID,
		})

		if err != nil {
			logrus.Infof("error while getting content: %s\n", err)
			util.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		var resp FileResponse

		resp.ID = intID
		resp.Content = string(content)
		resp.Filename = filename

		err = util.WriteJSON(w, http.StatusOK, resp)
		if err != nil {
			logrus.Infof("error while writing json response: %s\n", err)
			util.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		logrus.Infof("successfully got content of %s file\n", filename)
	}
}
