package client

import (
	"backend/internal/service"
)

type Handler struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{authService: authService}
}
