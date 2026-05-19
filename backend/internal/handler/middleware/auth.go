package middleware

import (
	"backend/internal/logging"
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/authctx"
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

// WithClientAuth wraps a handler with client authentication.
// Admins are also allowed to access these endpoints.
func WithClientAuth(next http.Handler, authService service.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		authHeader := r.Header.Get("Authorization")
		authType, token, ok := parseAuthHeader(authHeader)
		if !ok {
			response.Error(ctx, w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if authType == adminAuthType {
			WithAdminAuth(next, authService).ServeHTTP(w, r)
			return
		}

		device, err := authService.ValidateDeviceToken(ctx, token)
		if err != nil {
			if errors.Is(err, service.ErrInvalidClientCredentials) {
				logger.DebugContext(ctx, "invalid client token", "token", token)
				response.Error(ctx, w, apierror.Unauthorized("invalid credentials"))
			} else {
				logger.ErrorContext(ctx, "error validating client token", "error", err)
				response.Error(ctx, w, apierror.UnknownInternalError(err))
			}
			return
		}

		deviceIDLogger := logger.With("device_id", device.ID)
		deviceIDCtx := authctx.WithDeviceID(ctx, device.ID)
		logCtx := logging.NewContext(deviceIDCtx, deviceIDLogger)
		r = r.WithContext(logCtx)

		next.ServeHTTP(w, r)
	})
}

// WithAdminAuth wraps a handler with admin authentication
func WithAdminAuth(next http.Handler, authService service.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		authHeader := r.Header.Get("Authorization")
		authType, token, ok := parseAuthHeader(authHeader)
		if !ok || authType != adminAuthType {
			logger.WarnContext(ctx, "invalid admin authorization header")
			response.Error(r.Context(), w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if err := authService.ValidateAdminToken(token); err != nil {
			logger.WarnContext(ctx, "invalid admin token", "token", token)
			response.Error(r.Context(), w, apierror.Unauthorized("invalid credentials"))
			return
		}

		adminLogger := logger.With("admin", true)
		adminCtx := authctx.WithAdmin(ctx)
		logCtx := logging.NewContext(adminCtx, adminLogger)
		r = r.WithContext(logCtx)

		next.ServeHTTP(w, r)
	})
}
