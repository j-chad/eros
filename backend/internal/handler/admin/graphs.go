package admin

import (
	"backend/internal/models"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"encoding/json"
	"net/http"
)

func (h *Handler) ListStartNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.adminService.ListStartNodes(r.Context())
	if err != nil {
		http.Error(w, "Failed to list start nodes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, nodes)
}

func (h *Handler) DeleteGraph(w http.ResponseWriter, r *http.Request) {
	startNodeID := r.PathValue("id")
	if startNodeID == "" {
		response.Error(w, apierror.BadRequest("Start node ID is required"))
		return
	}

	if err := h.adminService.DeleteGraph(r.Context(), startNodeID); err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateGraph(w http.ResponseWriter, r *http.Request) {
	var newGraphRequest models.NewGraphRequest
	if err := json.NewDecoder(r.Body).Decode(&newGraphRequest); err != nil {
		response.Error(w, apierror.BadRequest("invalid request body"))
	}

	if newGraphRequest.Title == "" {
		response.Error(w, apierror.BadRequest("graph title is required"))
		return
	}

	createdGraphID, err := h.adminService.CreateGraph(r.Context(), newGraphRequest)
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, createdGraphID)
}

func (h *Handler) GetGraph(w http.ResponseWriter, r *http.Request) {
	startNodeID := r.PathValue("id")
	if startNodeID == "" {
		response.Error(w, apierror.BadRequest("Start node ID is required"))
		return
	}

	graph, err := h.adminService.GetGraph(r.Context(), startNodeID)
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, graph)
}
