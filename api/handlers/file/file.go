package file

import (
	"FrankRGTask/api/fileHandler"
	"FrankRGTask/internal/util"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type FileResponse struct {
	Filename string `json:"filename"`
	ID       int    `json:"id"`
	Content  string `json:"content"`
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Warnf("%s\n", err)
		return
	}

	filename := chi.URLParam(r, "name")

	logrus.Info("HEREEEEEEE", intID, filename)

	ctx, cancel := context.WithTimeout(context.Background(), fileHandler.DBTimeout)
	defer cancel()

	content, err := fileHandler.Repo.GetContent(ctx, intID)
	if err != nil {
		logrus.Infof("%s\n", err)
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var resp FileResponse

	resp.ID = intID
	resp.Content = string(content)
	resp.Filename = filename

	fmt.Println(resp)

	util.WriteJSON(w, http.StatusOK, resp)
}
