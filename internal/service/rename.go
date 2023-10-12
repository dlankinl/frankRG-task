package service

import (
	"context"
	"fmt"
)

type RenameParams struct {
	ID      int
	Newname string
}

func (s Service) Rename(ctx context.Context, params RenameParams) error {
	err := s.repo.Rename(ctx, params.Newname, params.ID)
	if err != nil {
		return fmt.Errorf("renaming file: %w", err)
	}

	return nil
}
