package client

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) ListFavourChoices(w http.ResponseWriter, r *http.Request) {
	choices, err := h.favourService.ListFavourChoices(r.Context())
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, choices)
}

func (h *Handler) GetFavourCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.favourService.GetFavourCount(r.Context())
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, count)
}
