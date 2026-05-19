package client

import (
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"io"
	"net/http"
)

func (h *Handler) ListGraphs(w http.ResponseWriter, r *http.Request) {
	graphs, err := h.graphService.ListGraphs(r.Context())
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(r.Context(), w, http.StatusOK, graphs)
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

	response.JSON(r.Context(), w, http.StatusOK, graph)
}

func (h *Handler) UnlockNode(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 1024)) // limit to 1KB to prevent abuse
	if err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("failed to read request body"))
		return
	}

	result, err := h.graphService.UnlockNode(r.Context(), r.PathValue("id"), string(payload))
	if err != nil {
		if errors.Is(err, service.ErrNodeUnlockIncorrect) {
			response.Error(r.Context(), w, apierror.Forbidden("incorrect unlock payload"))
			return
		}

		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(r.Context(), w, http.StatusOK, result)
}
