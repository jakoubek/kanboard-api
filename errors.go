package kanboard

import (
	"errors"
	"fmt"
)

var (
	// ErrConnectionFailed indicates a connection to the Kanboard server failed.
	ErrConnectionFailed = errors.New("connection to Kanboard server failed")

	// ErrUnauthorized indicates authentication failed.
	ErrUnauthorized = errors.New("authentication failed: invalid credentials")

	// ErrForbidden indicates insufficient permissions.
	ErrForbidden = errors.New("access forbidden: insufficient permissions")
)

// APIError represents an error returned by the Kanboard API.
type APIError struct {
	Code    int
	Message string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("Kanboard API error (code %d): %s", e.Code, e.Message)
}
