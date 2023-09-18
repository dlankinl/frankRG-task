package models

import (
	_ "FrankRGTask/internal/logger"
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type File struct {
	ID          int
	Name        string
	Size        int64
	ModTime     time.Time
	IsDirectory bool
	Content     []byte
	ParentID    int
}

func NewFile(name string, size int64, modTime time.Time, isDirectory bool, content []byte, parentID int) *File {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	query := `INSERT INTO Files(name, size, modtime, isdirectory, content, parentid) 
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	_ = DB.QueryRowContext(ctx, query, name, size, modTime, isDirectory, content, parentID).Scan(&id)

	return &File{
		ID:          id,
		Name:        name,
		Size:        size,
		ModTime:     modTime,
		IsDirectory: isDirectory,
		Content:     content,
		ParentID:    parentID,
	}
}

func (f *File) ChangeContent(content []byte) {
	fn := "internal.models.file.ChangeContent"

	if bytes.Equal(f.Content, content) {
		logrus.Infof("%s: %s\n", fn, "trying to insert same content")
		return
	}
	f.Content = content
}

//func ChangeFileContent()
