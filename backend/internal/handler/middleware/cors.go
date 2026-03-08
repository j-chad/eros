package middleware

import "net/http"

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

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin string, allowed []string) bool {
	for _, allowedOrigin := range allowed {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}
