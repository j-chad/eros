package middleware

import "net/http"

const ctfHeader = "X-Request-Trace"
const ctfFlag = "eros{h34d3rs_sp34k_l0ud3r}"

// WithCTF adds headers that expose a CTF flag to the user
func WithCTF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ctfHeader, ctfFlag)
		next.ServeHTTP(w, r)
	})
}
