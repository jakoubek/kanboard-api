package kanboard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPITokenAuth_Apply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "jsonrpc" {
			t.Errorf("expected username=jsonrpc, got %s", username)
		}
		if password != "my-api-token-12345" {
			t.Errorf("expected password=my-api-token-12345, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("my-api-token-12345")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPITokenAuth_CustomUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "custom-user" {
			t.Errorf("expected username=custom-user, got %s", username)
		}
		if password != "my-api-token-12345" {
			t.Errorf("expected password=my-api-token-12345, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPITokenUser("my-api-token-12345", "custom-user")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPITokenAuth_EmptyUserDefaultsToJsonrpc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "jsonrpc" {
			t.Errorf("expected username=jsonrpc (default), got %s", username)
		}
		if password != "my-api-token-12345" {
			t.Errorf("expected password=my-api-token-12345, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Using WithAPITokenUser with empty user should default to "jsonrpc"
	client := NewClient(server.URL).WithAPITokenUser("my-api-token-12345", "")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBasicAuth_Apply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "admin" {
			t.Errorf("expected username=admin, got %s", username)
		}
		if password != "secret-password" {
			t.Errorf("expected password=secret-password, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithBasicAuth("admin", "secret-password")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNoAuth_NoHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, ok := r.BasicAuth()
		if ok {
			t.Error("did not expect Basic Auth header when no auth configured")
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Client without any auth configured
	client := NewClient(server.URL)

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthenticator_Interface(t *testing.T) {
	// Verify both types implement Authenticator interface
	var _ Authenticator = &apiTokenAuth{}
	var _ Authenticator = &basicAuth{}
}

func TestAPITokenAuth_CustomHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should NOT have standard Authorization header via BasicAuth
		_, _, ok := r.BasicAuth()
		if ok {
			t.Error("did not expect Authorization header to be parsed as Basic Auth")
		}

		// Should have custom header
		customAuth := r.Header.Get("X-API-Auth")
		if customAuth == "" {
			t.Error("expected X-API-Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify the custom header has the correct Basic auth value
		expected := "Basic " + base64Encode("jsonrpc:my-api-token")
		if customAuth != expected {
			t.Errorf("expected X-API-Auth=%s, got %s", expected, customAuth)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).
		WithAuthHeader("X-API-Auth").
		WithAPIToken("my-api-token")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBasicAuth_CustomHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should NOT have standard Authorization header via BasicAuth
		_, _, ok := r.BasicAuth()
		if ok {
			t.Error("did not expect Authorization header to be parsed as Basic Auth")
		}

		// Should have custom header
		customAuth := r.Header.Get("X-Custom-Auth")
		if customAuth == "" {
			t.Error("expected X-Custom-Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify the custom header has the correct Basic auth value
		expected := "Basic " + base64Encode("admin:secret")
		if customAuth != expected {
			t.Errorf("expected X-Custom-Auth=%s, got %s", expected, customAuth)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).
		WithAuthHeader("X-Custom-Auth").
		WithBasicAuth("admin", "secret")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCustomHeader_WithCustomUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customAuth := r.Header.Get("X-API-Auth")
		if customAuth == "" {
			t.Error("expected X-API-Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify the custom header uses the custom username
		expected := "Basic " + base64Encode("custom-user:my-token")
		if customAuth != expected {
			t.Errorf("expected X-API-Auth=%s, got %s", expected, customAuth)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).
		WithAuthHeader("X-API-Auth").
		WithAPITokenUser("my-token", "custom-user")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// base64Encode is a test helper to generate expected auth values.
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestClient_FluentAuthConfiguration(t *testing.T) {
	// Test that fluent methods return the same client instance
	client := NewClient("https://example.com")

	client2 := client.WithAPIToken("token")
	if client != client2 {
		t.Error("WithAPIToken should return the same client instance")
	}

	client3 := client.WithBasicAuth("user", "pass")
	if client != client3 {
		t.Error("WithBasicAuth should return the same client instance")
	}

	client4 := client.WithAPITokenUser("token", "custom-user")
	if client != client4 {
		t.Error("WithAPITokenUser should return the same client instance")
	}
}

func TestAuthOverwrite(t *testing.T) {
	// Test that setting auth multiple times overwrites previous auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Should use the last configured auth (BasicAuth)
		if username != "final-user" {
			t.Errorf("expected username=final-user, got %s", username)
		}
		if password != "final-pass" {
			t.Errorf("expected password=final-pass, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Configure API token first, then overwrite with BasicAuth
	client := NewClient(server.URL).
		WithAPIToken("initial-token").
		WithBasicAuth("final-user", "final-pass")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
