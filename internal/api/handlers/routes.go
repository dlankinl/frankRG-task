package handlers

import fileService "FrankRGTask/internal/service"

type FileService interface {
}

type Service struct {
	fileService fileService.Service
}

func NewHandler(service fileService.Service) Service {
	return Service{
		fileService: service,
	}
}
