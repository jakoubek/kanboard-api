package kanboard

import (
	"context"
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
