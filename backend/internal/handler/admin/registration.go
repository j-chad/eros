package admin

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) CreateRegistrationCode(w http.ResponseWriter, r *http.Request) {
	registrationCode, err := h.adminService.CreateRegistrationCode(r.Context())
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, registrationCode)
}

func (h *Handler) InvalidateRegistrationCode(w http.ResponseWriter, r *http.Request) {
	if err := h.adminService.InvalidateRegistrationCode(r.Context()); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}

func (h *Handler) GetRegistrationCode(w http.ResponseWriter, r *http.Request) {
	registrationCode, err := h.adminService.GetRegistrationCode(r.Context())
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	if registrationCode == nil {
		response.Error(w, apierror.NotFound("no active registration code found"))
		return
	}

	response.JSON(w, http.StatusOK, registrationCode)
}
