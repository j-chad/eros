package admin

import (
	"backend/internal/service"
)

type Handler struct {
	adminService *service.AdminService
	pushService  *service.PushService
}

func NewHandler(adminService *service.AdminService, pushService *service.PushService) *Handler {
	return &Handler{adminService: adminService, pushService: pushService}
}
