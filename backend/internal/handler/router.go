package handler

import (
	"backend/internal/repository"
	"net/http"
)

func NewAppHandler(db repository.Repository) http.Handler {
	mux := http.NewServeMux()
	handler := &AppHandler{db: db, mux: mux}

	// ADMIN ROUTES
	// Device Management
	mux.Handle("POST /admin/devices", http.HandlerFunc(handler.adminRegisterDeviceHandler))
	mux.Handle("GET /admin/devices", http.HandlerFunc(handler.adminListDevicesHandler))
	mux.Handle("DELETE /admin/devices/:id", http.HandlerFunc(handler.adminDeleteDeviceHandler))

	return handler
}

type AppHandler struct {
	db  repository.Repository
	mux http.Handler
}

func (h *AppHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.mux.ServeHTTP(writer, request)
}
