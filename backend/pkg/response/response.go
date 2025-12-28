package response

import (
	"backend/pkg/apierror"
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

// Error writes an error response
func Error(w http.ResponseWriter, err *apierror.APIError) {
	if err.StatusCode >= 500 {
		log.Printf("internal error: %v", err.Err)
	}

	// Build error response
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    err.Code,
			"message": err.Message,
		},
	}

	// Add details if present
	if err.Details != nil && len(err.Details) > 0 {
		errorResponse["error"].(map[string]interface{})["details"] = err.Details
	}

	JSON(w, err.StatusCode, errorResponse)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
