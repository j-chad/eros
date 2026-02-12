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
	favourService *service.FavourService,
	graphService *service.GraphService,
) *Handler {
	handler := &Handler{
		auth:   authService,
		admin:  admin.NewHandler(adminService),
		client: client.NewHandler(authService, favourService, graphService),
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
	adminMux.HandleFunc("GET /api/admin/devices", h.admin.ListDevices)
	adminMux.HandleFunc("DELETE /api/admin/devices/{id}", h.admin.RevokeDevice)
	adminMux.HandleFunc("PATCH /api/admin/devices/{id}", h.admin.UpdateDeviceInfo)
	adminMux.HandleFunc("POST /api/admin/favours/choices", h.admin.CreateFavourChoice)
	adminMux.HandleFunc("PUT /api/admin/favours/choices/{id}", h.admin.UpdateFavourChoice)
	adminMux.HandleFunc("DELETE /api/admin/favours/choices/{id}", h.admin.DeleteFavourChoice)
	adminMux.HandleFunc("PUT /api/admin/favours", h.admin.UpdateFavourCount)
	adminMux.HandleFunc("PATCH /api/admin/favours/requests/{id}", h.admin.UpdateFavourRequestStatus)
	adminMux.HandleFunc("GET /api/admin/graphs", h.admin.ListStartNodes)
	adminMux.HandleFunc("DELETE /api/admin/graphs/{id}", h.admin.DeleteGraph)
	adminMux.HandleFunc("POST /api/admin/graphs", h.admin.CreateGraph)
	adminMux.HandleFunc("GET /api/admin/graphs/{id}", h.admin.GetGraph)
	adminMux.HandleFunc("PUT /api/admin/graphs/{id}", h.admin.UpdateGraph)
	adminMux.HandleFunc("/", routeNotFound)
	mux.Handle("/api/admin/", withAdminAuth(adminMux, *h.auth))

	clientMux := http.NewServeMux()
	clientMux.HandleFunc("GET /api/favours/choices", h.client.ListFavourChoices)
	clientMux.HandleFunc("GET /api/favours", h.client.GetFavourCount)
	clientMux.HandleFunc("POST /api/favours/request", h.client.RequestFavour)
	clientMux.HandleFunc("GET /api/favours/requests", h.client.ListFavourRequests)
	clientMux.HandleFunc("DELETE /api/favours/request/{id}", h.client.DeleteFavourRequest)
	clientMux.HandleFunc("GET /api/graphs", h.client.ListGraphs)
	clientMux.HandleFunc("GET /api/graphs/{id}", h.client.GetGraph)
	mux.Handle("/api/", withClientAuth(clientMux, *h.auth))

	mux.HandleFunc("/", routeNotFound)
}

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	response.Error(r.Context(), w, apierror.NotFound("no such route").
		WithDetail("path", r.URL.Path).
		WithDetail("method", r.Method))
}
