package admin

import (
	"backend/internal/handler/utils"
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) SendPushMessage(w http.ResponseWriter, r *http.Request) {
	var sendRequest models.PushRequest
	if err := utils.SafeJSONDecode(r.Body, &sendRequest); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	if sendRequest.Message.Title == "" {
		response.Error(r.Context(), w, apierror.BadRequest("title is required"))
	}
	if sendRequest.Message.Body == "" {
		response.Error(r.Context(), w, apierror.BadRequest("body is required"))
	}

	result, err := h.pushService.SendMessage(r.Context(), sendRequest)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(r.Context(), w, http.StatusOK, result)
}

func (h *Handler) ListPushSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.pushService.ListSubscriptions(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(r.Context(), w, http.StatusOK, subs)
}

func (h *Handler) DeletePushSubscription(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("device_id is required"))
		return
	}

	if err := h.pushService.Unregister(r.Context(), deviceID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}
