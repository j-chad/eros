package handler

import (
	"backend/internal/handler/admin"
	"backend/internal/handler/client"
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

type Handler struct {
	auth   *service.AuthService
	admin  *admin.Handler
	client *client.Handler
}

func NewHandler(
	authService *service.AuthService,
	adminService *service.AdminService,
) *Handler {
	handler := &Handler{
		auth:   authService,
		admin:  admin.NewHandler(adminService),
		client: client.NewHandler(authService),
	}
	return handler
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/device", h.client.RegisterDevice)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /api/admin/ping", h.admin.Ping)
	adminMux.HandleFunc("POST /api/admin/registration-codes", h.admin.CreateRegistrationCode)
	adminMux.HandleFunc("GET /api/admin/registration-codes", h.admin.GetRegistrationCode)
	adminMux.HandleFunc("DELETE /api/admin/registration-codes", h.admin.InvalidateRegistrationCode)
	adminMux.HandleFunc("/", routeNotFound)
	mux.Handle("/api/admin/", withAdminAuth(adminMux, *h.auth))

	clientMux := http.NewServeMux()
	clientMux.HandleFunc("GET /api/test", h.client.TestHandler)
	mux.Handle("/api/", withClientAuth(clientMux, *h.auth))

	mux.HandleFunc("/", routeNotFound)
}

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	response.Error(w, apierror.NotFound("no such route").
		WithDetail("path", r.URL.Path).
		WithDetail("method", r.Method))
}
