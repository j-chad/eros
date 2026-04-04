package middleware

import (
	"backend/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusWriter_CapturesStatus(t *testing.T) {
	w := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: w}

	sw.WriteHeader(http.StatusNotFound)
	testutil.Equal(t, sw.status, http.StatusNotFound)
}

func TestStatusWriter_DefaultsTo200OnWrite(t *testing.T) {
	w := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: w}

	sw.Write([]byte("hello"))
	testutil.Equal(t, sw.status, http.StatusOK)
}

func TestStatusWriter_PreservesExplicitStatus(t *testing.T) {
	w := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: w}

	sw.WriteHeader(http.StatusCreated)
	sw.Write([]byte("body"))
	testutil.Equal(t, sw.status, http.StatusCreated)
}
