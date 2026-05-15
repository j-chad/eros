package client

import (
	"backend/internal/handler/utils"
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
	"net/url"
)

func (h *Handler) SubscribePush(w http.ResponseWriter, r *http.Request) {
	var subscription models.PushSubscription

	if err := utils.SafeJSONDecode(r.Body, &subscription); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	u, err := url.Parse(subscription.Endpoint)
	if err != nil || u.Host == "" || u.Scheme != "https" {
		response.Error(r.Context(), w, apierror.BadRequest("invalid endpoint"))
	}
}
