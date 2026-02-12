package client

import (
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"encoding/json"
	"net/http"
)

func (h *Handler) ListFavourChoices(w http.ResponseWriter, r *http.Request) {
	choices, err := h.favourService.ListFavourChoices(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, choices)
}

func (h *Handler) GetFavourCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.favourService.GetFavourCount(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, count)
}

func (h *Handler) RequestFavour(w http.ResponseWriter, r *http.Request) {
	var favour models.FavourRequest

	if err := json.NewDecoder(r.Body).Decode(&favour); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	if favour.ChoiceID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("favour choice ID is required"))
		return
	}

	if err := h.favourService.RequestFavour(r.Context(), &favour); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, favour)
}

func (h *Handler) ListFavourRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.favourService.ListFavourRequests(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, requests)
}

func (h *Handler) DeleteFavourRequest(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("id")
	if requestID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("favour request ID is required"))
		return
	}

	if err := h.favourService.DeleteFavourRequest(r.Context(), requestID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
