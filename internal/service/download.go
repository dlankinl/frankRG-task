package service

import (
	"context"
	"io"
)

type DownloadParams struct {
	FileId int
}

func (s Service) Download(ctx context.Context, writer io.Writer, params DownloadParams) error {
	err := s.repo.Download(ctx, params.FileId, func(reader io.Reader) error {
		_, err := io.Copy(writer, reader)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
