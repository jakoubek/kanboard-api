package kanboard

import (
	"errors"
	"fmt"
)

// Network errors
var (
	// ErrConnectionFailed indicates a connection to the Kanboard server failed.
	ErrConnectionFailed = errors.New("connection to Kanboard server failed")

	// ErrTimeout indicates a request timed out.
	ErrTimeout = errors.New("request timed out")

	// ErrTooManyRedirects indicates the server returned too many redirects.
	ErrTooManyRedirects = errors.New("too many redirects")
)

// Authentication errors
var (
	// ErrUnauthorized indicates authentication failed.
	ErrUnauthorized = errors.New("authentication failed: invalid credentials")

	// ErrForbidden indicates insufficient permissions.
	ErrForbidden = errors.New("access forbidden: insufficient permissions")
)

// Resource errors
var (
	// ErrNotFound indicates a resource was not found.
	ErrNotFound = errors.New("resource not found")

	// ErrProjectNotFound indicates the specified project was not found.
	ErrProjectNotFound = errors.New("project not found")

	// ErrTaskNotFound indicates the specified task was not found.
	ErrTaskNotFound = errors.New("task not found")

	// ErrColumnNotFound indicates the specified column was not found.
	ErrColumnNotFound = errors.New("column not found")

	// ErrCommentNotFound indicates the specified comment was not found.
	ErrCommentNotFound = errors.New("comment not found")

	// ErrCategoryNotFound indicates the specified category was not found.
	ErrCategoryNotFound = errors.New("category not found")
)

// Logic errors
var (
	// ErrAlreadyInLastColumn indicates a task is already in the last column.
	ErrAlreadyInLastColumn = errors.New("task is already in the last column")

	// ErrAlreadyInFirstColumn indicates a task is already in the first column.
	ErrAlreadyInFirstColumn = errors.New("task is already in the first column")

	// ErrTaskClosed indicates a task is already closed.
	ErrTaskClosed = errors.New("task is already closed")

	// ErrTaskOpen indicates a task is already open.
	ErrTaskOpen = errors.New("task is already open")
)

// Validation errors
var (
	// ErrEmptyTitle indicates a task title cannot be empty.
	ErrEmptyTitle = errors.New("task title cannot be empty")

	// ErrInvalidProjectID indicates an invalid project ID was provided.
	ErrInvalidProjectID = errors.New("invalid project ID")
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

// IsNotFound returns true if the error indicates a resource was not found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrProjectNotFound) ||
		errors.Is(err, ErrTaskNotFound) ||
		errors.Is(err, ErrColumnNotFound) ||
		errors.Is(err, ErrCommentNotFound) ||
		errors.Is(err, ErrCategoryNotFound)
}

// IsUnauthorized returns true if the error indicates an authentication failure.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsAPIError returns true if the error is an APIError from the Kanboard API.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}
