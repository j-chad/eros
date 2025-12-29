package admin

import (
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) CreateFavourChoice(w http.ResponseWriter, r *http.Request) {
	var choice models.FavourChoice

	if err := json.NewDecoder(r.Body).Decode(&choice); err != nil {
		response.Error(w, apierror.BadRequest("invalid request body"))
		return
	}

	if choice.Label == "" {
		response.Error(w, apierror.BadRequest("label is required"))
		return
	}

	if choice.Description != nil && *choice.Description == "" {
		choice.Description = nil
	}

	choice.CreatedAt = time.Now()
	choice.UpdatedAt = time.Now()

	if err := h.adminService.CreateFavourChoice(r.Context(), &choice); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, choice)
}

func (h *Handler) UpdateFavourChoice(w http.ResponseWriter, r *http.Request) {
	choiceID := r.PathValue("id")
	if choiceID == "" {
		response.Error(w, apierror.BadRequest("choice ID is required"))
		return
	}

	var choice models.FavourChoice
	if err := json.NewDecoder(r.Body).Decode(&choice); err != nil {
		response.Error(w, apierror.BadRequest("invalid request body"))
		return
	}

	if choice.Label == "" {
		response.Error(w, apierror.BadRequest("label is required"))
		return
	}

	if choice.Description != nil && *choice.Description == "" {
		choice.Description = nil
	}

	choice.ID = choiceID
	choice.UpdatedAt = time.Now()

	if err := h.adminService.UpdateFavourChoice(r.Context(), choice); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, choice)
}

func (h *Handler) DeleteFavourChoice(w http.ResponseWriter, r *http.Request) {
	choiceID := r.PathValue("id")
	if choiceID == "" {
		response.Error(w, apierror.BadRequest("choice ID is required"))
		return
	}

	if err := h.adminService.DeleteFavourChoice(r.Context(), choiceID); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}

func (h *Handler) UpdateFavourCount(w http.ResponseWriter, r *http.Request) {
	var count int
	if err := json.NewDecoder(r.Body).Decode(&count); err != nil {
		response.Error(w, apierror.BadRequest("invalid request body"))
		return
	}

	if count < 0 {
		response.Error(w, apierror.BadRequest("favour count cannot be negative"))
		return
	}

	if err := h.adminService.UpdateFavourCount(r.Context(), count); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}
