package service

import (
	"context"
	"fmt"
)

type DeleteParams struct {
	ID int
}

func (s Service) Delete(ctx context.Context, params DeleteParams) error {
	err := s.repo.DeleteFile(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("deleting file/directory by id: %w", err)
	}

	return nil
}
