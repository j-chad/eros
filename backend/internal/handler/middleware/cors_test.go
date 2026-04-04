package middleware

import (
	"backend/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name    string
		origin  string
		allowed []string
		want    bool
	}{
		{"exact match", "http://localhost:3000", []string{"http://localhost:3000"}, true},
		{"no match", "http://evil.com", []string{"http://localhost:3000"}, false},
		{"wildcard", "http://anything.com", []string{"*"}, true},
		{"empty origin", "", []string{"http://localhost:3000"}, false},
		{"empty allowlist", "http://localhost:3000", []string{}, false},
		{"nil allowlist", "http://localhost:3000", nil, false},
		{"multiple allowed", "http://b.com", []string{"http://a.com", "http://b.com"}, true},
		{"wildcard with others", "http://x.com", []string{"http://a.com", "*"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.Equal(t, isOriginAllowed(tt.origin, tt.allowed), tt.want)
		})
	}
}

func TestWithCORS_SetsHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := WithCORS(inner, []string{"http://app.com"})

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://app.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	testutil.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "http://app.com")
	testutil.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "true")
	testutil.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	testutil.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type, Authorization")
	testutil.Equal(t, w.Header().Get("Access-Control-Max-Age"), "172800")
	testutil.Equal(t, w.Code, http.StatusOK)
}

func TestWithCORS_DisallowedOrigin(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := WithCORS(inner, []string{"http://app.com"})

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	testutil.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "")
	testutil.Equal(t, w.Code, http.StatusOK)
}

func TestWithCORS_Preflight(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := WithCORS(inner, []string{"http://app.com"})

	req := httptest.NewRequest("OPTIONS", "/api/test", nil)
	req.Header.Set("Origin", "http://app.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	testutil.Equal(t, w.Code, http.StatusNoContent)
	testutil.False(t, called, "inner handler should not be called for preflight")
}
