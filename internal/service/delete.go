package service

import (
	"context"
)

type DeleteParams struct {
	ID int
}

func (s Service) Delete(ctx context.Context, params DeleteParams) (int, error) {
	deletedRows, err := s.repo.DeleteByID(ctx, params.ID)
	if err != nil {
		return 0, err
	}

	return deletedRows, nil
}
