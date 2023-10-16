package service

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type FileViewParams struct {
	ID       int
	DirPath  string
	Filename string
}

func (s Service) GetContent(ctx context.Context, params FileViewParams) error {
	if _, err := os.Stat(params.DirPath); errors.Is(err, os.ErrNotExist) {
		makeDirErr := os.MkdirAll(params.DirPath, 0755)
		if makeDirErr != nil {
			err = fmt.Errorf("creating dir: %w", makeDirErr)
		}
	}

	file, err := os.Create(params.Filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = fmt.Errorf("file closing: %w", closeErr)
		}
	}()

	err = s.repo.GetContent(ctx, params.ID, file)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	return nil
}
