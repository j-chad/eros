package admin

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"mime"
	"net/http"
	"path/filepath"
)

const FileSizeLimit = 10 * 1024 * 1024 // 10 MB

func (h *Handler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("node_id")
	if nodeID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("node_id query parameter is required"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(r.Context(), w, apierror.BadRequest("file is required"))
		return
	}
	defer file.Close()

	if header.Size == 0 {
		response.Error(r.Context(), w, apierror.BadRequest("file cannot be empty"))
		return
	}

	if header.Size > FileSizeLimit {
		response.Error(r.Context(), w, apierror.BadRequest("file size exceeds 10 MB limit"))
		return
	}

	if header.Filename == "" {
		response.Error(r.Context(), w, apierror.BadRequest("file name cannot be empty"))
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(header.Filename))
		if mimeType == "" {
			response.Error(r.Context(), w, apierror.BadRequest("unable to determine file MIME type"))
			return
		}
	}

	fileModel, err := h.adminService.UploadFile(r.Context(), nodeID, header.Filename, mimeType, header.Size, file)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusCreated, fileModel)
}

func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("node_id")
	if nodeID == "" {
		response.Error(r.Context(), w, apierror.BadRequest("node_id query parameter is required"))
		return
	}

	files, err := h.adminService.ListFiles(r.Context(), nodeID)
	if err != nil {
		response.Error(r.Context(), w, apierror.UnknownInternalError(err))
		return
	}

	response.JSON(w, http.StatusOK, files)
}
