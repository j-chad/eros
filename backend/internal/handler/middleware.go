package handler

import (
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"errors"
	"fmt"
	"net/http"
)

const adminAuthType = "Admin"
const clientAuthType = "Bearer"

func parseAuthHeader(header string) (authType, token string, ok bool) {
	if header == "" {
		return "", "", false
	}

	n, err := fmt.Sscanf(header, "%s %s", &authType, &token)
	if err != nil || n != 2 {
		return "", "", false
	}

	if authType != adminAuthType && authType != clientAuthType {
		return "", "", false
	}

	return authType, token, true
}

// withClientAuth wraps a handler with client authentication.
// Admins are also allowed to access these endpoints.
func withClientAuth(next http.Handler, authService service.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authType, token, ok := parseAuthHeader(authHeader)
		if !ok {
			response.Error(w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if authType == adminAuthType {
			withAdminAuth(next, authService).ServeHTTP(w, r)
			return
		}

		if err := authService.ValidateDeviceToken(token); err != nil {
			if errors.Is(err, service.ErrInvalidClientCredentials) {
				response.Error(w, apierror.Unauthorized("invalid credentials"))
			} else {
				response.Error(w, apierror.UnknownInternalError(err))
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// withAdminAuth wraps a handler with admin authentication
func withAdminAuth(next http.Handler, authService service.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authType, token, ok := parseAuthHeader(authHeader)
		if !ok || authType != adminAuthType {
			response.Error(w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if err := authService.ValidateAdminToken(token); err != nil {
			response.Error(w, apierror.Unauthorized("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithCORS(next http.Handler, allowedOrigins []string) http.Handler {
	fmt.Println("CORS allowed origins:", allowedOrigins)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Allow credentials (cookies, authorization headers)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Max age for preflight cache (48 hours)
		w.Header().Set("Access-Control-Max-Age", "172800")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin string, allowed []string) bool {
	for _, allowedOrigin := range allowed {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}
