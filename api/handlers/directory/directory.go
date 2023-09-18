package directory

import (
	"FrankRGTask/internal/models"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"text/template"
)

//func DirHandler(w http.ResponseWriter, r *http.Request) {
//	fn := "api.handlers.directory.DirHandler"
//
//	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
//	defer cancel()
//
//	query := `SELECT * FROM Files WHERE parentid = $1`
//
//	rows, err := models.DB.QueryContext(ctx, query, 1)
//	if err != nil {
//		logrus.Warnf("%s: %s\n", fn, err)
//		return
//	}
//
//	var filesList []models.File
//	for rows.Next() {
//		var file models.File
//
//		s := reflect.ValueOf(&file).Elem()
//		numCols := s.NumField()
//		columns := make([]interface{}, numCols)
//		for i := 0; i < numCols; i++ {
//			field := s.Field(i)
//			columns[i] = field.Addr().Interface()
//		}
//
//		err = rows.Scan(columns...)
//		if err != nil {
//			logrus.Warnf("%s: %s\n", fn, err)
//			return
//		}
//		filesList = append(filesList, file)
//	}
//
//	template := template.Must(template.ParseFiles("front/index.html"))
//	template.Execute(w, filesList)
//}

func DirHandler(w http.ResponseWriter, r *http.Request) {
	fn := "api.handlers.directory.DirHandler"

	dirName := chi.URLParam(r, "name")

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	query1 := `SELECT id FROM Files WHERE name = $1 AND size = 0`

	var id int

	_ = models.DB.QueryRowContext(ctx, query1, dirName).Scan(&id)

	query := `SELECT * FROM Files WHERE parentid = $1`

	//rows, err := models.DB.QueryContext(ctx, query, 1)
	rows, err := models.DB.QueryContext(ctx, query, id)
	if err != nil {
		logrus.Warnf("%s: %s\n", fn, err)
		return
	}

	var filesList []models.File
	for rows.Next() {
		var file models.File

		s := reflect.ValueOf(&file).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err = rows.Scan(columns...)
		if err != nil {
			logrus.Warnf("%s: %s\n", fn, err)
			return
		}
		filesList = append(filesList, file)
	}

	template := template.Must(template.ParseFiles("front/index.html"))
	template.Execute(w, filesList)
}
