package admin

import (
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
	"time"
)

func (h *Handler) CreateFavourChoice(w http.ResponseWriter, r *http.Request) {
	label := r.FormValue("label")
	if label == "" {
		response.Error(w, apierror.BadRequest("label is required"))
		return
	}

	description := r.FormValue("description")
	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}

	canMessage := r.FormValue("can_message")
	canMessageBool := false
	if canMessage == "true" {
		canMessageBool = true
	}

	choice := models.FavourChoice{
		Label:       label,
		Description: descriptionPtr,
		CanMessage:  canMessageBool,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

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

	label := r.FormValue("label")
	if label == "" {
		response.Error(w, apierror.BadRequest("label is required"))
		return
	}

	description := r.FormValue("description")
	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}

	canMessage := r.FormValue("can_message")
	canMessageBool := false
	if canMessage == "true" {
		canMessageBool = true
	}

	choice := models.FavourChoice{
		ID:          choiceID,
		Label:       label,
		Description: descriptionPtr,
		CanMessage:  canMessageBool,
	}

	if err := h.adminService.UpdateFavourChoice(r.Context(), choice); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
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
