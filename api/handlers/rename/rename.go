package rename

import (
	"FrankRGTask/api/fileHandler"
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FileRequest struct {
	ID      int    `json:"id"`
	Newname string `json:"new_name"`
}

func RenameFile(w http.ResponseWriter, r *http.Request) {
	var fileResp FileRequest
	err := json.NewDecoder(r.Body).Decode(&fileResp)
	if err != nil {
		logrus.Warnf("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	//query := `UPDATE files
	//		SET name = $1
	//		WHERE id = $2
	//	`
	//
	//res, err := models.DB.ExecContext(ctx, query, fileResp.Newname, fileResp.ID)
	//if err != nil {
	//	logrus.Warnf("%s\n", err)
	//	util.ErrorJSON(w, err, http.StatusBadRequest)
	//	return
	//}
	//
	//rowsAffected, err := res.RowsAffected()
	//if err != nil {
	//	logrus.Warnf("%s\n", err)
	//	util.ErrorJSON(w, err, http.StatusBadRequest)
	//	return
	//}
	//
	//if rowsAffected == 0 {
	//	logrus.Warnf("%s\n", err)
	//	util.ErrorJSON(w, errors.New("nothing found to update"), http.StatusNotFound)
	//	return
	//}

	err = fileHandler.Repo.Rename(ctx, fileResp.Newname, fileResp.ID)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	util.WriteJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
	logrus.Infof("file with id=%d was successfully renamed on '%s'\n", fileResp.ID, fileResp.Newname)
}
