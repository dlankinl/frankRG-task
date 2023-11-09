package handlers

import (
	_ "FrankRGTask/internal/logger"
	fileService "FrankRGTask/internal/service"
	"FrankRGTask/internal/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

func (s Handlers) ListDirFiles(w http.ResponseWriter, r *http.Request) {
	type FileRequest struct {
		Name string
	}

	dirName := chi.URLParam(r, "name")

	files, err := s.fileService.ListDirFiles(r.Context(), fileService.DirParams{
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
