package client

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

func (h *Handler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.graphService.GetGraph(r.Context(), r.PathValue("id"))
	if err != nil {
		response.Error(w, apierror.UnknownInternalError(err))
		return
	}

	if graph == nil {
		response.Error(w, apierror.NotFound("graph not found"))
		return
	}

	response.JSON(w, http.StatusOK, graph)
}
