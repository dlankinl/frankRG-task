package rename

import (
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
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	query := `UPDATE files
			SET name = $1
			WHERE id = $2
		`

	_, err = models.DB.ExecContext(ctx, query, fileResp.Newname, fileResp.ID)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	util.WriteJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
	logrus.Infof("file with id=%d was successfully renamed on '%s'\n", fileResp.ID, fileResp.Newname)
}