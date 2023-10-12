package service

import (
	"context"
	"fmt"
)

type FileViewParams struct {
	ID int
}

func (s Service) GetContent(ctx context.Context, params FileViewParams) ([]byte, error) {
	content, err := s.repo.GetContent(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("getting content of file: %w", err)
	}

	return content, nil
}
