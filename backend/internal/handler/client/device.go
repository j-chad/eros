package client

import (
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"net/http"
)

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	registrationCode := r.FormValue("registration_code")
	if registrationCode == "" {
		response.Error(w, apierror.BadRequest("registration_code is required"))
		return
	}

	deviceInfo := r.FormValue("device_info")

	token, err := h.authService.RegisterDevice(r.Context(), registrationCode, deviceInfo)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRegistrationCode) {
			response.Error(w, apierror.Unauthorized("invalid registration code"))
		} else {
			response.Error(w, apierror.UnknownInternalError(err))
		}

		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"token": token})
}
