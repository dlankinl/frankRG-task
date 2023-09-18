package delete

import (
	"FrankRGTask/internal/models"
	"FrankRGTask/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	logrus.Info("ID: ", intID)

	ctx, cancel := context.WithTimeout(context.Background(), models.DBTimeout)
	defer cancel()

	query := `
		WITH RECURSIVE DirectoryHierarchy AS (
		    SELECT id FROM files WHERE id = $1           
		    UNION ALL 
		    SELECT f.id FROM files f
			INNER JOIN DirectoryHierarchy dh ON f.parentid = dh.id
		)
		SELECT id FROM DirectoryHierarchy;
		`

	rows, err := models.DB.QueryContext(ctx, query, intID)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	var idsToDelete []int
	for rows.Next() {
		var idInner int
		if err = rows.Scan(&idInner); err != nil {
			logrus.Warnf("%s\n", err)
			return
		}
		idsToDelete = append(idsToDelete, idInner)
	}

	deleteQuery := `DELETE FROM files WHERE id = ANY($1::integer[])`

	pgIntArray := pq.Array(idsToDelete)
	_, err = models.DB.ExecContext(ctx, deleteQuery, pgIntArray)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	util.WriteJSON(
		w,
		http.StatusOK,
		struct {
			Status      string `json:"status"`
			DeletedRows int    `json:"deleted_rows"`
		}{
			Status:      "OK",
			DeletedRows: len(idsToDelete),
		},
	)
	logrus.Infof("successfully deleted %d\n", len(idsToDelete))
}
