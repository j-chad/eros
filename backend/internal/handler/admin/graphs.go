package admin

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

func (h *Handler) ListStartNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.adminService.ListGraphs(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(fmt.Errorf("failed to list start nodes: %w", err)))
		return
	}

	response.JSON(w, http.StatusOK, nodes)
}

func (h *Handler) DeleteGraph(w http.ResponseWriter, r *http.Request) {
	graphID := r.PathValue("id")
	if graphID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("Graph ID is required"))
		return
	}

	if err := h.adminService.DeleteGraph(r.Context(), graphID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateGraph(w http.ResponseWriter, r *http.Request) {
	var newGraphRequest models.NewGraphRequest
	if err := utils.SafeJSONDecode(r.Body, &newGraphRequest); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	if newGraphRequest.Title == "" {
		response.Error(r.Context(), w, apierror.BadRequest("graph title is required"))
		return
	}

	createdGraphID, err := h.adminService.CreateGraph(r.Context(), newGraphRequest)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, createdGraphID)
}

func (h *Handler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graphID := r.PathValue("id")
	if graphID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("Graph ID is required"))
		return
	}

	graph, err := h.adminService.GetGraph(r.Context(), graphID)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	if graph == nil {
		response.Error(r.Context(), w, apierror.NotFound("Graph not found"))
		return
	}

	response.JSON(w, http.StatusOK, graph)
}

func (h *Handler) UnlockNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("node_id")
	if nodeID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("Node ID is required"))
		return
	}

	if err := h.adminService.AdminUnlockNode(r.Context(), nodeID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}

func (h *Handler) LockNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("node_id")
	if nodeID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("Node ID is required"))
		return
	}

	if err := h.adminService.AdminLockNode(r.Context(), nodeID); err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.NoContent(w)
}

func (h *Handler) UpdateGraph(w http.ResponseWriter, r *http.Request) {
	var graphID = r.PathValue("id")
	if graphID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("Graph ID is required"))
		return
	}

	var graph models.Graph
	if err := utils.SafeJSONDecode(r.Body, &graph); err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("invalid request body"))
		return
	}

	graph.ID = graphID

	err := h.adminService.UpdateGraph(r.Context(), graph)
	if err != nil {
		if errors.Is(err, service.ErrInvalidGraph) {
			response.Error(r.Context(), w, apierror.BadRequest("invalid graph data"))
			return
		}

		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
