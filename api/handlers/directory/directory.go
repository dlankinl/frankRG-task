package directory

import (
	"FrankRGTask/api/fileHandler"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

func DirHandler(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	//query1 := `SELECT id FROM Files WHERE name = $1 AND size = 0`
	//
	//var id int
	//
	//err := models.DB.QueryRowContext(ctx, query1, dirName).Scan(&id)
	//
	//if errors.Is(err, sql.ErrNoRows) {
	//	logrus.Infof("%s\n", err)
	//	util.ErrorJSON(w, errors.New("dir wasn't found"), http.StatusNotFound)
	//	return
	//}
	//
	//query := `SELECT * FROM Files WHERE parentid = $1`
	//
	//rows, err := models.DB.QueryContext(ctx, query, id)
	//
	//if errors.Is(err, sql.ErrNoRows) {
	//	logrus.Infof("%s\n", err)
	//	util.ErrorJSON(w, errors.New("no parentDir id found"), http.StatusNotFound)
	//	return
	//}
	//if err != nil {
	//	logrus.Warnf("%s\n", err)
	//	util.ErrorJSON(w, err, http.StatusBadRequest)
	//	return
	//}
	//
	//var filesList []models.File
	//for rows.Next() {
	//	var file models.File
	//
	//	s := reflect.ValueOf(&file).Elem()
	//	numCols := s.NumField()
	//	columns := make([]interface{}, numCols)
	//	for i := 0; i < numCols; i++ {
	//		field := s.Field(i)
	//		columns[i] = field.Addr().Interface()
	//	}
	//
	//	err = rows.Scan(columns...)
	//	if err != nil {
	//		logrus.Warnf("%s\n", err)
	//		util.ErrorJSON(w, err, http.StatusBadRequest)
	//		return
	//	}
	//	filesList = append(filesList, file)
	//}

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
