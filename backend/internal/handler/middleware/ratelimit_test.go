package middleware

import (
	"backend/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPRateLimit_AllowsUnderLimit(t *testing.T) {
	limiter := NewIPRateLimit(3)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for range 3 {
		req := httptest.NewRequest("POST", "/api/device", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		testutil.Equal(t, rec.Code, http.StatusOK)
	}
}

func TestIPRateLimit_BlocksOverLimit(t *testing.T) {
	limiter := NewIPRateLimit(3)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for range 3 {
		req := httptest.NewRequest("POST", "/api/device", nil)
		req.RemoteAddr = "10.0.0.1:9999"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	// 4th request should be rejected
	req := httptest.NewRequest("POST", "/api/device", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	testutil.Equal(t, rec.Code, http.StatusTooManyRequests)
}

func TestIPRateLimit_IndependentBucketsPerIP(t *testing.T) {
	limiter := NewIPRateLimit(1)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust IP A
	req := httptest.NewRequest("POST", "/api/device", nil)
	req.RemoteAddr = "1.1.1.1:1000"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	testutil.Equal(t, rec.Code, http.StatusOK)

	// IP A is now blocked
	req = httptest.NewRequest("POST", "/api/device", nil)
	req.RemoteAddr = "1.1.1.1:1000"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	testutil.Equal(t, rec.Code, http.StatusTooManyRequests)

	// IP B should still work
	req = httptest.NewRequest("POST", "/api/device", nil)
	req.RemoteAddr = "2.2.2.2:2000"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
}
