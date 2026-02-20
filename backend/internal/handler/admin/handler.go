package admin

import (
	"backend/internal/service"
)

type Handler struct {
	adminService *service.AdminService
}

func NewHandler(adminService *service.AdminService) *Handler {
	return &Handler{adminService: adminService}
}
