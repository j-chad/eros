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
	mux    *http.ServeMux
	auth   *service.AuthService
	admin  *admin.Handler
	client *client.Handler
}

func NewHandler(
	authService *service.AuthService,
	adminService *service.AdminService,
) *Handler {
	handler := &Handler{
		mux:    http.NewServeMux(),
		auth:   authService,
		admin:  admin.NewHandler(adminService),
		client: client.NewHandler(authService),
	}

	handler.registerRoutes()
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	h.mux.HandleFunc("POST /api/device", h.client.RegisterDevice)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("POST /api/admin/registration-codes", h.admin.CreateRegistrationCode)
	adminMux.HandleFunc("GET /api/admin/registration-codes", h.admin.GetRegistrationCode)
	adminMux.HandleFunc("DELETE /api/admin/registration-codes", h.admin.InvalidateRegistrationCode)
	adminMux.HandleFunc("/", routeNotFound)
	h.mux.Handle("/api/admin/", withAdminAuth(adminMux, *h.auth))

	clientMux := http.NewServeMux()
	clientMux.HandleFunc("GET /api/test", h.client.TestHandler)
	h.mux.Handle("/api/", withClientAuth(clientMux, *h.auth))

	h.mux.HandleFunc("/", routeNotFound)
}

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	response.Error(w, apierror.NotFound("no such route").
		WithDetail("path", r.URL.Path).
		WithDetail("method", r.Method))
}
