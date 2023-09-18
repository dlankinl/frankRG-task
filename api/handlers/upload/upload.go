package upload

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var file models.File

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2 MB

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		logrus.Warnf("error retrieving the file: %s\n", err)
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	fileCreated := models.NewFile(handler.Filename, handler.Size, time.Now(), false, fileBytes, 1)
	util.WriteJSON(w, http.StatusOK, fileCreated)

	logrus.Infof("file '%s' was successfully uploaded\n", handler.Filename)
}

//
//func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
//	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
//		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
//		return
//	}
//
//	file, fileHeader, err := r.FormFile("file")
//	if err != nil {
//		util.ErrorJSON(w, err, http.StatusBadRequest)
//		return
//	}
//
//	defer file.Close()
//
//	fmt.Println(fileHeader, file)
//}
