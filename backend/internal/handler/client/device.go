package client

import (
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"fmt"
	"net/http"
)

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	registrationCode := r.FormValue("registration_code")
	if registrationCode == "" {
		response.Error(r.Context(), w, apierror.BadRequest("registration_code is required"))
		return
	}

	clientDeviceInfo := r.FormValue("device_info")
	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr
	deviceInfo := fmt.Sprintf("UserAgent: %s; IPAddress: %s; Info: %s", userAgent, ipAddress, clientDeviceInfo)

	token, err := h.authService.RegisterDevice(r.Context(), registrationCode, deviceInfo)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRegistrationCode) {
			response.Error(r.Context(), w, apierror.Unauthorized("invalid registration code"))
		} else {
			response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		}

		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"token": token})
}
