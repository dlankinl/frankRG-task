package models

import (
	_ "FrankRGTask/internal/logger"
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
