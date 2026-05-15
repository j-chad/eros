package client

import (
	"backend/internal/service"
)

type Handler struct {
	authService   *service.AuthService
	favourService *service.FavourService
	graphService  *service.GraphService
	fileService   *service.FileService
	pushService   *service.PushService
}

func NewHandler(
	authService *service.AuthService,
	favourService *service.FavourService,
	graphService *service.GraphService,
	fileService *service.FileService,
	pushService *service.PushService,
) *Handler {
	return &Handler{authService: authService, favourService: favourService, graphService: graphService, fileService: fileService, pushService: pushService}
}
