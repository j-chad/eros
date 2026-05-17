package middleware

import (
	"backend/internal/config"
	"net/http"
)

func WithCORS(next http.Handler, conf config.CORSConfig) http.Handler {
	wildcard := false
	for _, o := range conf.AllowedOrigins {
		if o == "*" {
			wildcard = true
			break
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if wildcard {
			// Wildcard mode: allow all origins but never send credentials.
			// Browsers refuse Access-Control-Allow-Credentials with a wildcard origin.
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if isOriginAllowed(origin, conf.AllowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		}

		if conf.AllowPrivateNetwork {
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
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

// isOriginAllowed checks if an origin is in the allowed list.
func isOriginAllowed(origin string, allowed []string) bool {
	for _, allowedOrigin := range allowed {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}
