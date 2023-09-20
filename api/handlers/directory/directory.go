package directory

import (
	"FrankRGTask/api/fileHandler"
	"FrankRGTask/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

func DirHandler(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), fileHandler.DBTimeout)
	defer cancel()

	parentID, err := fileHandler.Repo.GetParent(ctx, dirName)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	filesList, err := fileHandler.Repo.GetFilesInDir(ctx, parentID)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	template := template.Must(template.ParseFiles("front/index.html"))
	template.Execute(w, filesList)
}
