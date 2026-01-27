package kanboard

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "invalid request",
			err:      &APIError{Code: -32600, Message: "Invalid Request"},
			expected: "Kanboard API error (code -32600): Invalid Request",
		},
		{
			name:     "method not found",
			err:      &APIError{Code: -32601, Message: "Method not found"},
			expected: "Kanboard API error (code -32601): Method not found",
		},
		{
			name:     "custom error",
			err:      &APIError{Code: 1001, Message: "Task not found"},
			expected: "Kanboard API error (code 1001): Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("APIError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrNotFound", ErrNotFound, true},
		{"ErrProjectNotFound", ErrProjectNotFound, true},
		{"ErrTaskNotFound", ErrTaskNotFound, true},
		{"ErrColumnNotFound", ErrColumnNotFound, true},
		{"ErrCommentNotFound", ErrCommentNotFound, true},
		{"wrapped ErrNotFound", fmt.Errorf("context: %w", ErrNotFound), true},
		{"wrapped ErrTaskNotFound", fmt.Errorf("getting task: %w", ErrTaskNotFound), true},
		{"ErrUnauthorized", ErrUnauthorized, false},
		{"ErrConnectionFailed", ErrConnectionFailed, false},
		{"generic error", errors.New("some error"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.expected {
				t.Errorf("IsNotFound(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrUnauthorized", ErrUnauthorized, true},
		{"wrapped ErrUnauthorized", fmt.Errorf("auth failed: %w", ErrUnauthorized), true},
		{"ErrForbidden", ErrForbidden, false},
		{"ErrNotFound", ErrNotFound, false},
		{"generic error", errors.New("some error"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.err); got != tt.expected {
				t.Errorf("IsUnauthorized(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestIsAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"APIError", &APIError{Code: -32600, Message: "Invalid"}, true},
		{"wrapped APIError", fmt.Errorf("call failed: %w", &APIError{Code: -32600, Message: "Invalid"}), true},
		{"ErrUnauthorized", ErrUnauthorized, false},
		{"ErrNotFound", ErrNotFound, false},
		{"generic error", errors.New("some error"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAPIError(tt.err); got != tt.expected {
				t.Errorf("IsAPIError(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestErrorsIs(t *testing.T) {
	// Test that errors.Is works correctly with sentinel errors
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{"direct match", ErrTaskNotFound, ErrTaskNotFound, true},
		{"wrapped match", fmt.Errorf("ctx: %w", ErrTaskNotFound), ErrTaskNotFound, true},
		{"double wrapped", fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", ErrTaskNotFound)), ErrTaskNotFound, true},
		{"different error", ErrTaskNotFound, ErrProjectNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.expected {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, got, tt.expected)
			}
		})
	}
}

func TestErrorsAs(t *testing.T) {
	// Test that errors.As works correctly with APIError
	apiErr := &APIError{Code: -32600, Message: "Invalid Request"}
	wrappedErr := fmt.Errorf("call failed: %w", apiErr)

	var target *APIError

	// Direct APIError
	if !errors.As(apiErr, &target) {
		t.Error("errors.As should match direct APIError")
	}
	if target.Code != -32600 {
		t.Errorf("expected Code=-32600, got %d", target.Code)
	}

	// Wrapped APIError
	target = nil
	if !errors.As(wrappedErr, &target) {
		t.Error("errors.As should match wrapped APIError")
	}
	if target.Message != "Invalid Request" {
		t.Errorf("expected Message='Invalid Request', got %s", target.Message)
	}

	// Non-APIError
	target = nil
	if errors.As(ErrNotFound, &target) {
		t.Error("errors.As should not match non-APIError")
	}
}

func TestSentinelErrorMessages(t *testing.T) {
	// Ensure all sentinel errors have meaningful messages
	sentinels := []struct {
		name string
		err  error
	}{
		{"ErrConnectionFailed", ErrConnectionFailed},
		{"ErrTimeout", ErrTimeout},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrForbidden", ErrForbidden},
		{"ErrNotFound", ErrNotFound},
		{"ErrProjectNotFound", ErrProjectNotFound},
		{"ErrTaskNotFound", ErrTaskNotFound},
		{"ErrColumnNotFound", ErrColumnNotFound},
		{"ErrCommentNotFound", ErrCommentNotFound},
		{"ErrAlreadyInLastColumn", ErrAlreadyInLastColumn},
		{"ErrAlreadyInFirstColumn", ErrAlreadyInFirstColumn},
		{"ErrTaskClosed", ErrTaskClosed},
		{"ErrTaskOpen", ErrTaskOpen},
		{"ErrEmptyTitle", ErrEmptyTitle},
		{"ErrInvalidProjectID", ErrInvalidProjectID},
	}

	for _, s := range sentinels {
		t.Run(s.name, func(t *testing.T) {
			if s.err == nil {
				t.Errorf("%s should not be nil", s.name)
			}
			if s.err.Error() == "" {
				t.Errorf("%s should have a non-empty error message", s.name)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test error wrapping preserves context
	originalErr := ErrTaskNotFound
	wrappedOnce := fmt.Errorf("getting task %d: %w", 42, originalErr)
	wrappedTwice := fmt.Errorf("in board scope: %w", wrappedOnce)

	// Should preserve original error
	if !errors.Is(wrappedTwice, ErrTaskNotFound) {
		t.Error("wrapped error should match original with errors.Is")
	}

	// Should include context in message
	if wrappedTwice.Error() != "in board scope: getting task 42: task not found" {
		t.Errorf("unexpected error message: %s", wrappedTwice.Error())
	}
}

func TestOperationFailedError(t *testing.T) {
	tests := []struct {
		name           string
		err            *OperationFailedError
		expectedSubstr []string
	}{
		{
			name: "with hints",
			err: &OperationFailedError{
				Operation: "moveTaskPosition(task=42, column=5, project=1)",
				Hints:     []string{"task may not exist", "column may not belong to project"},
			},
			expectedSubstr: []string{
				"moveTaskPosition",
				"operation failed",
				"possible causes",
				"task may not exist",
				"column may not belong to project",
			},
		},
		{
			name: "without hints",
			err: &OperationFailedError{
				Operation: "someOperation",
				Hints:     nil,
			},
			expectedSubstr: []string{"someOperation", "operation failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, substr := range tt.expectedSubstr {
				if !containsSubstr(errMsg, substr) {
					t.Errorf("error message %q should contain %q", errMsg, substr)
				}
			}
		})
	}
}

func TestIsOperationFailed(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"OperationFailedError", &OperationFailedError{Operation: "test"}, true},
		{"wrapped OperationFailedError", fmt.Errorf("call failed: %w", &OperationFailedError{Operation: "test"}), true},
		{"ErrUnauthorized", ErrUnauthorized, false},
		{"generic error", errors.New("some error"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOperationFailed(tt.err); got != tt.expected {
				t.Errorf("IsOperationFailed(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func containsSubstr(s, substr string) bool {
	return strings.Contains(s, substr)
}
