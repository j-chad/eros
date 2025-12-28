package admin

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.adminService.ListDevices(r.Context())
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, devices)
}

func (h *Handler) RevokeDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	if deviceID == "" {
		response.Error(w, apierror.BadRequest("device ID is required"))
		return
	}

	if err := h.adminService.RevokeDevice(r.Context(), deviceID); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}
