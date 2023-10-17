package handlers

import fileService "FrankRGTask/internal/service"

type FileService interface {
}

type Handlers struct {
	fileService fileService.Service
}

func NewHandler(service fileService.Service) Handlers {
	return Handlers{
		fileService: service,
	}
}
