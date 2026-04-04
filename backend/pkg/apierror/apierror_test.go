package apierror

import (
	"backend/internal/testutil"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestBadRequest(t *testing.T) {
	err := BadRequest("bad input")
	testutil.Equal(t, err.Code, "BAD_REQUEST")
	testutil.Equal(t, err.Message, "bad input")
	testutil.Equal(t, err.StatusCode, http.StatusBadRequest)
}

func TestNotFound(t *testing.T) {
	err := NotFound("not here")
	testutil.Equal(t, err.Code, "NOT_FOUND")
	testutil.Equal(t, err.StatusCode, http.StatusNotFound)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("nope")
	testutil.Equal(t, err.Code, "UNAUTHORIZED")
	testutil.Equal(t, err.StatusCode, http.StatusUnauthorized)
}

func TestForbidden(t *testing.T) {
	err := Forbidden("denied")
	testutil.Equal(t, err.Code, "FORBIDDEN")
	testutil.Equal(t, err.StatusCode, http.StatusForbidden)
}

func TestUnknownInternalError_WrapsPlainError(t *testing.T) {
	cause := fmt.Errorf("db on fire")
	err := UnknownInternalError(cause)

	testutil.Equal(t, err.Code, "UNKNOWN_INTERNAL_ERROR")
	testutil.Equal(t, err.StatusCode, http.StatusInternalServerError)
	testutil.True(t, errors.Is(err, cause), "wrapped error should be the original cause")
}

func TestUnknownInternalError_UnwrapsExistingAPIError(t *testing.T) {
	original := NotFound("graph not found")
	wrapped := fmt.Errorf("service layer: %w", original)

	result := UnknownInternalError(wrapped)
	testutil.Equal(t, result, original)
	testutil.Equal(t, result.StatusCode, http.StatusNotFound)
}

func TestError_WithWrappedError(t *testing.T) {
	err := BadRequest("bad")
	err.Err = fmt.Errorf("boom")
	testutil.Equal(t, err.Error(), "bad: boom")
}

func TestError_WithoutWrappedError(t *testing.T) {
	err := BadRequest("just bad")
	testutil.Equal(t, err.Error(), "just bad")
}

func TestUnwrap(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := UnknownInternalError(cause)
	testutil.True(t, errors.Is(err, cause), "errors.Is should find the wrapped cause")
}

func TestWithDetail_InitializesMap(t *testing.T) {
	err := BadRequest("bad").WithDetail("field", "name")
	testutil.NotNil(t, err.Details)
	testutil.Equal(t, err.Details["field"], "name")
}

func TestWithDetail_AddsToExistingMap(t *testing.T) {
	err := BadRequest("bad").WithDetail("a", 1).WithDetail("b", 2)
	testutil.Equal(t, len(err.Details), 2)
}

func TestWithDetails_ReplacesMap(t *testing.T) {
	err := BadRequest("bad").WithDetail("old", "value")
	err.WithDetails(map[string]any{"new": "value"})

	_, hasOld := err.Details["old"]
	testutil.False(t, hasOld, "old key should be gone after WithDetails")
	testutil.Equal(t, err.Details["new"], "value")
}

func TestWithDetail_Fluent(t *testing.T) {
	err := BadRequest("bad").WithDetail("a", 1).WithDetail("b", 2)
	testutil.Equal(t, err.Code, "BAD_REQUEST")
}
