package service

import (
	"FrankRGTask/internal/models"
	"context"
	"time"
)

type UploadParams struct {
	Name        string
	Size        int64
	IsDirectory bool
	Content     []byte
	ParentDir   string
}

func (s Service) Upload(ctx context.Context, params UploadParams) error {
	parentID, err := s.repo.GetParent(ctx, params.ParentDir)
	if err != nil {
		return err
	}

	file := models.File{
		Name:        params.Name,
		Size:        params.Size,
		ModTime:     time.Now().UTC(),
		IsDirectory: params.IsDirectory,
		Content:     params.Content,
		ParentID:    parentID,
	}

	err = s.repo.Create(ctx, &file)
	if err != nil {
		return err
	}

	return nil
}
