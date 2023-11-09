package service

import (
	"FrankRGTask/internal/models"
	"context"
	"fmt"
)

type DirParams struct {
	Name string
}

func (s Service) ListDirFiles(ctx context.Context, params DirParams) ([]models.File, error) {
	id, err := s.repo.GetParent(ctx, params.Name)
	if err != nil {
		return nil, fmt.Errorf("getting id of parent directory: %w", err)
	}

	files, err := s.repo.GetFilesInDir(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting files in directory: %w", err)
	}

	return files, nil
}
