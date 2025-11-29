package api

import "fmt"

// PermisError represents an error from the Permissio.io API.
type PermisError struct {
	// Message is the error message.
	Message string

	// Code is the error code (if provided).
	Code string

	// StatusCode is the HTTP status code.
	StatusCode int

	// Details contains additional error details.
	Details map[string]interface{}
}

// Error implements the error interface.
func (e *PermisError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s (status: %d)", e.Code, e.Message, e.StatusCode)
	}
	return fmt.Sprintf("%s (status: %d)", e.Message, e.StatusCode)
}

// IsNotFound returns true if this is a 404 error.
func (e *PermisError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if this is a 401 error.
func (e *PermisError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if this is a 403 error.
func (e *PermisError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsBadRequest returns true if this is a 400 error.
func (e *PermisError) IsBadRequest() bool {
	return e.StatusCode == 400
}

// IsServerError returns true if this is a 5xx error.
func (e *PermisError) IsServerError() bool {
	return e.StatusCode >= 500
}

// NewPermisError creates a new PermisError.
func NewPermisError(message, code string, statusCode int) *PermisError {
	return &PermisError{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

// AccessDeniedError creates an access denied error.
func AccessDeniedError(message string) *PermisError {
	return &PermisError{
		Message:    message,
		Code:       "ACCESS_DENIED",
		StatusCode: 403,
	}
}
