package service

import "context"

type FileViewParams struct {
	ID int
}

func (s Service) GetContent(ctx context.Context, params FileViewParams) ([]byte, error) {
	content, err := s.repo.GetContent(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return content, nil
}
