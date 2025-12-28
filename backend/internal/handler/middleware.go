package handler

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
	"crypto/subtle"
	"net/http"
	"os"
)

// withAdminAuth wraps a handler with admin authentication
func withAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		expectedKey := os.Getenv("ADMIN_API_KEY")

		if apiKey == "" || expectedKey == "" {
			response.Error(w, apierror.Unauthorized("api key not provided"))
		}

		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
			response.Error(w, apierror.Forbidden("invalid admin credentials"))
		}

		next.ServeHTTP(w, r)
	})
}
