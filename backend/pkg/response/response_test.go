package response

import (
	"backend/internal/testutil"
	"backend/pkg/apierror"
	"backend/pkg/authctx"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_StatusAndContentType(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(context.Background(), w, http.StatusCreated, map[string]string{"id": "123"})

	testutil.Equal(t, w.Code, http.StatusCreated)
	testutil.Equal(t, w.Header().Get("Content-Type"), "application/json")
}

func TestJSON_Body(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(context.Background(), w, http.StatusOK, map[string]int{"count": 42})

	var body map[string]int
	testutil.NilErr(t, json.NewDecoder(w.Body).Decode(&body))
	testutil.Equal(t, body["count"], 42)
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	NoContent(w)
	testutil.Equal(t, w.Code, http.StatusNoContent)
	testutil.Equal(t, w.Body.Len(), 0)
}

func TestError_BasicFields(t *testing.T) {
	w := httptest.NewRecorder()
	Error(context.Background(), w, apierror.BadRequest("missing field"))

	testutil.Equal(t, w.Code, http.StatusBadRequest)

	var body map[string]map[string]any
	testutil.NilErr(t, json.NewDecoder(w.Body).Decode(&body))
	testutil.Equal(t, body["error"]["code"], "BAD_REQUEST")
	testutil.Equal(t, body["error"]["message"], "missing field")
}

func TestError_WithDetails(t *testing.T) {
	w := httptest.NewRecorder()
	Error(context.Background(), w, apierror.BadRequest("invalid").WithDetail("field", "email"))

	var body map[string]map[string]any
	json.NewDecoder(w.Body).Decode(&body)
	details := body["error"]["details"].(map[string]any)
	testutil.Equal(t, details["field"], "email")
}

func TestError_HidesInternalErrorFromClient(t *testing.T) {
	w := httptest.NewRecorder()
	Error(context.Background(), w, apierror.UnknownInternalError(fmt.Errorf("db connection failed")))

	var body map[string]map[string]any
	json.NewDecoder(w.Body).Decode(&body)
	_, hasInternal := body["error"]["internal"]
	testutil.False(t, hasInternal, "internal error should not be exposed to non-admin")
}

func TestError_ExposesInternalErrorToAdmin(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := authctx.WithAdmin(context.Background())
	Error(ctx, w, apierror.UnknownInternalError(fmt.Errorf("db connection failed")))

	var body map[string]map[string]any
	json.NewDecoder(w.Body).Decode(&body)
	testutil.Equal(t, body["error"]["internal"], "db connection failed")
}
