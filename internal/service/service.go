package service

import (
	"FrankRGTask/internal/models"
	"context"
	"io"
)

type Repository interface {
	CreateDir(ctx context.Context, file *models.File) error
	Create(ctx context.Context, file *models.File, content []byte) error
	Upload(ctx context.Context, file *models.File, reader io.Reader) error
	GetParent(ctx context.Context, name string) (int, error)
	Rename(ctx context.Context, newName string, id int) error
	FindFilesRecursive(ctx context.Context, id int) ([]int, error)
	GetFilesInDir(ctx context.Context, id int) ([]models.File, error)
	Download(ctx context.Context, id int, readFn func(reader io.Reader) error) error
	DeleteFile(ctx context.Context, id int) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return Service{
		repo: r,
	}
}
