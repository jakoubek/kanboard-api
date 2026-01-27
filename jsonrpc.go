package kanboard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      int64       `json:"id"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error.
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("JSON-RPC error (code %d): %s", e.Code, e.Message)
}

// nextRequestID returns a random request ID.
func nextRequestID() int64 {
	return rand.Int64()
}

// call sends a JSON-RPC request and parses the response.
// The result parameter should be a pointer to the expected result type.
func (c *Client) call(ctx context.Context, method string, params interface{}, result interface{}) error {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		ID:      nextRequestID(),
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.logger != nil {
		c.logger.Debug("JSON-RPC request", "method", method, "body", string(body))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.auth != nil {
		c.auth.Apply(httpReq)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// Preserve specific errors that shouldn't be wrapped as connection failures
		if errors.Is(err, ErrTooManyRedirects) {
			return ErrTooManyRedirects
		}
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if resp.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	// Check for HTML response (indicates redirect to login page or misconfiguration)
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return fmt.Errorf("%w: received HTML instead of JSON (possible redirect to login page)", ErrUnauthorized)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.logger != nil {
		c.logger.Debug("JSON-RPC response", "method", method, "body", string(respBody))
	}

	// Check if response body is HTML (backup check if Content-Type header is wrong/missing)
	if len(respBody) > 0 {
		trimmed := bytes.TrimLeft(respBody, " \t\n\r")
		if len(trimmed) > 0 && trimmed[0] == '<' {
			// Response is HTML, not JSON - extract a preview for debugging
			preview := string(respBody)
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			return fmt.Errorf("server returned HTML instead of JSON (possible auth error or server error): %s", preview)
		}
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return &APIError{
			Code:    rpcResp.Error.Code,
			Message: rpcResp.Error.Message,
		}
	}

	if result != nil && rpcResp.Result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return nil
}
