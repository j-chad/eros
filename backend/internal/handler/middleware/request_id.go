package middleware

import (
	"backend/internal/crypto"
	"backend/internal/logging"
	"encoding/hex"
	"log/slog"
	"net/http"
)

const RequestIDHeader = "X-Request-ID"

func generateRequestID() string {
	bytes := crypto.RandomBytes(4)
	return hex.EncodeToString(bytes)
}

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		w.Header().Add(RequestIDHeader, requestID)

		logger := logging.FromContext(r.Context())
		requestLogger := logger.With(slog.String("request_id", requestID))
		r = r.WithContext(logging.NewContext(r.Context(), requestLogger))

		next.ServeHTTP(w, r)
	})
}
