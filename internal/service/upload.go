package service

import (
	"FrankRGTask/internal/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type UploadParams struct {
	Name        string
	Size        int64
	IsDirectory bool
	Path        sql.NullString
	ParentDir   string
}

func (s Service) Upload(ctx context.Context, params UploadParams) (*models.File, error) {
	fileReader, err := os.OpenFile(params.Path.String, os.O_RDONLY, 0400)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	parentID, err := s.repo.GetParent(ctx, params.ParentDir)
	if err != nil {
		return nil, fmt.Errorf("getting id of parent directory: %w", err)
	}

	file := models.File{
		Name:        params.Name,
		Size:        params.Size,
		ModTime:     time.Now().UTC(),
		IsDirectory: params.IsDirectory,
		Path:        params.Path,
		ParentID:    parentID,
	}

	err = s.repo.Create(ctx, &file, fileReader)
	if err != nil {
		return nil, fmt.Errorf("creating file/directory: %w", err)
	}

	return &file, nil
}
