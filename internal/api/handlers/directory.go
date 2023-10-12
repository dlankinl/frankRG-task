package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

func (s service) ListDirFiles() http.HandlerFunc {
	type FileRequest struct {
		Name string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dirName := chi.URLParam(r, "name")

		ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
		defer cancel()

		files, err := s.fileService.ListDirFiles(ctx, fileService.DirParams{
			Name: dirName,
		})

		if err != nil {
			logrus.Warnf("%s\n", err)
			util.ErrorJSON(w, errors.New("couldn't find files list in directory"), http.StatusBadRequest)
			return
		}

		htmlTempl := template.Must(template.ParseFiles("front/index.html"))
		err = htmlTempl.Execute(w, files)
		if err != nil {
			logrus.Warnf("%s\n", err)
			util.ErrorJSON(w, errors.New("couldn't execute html template"), http.StatusBadRequest)
		}
	}
}
