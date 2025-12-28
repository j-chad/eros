package handler

import (
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"net/http"
)

// withAdminAuth wraps a handler with admin authentication
func withAdminAuth(next http.Handler, authService service.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			response.Error(w, apierror.Unauthorized("api key not provided"))
			return
		}

		if err := authService.ValidateAdminAPIKey(apiKey); err != nil {
			response.Error(w, apierror.Unauthorized("invalid admin credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
