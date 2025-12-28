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
