package handler

import (
	"backend/internal/logger"
	"backend/internal/service"
	"backend/pkg/apierror"
	"backend/pkg/authctx"
	"backend/pkg/response"
	"errors"
	"fmt"
	"log/slog"
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
			response.Error(r.Context(), w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if authType == adminAuthType {
			withAdminAuth(next, authService).ServeHTTP(w, r)
			return
		}

		if err := authService.ValidateDeviceToken(token); err != nil {
			if errors.Is(err, service.ErrInvalidClientCredentials) {
				response.Error(r.Context(), w, apierror.Unauthorized("invalid credentials"))
			} else {
				response.Error(r.Context(), w, apierror.UnknownInternalError(err))
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
			response.Error(r.Context(), w, apierror.Unauthorized("invalid authorization header"))
			return
		}

		if err := authService.ValidateAdminToken(token); err != nil {
			response.Error(r.Context(), w, apierror.Unauthorized("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r.WithContext(authctx.WithAdmin(r.Context())))
	})
}

func WithCORS(next http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Allow credentials (cookies, authorization headers)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

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

// statusWriter wraps http.ResponseWriter to capture the status code.
type statusWriter struct {
	http.ResponseWriter
	code    int
	written bool
}

func (w *statusWriter) WriteHeader(code int) {
	if !w.written {
		w.code = code
		w.written = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.code = http.StatusOK
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

func WithTracing(headerName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var traceID string
		if traceID = r.Header.Get(headerName); traceID != "" {
			ctx = logger.WithExistingTrace(ctx, traceID)
		} else {
			ctx, traceID = logger.WithTrace(ctx)
		}

		ctx, span := logger.StartSpan(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		span.Logger().Debug("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"query", r.URL.RawQuery,
			"user_agent", r.UserAgent(),
		)

		w.Header().Set(headerName, traceID)

		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r.WithContext(ctx))

		level := slog.LevelInfo
		switch {
		case sw.code >= 500:
			level = slog.LevelError
		case sw.code >= 400:
			level = slog.LevelWarn
		}
		span.Logger().Log(ctx, level, "request completed", "status", sw.code)
		span.End()
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
