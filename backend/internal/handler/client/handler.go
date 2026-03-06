package client

import (
	"backend/internal/service"
)

type Handler struct {
	authService   *service.AuthService
	favourService *service.FavourService
	graphService  *service.GraphService
	fileService   *service.FileService
}

func NewHandler(
	authService *service.AuthService,
	favourService *service.FavourService,
	graphService *service.GraphService,
	fileService *service.FileService,
) *Handler {
	return &Handler{authService: authService, favourService: favourService, graphService: graphService, fileService: fileService}
}
