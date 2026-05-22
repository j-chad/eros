package admin

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"io"
	"net/http"
)

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.authService.ListDevices(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(r.Context(), w, http.StatusOK, devices)
}

func (h *Handler) RevokeDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	if deviceID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("device ID is required"))
		return
	}

	if err := h.authService.RevokeDevice(r.Context(), deviceID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}

func (h *Handler) UpdateDeviceInfo(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	if deviceID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("device ID is required"))
		return
	}

	// interpret device info as raw string from body
	const maxSize = 1024 // 1KB
	buffer, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
	if err != nil && !errors.Is(err, io.EOF) {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err).WithDetail("msg", "failed to read device info"))
		return
	}
	deviceInfo := string(buffer)

	if err := h.authService.UpdateDeviceInfo(r.Context(), deviceID, deviceInfo); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}
