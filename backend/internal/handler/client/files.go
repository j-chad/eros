package client

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"io"
	"net/http"
)

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fileID := r.PathValue("fileID")

	file, err := h.fileService.GetFile(ctx, fileID)
	if err != nil {
		response.Error(ctx, w, apierror.UnknownInternalError(err))
		return
	}
	if file == nil {
		response.Error(ctx, w, apierror.NotFound("file not found"))
		return
	}

	// Access control: check whether the client can reach the owning node.
	node, err := h.graphService.GetAccessibleNode(ctx, file.NodeID)
	if err != nil {
		response.Error(ctx, w, apierror.UnknownInternalError(err))
		return
	}
	if node == nil {
		// Return 404 (not 403) to avoid revealing file existence.
		response.Error(ctx, w, apierror.NotFound("file not found"))
		return
	}

	// If S3, redirect to presigned URL.
	if h.fileService.IsPresignCapable() {
		presigned, err := h.fileService.PresignURL(ctx, file.StorageKey)
		if err != nil {
			response.Error(ctx, w, apierror.UnknownInternalError(err))
			return
		}
		http.Redirect(w, r, presigned, http.StatusFound)
		return
	}

	// Local storage: stream the file.
	reader, _, err := h.fileService.GetFileStream(ctx, fileID)
	if err != nil {
		response.Error(ctx, w, apierror.UnknownInternalError(err))
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+file.Filename+"\"")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	w.WriteHeader(http.StatusOK)

	io.Copy(w, reader)
}
