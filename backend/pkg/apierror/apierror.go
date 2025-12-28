package apierror

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"` // Not serialized to JSON
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"` // Original error, not exposed
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// WithDetails adds additional context to the error
func (e *APIError) WithDetails(details map[string]interface{}) *APIError {
	e.Details = details
	return e
}

// WithDetail adds a single detail field
func (e *APIError) WithDetail(key string, value interface{}) *APIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// Unauthorized creates a 401 error
func Unauthorized(message string) *APIError {
	return &APIError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func UnknownInternalError(err error) *APIError {
	return &APIError{
		Code:       "UNKNOWN_INTERNAL_ERROR",
		Message:    "An unknown internal error occurred",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func Forbidden(message string) *APIError {
	return &APIError{
		Code:       "FORBIDDEN",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NotFound(message string) *APIError {
	return &APIError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func BadRequest(message string) *APIError {
	return &APIError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
