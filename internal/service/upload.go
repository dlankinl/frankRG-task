package service

import (
	"FrankRGTask/internal/models"
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"
)

type UploadParams struct {
	Name      string
	Size      int64
	Path      sql.NullString
	ParentDir string
}

func (s Service) Upload(ctx context.Context, req io.Reader, params UploadParams) (*models.File, error) {
	parentID, err := s.repo.GetParent(ctx, params.ParentDir)
	if err != nil {
		return nil, fmt.Errorf("getting id of parent directory: %w", err)
	}

	file := models.File{
		Name:        params.Name,
		Size:        params.Size,
		ModTime:     time.Now().UTC(),
		IsDirectory: false,
		Path:        params.Path,
		ParentID:    parentID,
	}

	err = s.repo.Upload(ctx, &file, req)
	if err != nil {
		return nil, fmt.Errorf("creating file/directory: %w", err)
	}

	return &file, nil
}
