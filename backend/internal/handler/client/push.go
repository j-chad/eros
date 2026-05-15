package client

import (
	"backend/internal/handler/utils"
	"backend/internal/models"
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/authctx"
	"backend/pkg/response"
	"errors"
	"net/http"
)

func (h *Handler) SubscribePush(w http.ResponseWriter, r *http.Request) {
	var subscription models.PushSubscription

	if err := utils.SafeJSONDecode(r.Body, &subscription); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	deviceID, ok := authctx.DeviceID(r.Context())
	if !ok {
		response.Error(r.Context(), w, apierror.BadRequest("invalid device id"))
	}

	err := h.pushService.Register(r.Context(), deviceID, subscription)
	if err != nil {
		if errors.Is(err, service.ErrSubValidationFailed) {
			response.Error(r.Context(), w, apierror.BadRequest(err.Error()))
			return
		}

		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UnsubscribePush(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := authctx.DeviceID(r.Context())
	if !ok {
		response.Error(r.Context(), w, apierror.BadRequest("invalid device id"))
	}

	err := h.pushService.Unregister(r.Context(), deviceID)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetVAPIDPublicKey(w http.ResponseWriter, r *http.Request) {
	pubKey, err := h.pushService.GetVAPIDPublicKey()
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, pubKey)
}
