package handler

import (
	"backend/internal/handler/admin"
	"backend/internal/service"
	"net/http"
)

//func NewAppHandler(db repository.Repository) http.Handler {
//	mux := http.NewServeMux()
//	handler := &AppHandler{db: db, mux: mux}
//
//	// ADMIN ROUTES
//	// Device Management
//	mux.Handle("POST /admin/devices", http.HandlerFunc(handler.adminRegisterDeviceHandler))
//	mux.Handle("GET /admin/devices", http.HandlerFunc(handler.adminListDevicesHandler))
//	mux.Handle("DELETE /admin/devices/:id", http.HandlerFunc(handler.adminDeleteDeviceHandler))
//
//	return handler
//}

type Handler struct {
	admin *admin.Handler
}

func NewHandler(
	adminService *service.AdminService,
) *Handler {
	return &Handler{
		admin: admin.NewHandler(adminService),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("POST /api/admin/registration-codes", h.admin.CreateRegistrationCode)
	adminMux.HandleFunc("GET /api/admin/registration-codes", h.admin.ListRegistrationCodes)

	mux.Handle("/api/admin/", withAdminAuth(adminMux))
}
