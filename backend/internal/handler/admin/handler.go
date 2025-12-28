package admin

import (
	"backend/internal/service"
	"backend/pkg/response"
	"net/http"
)

type Handler struct {
	adminService *service.AdminService
}

func NewHandler(adminService *service.AdminService) *Handler {
	return &Handler{adminService: adminService}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	response.NoContent(w)
}
