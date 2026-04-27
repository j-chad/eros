package client

import (
	"backend/internal/handler/utils"
	"backend/internal/models"
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"fmt"
	"net/http"
)

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterDeviceRequest
	if err := utils.SafeJSONDecode(r.Body, &req); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(r.Context(), w, err)
		return
	}

	registrationCode := req.RegistrationCode
	clientDeviceInfo := req.DeviceInfo

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
