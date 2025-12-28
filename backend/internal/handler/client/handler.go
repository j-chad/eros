package client

import (
	"backend/internal/service"
	"backend/pkg/response"
	"net/http"
)

type Handler struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"message": "Client authentication successful"})
}
