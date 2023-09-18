package create

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir""`
}

var file models.File

func CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	fn := "api.handlers.create.CreateFileHandler"

	var fileResp FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileResp)
	if err != nil {
		logrus.Warnf("%s: %s\n", fn, err)
		return
	}
	fmt.Println(fileResp)

	fileCreated := models.NewFile(fileResp.Name, fileResp.Size, time.Now(), fileResp.IsDir, []byte(fileResp.Content), 1)
	util.WriteJSON(w, http.StatusOK, fileCreated)
}

//// TODO: should be POST!
//func CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
//	var fileResp FileRequest
//
//	err := json.NewDecoder(r.Body).Decode(&fileResp)
//	if err != nil {
//		logrus.Warnf("%s: %s\n", fn, err)
//		return
//	}
//
//	fileCreated := models.NewFile(fileResp.Name, 0, time.Now(), true, []byte{}, 1)
//	util.WriteJSON(w, http.StatusOK, fileCreated)
//}
