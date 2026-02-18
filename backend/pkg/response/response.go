package response

import (
	"backend/pkg/apierror"
	"backend/pkg/authctx"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

// Error writes an error response
func Error(ctx context.Context, w http.ResponseWriter, err *apierror.APIError) {
	if err.StatusCode >= 500 {
		log.Printf("internal error: %v", err.Err)
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
