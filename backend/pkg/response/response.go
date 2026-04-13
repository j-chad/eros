package response

import (
	"backend/internal/logging"
	"backend/pkg/apierror"
	"backend/pkg/authctx"
	"context"
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger := logging.FromContext(context.Background())
		logger.Error("error encoding JSON response", "error", err)
	}
}

// Error writes an error response
func Error(ctx context.Context, w http.ResponseWriter, err *apierror.APIError) {
	if err.StatusCode >= 500 {
		logger := logging.FromContext(ctx)
		logger.Error("internal error", "error", err.Err, "code", err.Code, "message", err.Message, "details", err.Details)
	}

	// Build error response
	errorResponse := map[string]any{
		"error": map[string]any{
			"code":    err.Code,
			"message": err.Message,
		},
	}

	// Add details if present
	if err.Details != nil && len(err.Details) > 0 {
		errorResponse["error"].(map[string]any)["details"] = err.Details
	}

	// Add internal error if admin
	if authctx.IsAdmin(ctx) && err.Err != nil {
		errorResponse["error"].(map[string]any)["internal"] = err.Err.Error()
	}

	JSON(w, err.StatusCode, errorResponse)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
