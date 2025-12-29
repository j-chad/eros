package client

import (
	"backend/internal/service"
)

type Handler struct {
	authService   *service.AuthService
	favourService *service.FavourService
}

func NewHandler(
	authService *service.AuthService,
	favourService *service.FavourService,
) *Handler {
	return &Handler{authService: authService, favourService: favourService}
}
