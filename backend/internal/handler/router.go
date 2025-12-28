package handler

import (
	"backend/internal/handler/admin"
	"backend/internal/service"
	"net/http"
)

type Handler struct {
	mux   *http.ServeMux
	admin *admin.Handler
}

func NewHandler(
	adminService *service.AdminService,
) *Handler {
	return &Handler{
		mux:   http.NewServeMux(),
		admin: admin.NewHandler(adminService),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	adminMux := http.NewServeMux()
	//adminMux.HandleFunc("POST /api/admin/registration-codes", h.admin.CreateRegistrationCode)
	//adminMux.HandleFunc("GET /api/admin/registration-codes", h.admin.ListRegistrationCodes)

	mux.Handle("/api/admin/", withAdminAuth(adminMux))
}
