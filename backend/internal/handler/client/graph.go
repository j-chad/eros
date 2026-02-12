package client

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) ListGraphs(w http.ResponseWriter, r *http.Request) {
	graphs, err := h.graphService.ListGraphs(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, graphs)
}

func (h *Handler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.graphService.GetGraph(r.Context(), r.PathValue("id"))
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	if graph == nil {
		response.Error(r.Context(), w, apierror.NotFound("graph not found"))
		return
	}

	response.JSON(w, http.StatusOK, graph)
}
