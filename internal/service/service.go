package service

import "FrankRGTask/internal/repository/file"

type Service struct {
	repo file.PostgresDB
}

func NewService(r file.PostgresDB) Service {
	return Service{
		repo: r,
	}
}
