package models

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"time"
)

type File struct {
	Name        string
	Size        int64
	ModTime     time.Time
	IsDirectory bool
	Content     []byte
	SubFiles    []File
	ParentID
}

func NewFile(name string, size int64, modTime time.Time, isDirectory bool, content []byte, subFiles []File) File {
	return File{
		Name:        name,
		Size:        size,
		ModTime:     modTime,
		IsDirectory: isDirectory,
		Content:     content,
		SubFiles:    subFiles,
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

func (f *File) AddSubFiles(subFiles []File) {
	fn := "internal.models.file.AddSubFiles"

	if len(subFiles) == 0 {
		logrus.Infof("%s: %s\n", fn, "subFiles length is zero")
		return
	}

	f.SubFiles = append(f.SubFiles, subFiles...)
}
