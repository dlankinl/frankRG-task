package service

import (
	_ "FrankRGTask/internal/logger"
	"FrankRGTask/internal/models"
	"context"
	"fmt"
	"time"
)

type CreateParams struct {
	Name        string
	Size        int64
	IsDirectory bool
	Content     []byte
	ParentDir   string
}

func (s Service) Create(ctx context.Context, params CreateParams) error {
	parentID, err := s.repo.GetParent(ctx, params.ParentDir)
	if err != nil {
		return fmt.Errorf("getting id of parent directory: %w", err)
	}

	file := models.File{
		Name:        params.Name,
		Size:        params.Size,
		ModTime:     time.Now().UTC(),
		IsDirectory: params.IsDirectory,
		ParentID:    parentID,
	}

	if params.IsDirectory {
		err = s.repo.CreateDir(ctx, &file)
	} else {
		err = s.repo.Create(ctx, &file, params.Content)
	}
	if err != nil {
		return fmt.Errorf("creating file/directory: %w", err)
	}

	return nil
}
