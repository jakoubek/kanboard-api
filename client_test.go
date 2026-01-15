package kanboard

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://kanboard.example.com")

	if client.baseURL != "https://kanboard.example.com" {
		t.Errorf("expected baseURL='https://kanboard.example.com', got %s", client.baseURL)
	}
	if client.endpoint != "https://kanboard.example.com/jsonrpc.php" {
		t.Errorf("expected endpoint='https://kanboard.example.com/jsonrpc.php', got %s", client.endpoint)
	}
	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("expected timeout=%v, got %v", DefaultTimeout, client.httpClient.Timeout)
	}
}

func TestNewClient_TrailingSlash(t *testing.T) {
	client := NewClient("https://kanboard.example.com/")

	if client.baseURL != "https://kanboard.example.com" {
		t.Errorf("trailing slash should be removed, got %s", client.baseURL)
	}
	if client.endpoint != "https://kanboard.example.com/jsonrpc.php" {
		t.Errorf("expected endpoint='https://kanboard.example.com/jsonrpc.php', got %s", client.endpoint)
	}
}

func TestNewClient_Subdirectory(t *testing.T) {
	client := NewClient("https://example.com/kanboard")

	if client.endpoint != "https://example.com/kanboard/jsonrpc.php" {
		t.Errorf("expected endpoint='https://example.com/kanboard/jsonrpc.php', got %s", client.endpoint)
	}
}

func TestDefaultTimeout(t *testing.T) {
	if DefaultTimeout != 30*time.Second {
		t.Errorf("expected DefaultTimeout=30s, got %v", DefaultTimeout)
	}
}

func TestClient_WithAPIToken(t *testing.T) {
	client := NewClient("https://example.com")
	result := client.WithAPIToken("my-token")

	// Should return same client instance
	if client != result {
		t.Error("WithAPIToken should return the same client instance")
	}

	if client.auth == nil {
		t.Error("auth should be set")
	}
}

func TestClient_WithBasicAuth(t *testing.T) {
	client := NewClient("https://example.com")
	result := client.WithBasicAuth("admin", "password")

	// Should return same client instance
	if client != result {
		t.Error("WithBasicAuth should return the same client instance")
	}

	if client.auth == nil {
		t.Error("auth should be set")
	}
}

func TestClient_WithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	client := NewClient("https://example.com")
	result := client.WithHTTPClient(customClient)

	// Should return same client instance
	if client != result {
		t.Error("WithHTTPClient should return the same client instance")
	}

	if client.httpClient != customClient {
		t.Error("httpClient should be set to custom client")
	}
}

func TestClient_WithTimeout(t *testing.T) {
	client := NewClient("https://example.com")
	originalClient := client.httpClient

	result := client.WithTimeout(60 * time.Second)

	// Should return same client instance
	if client != result {
		t.Error("WithTimeout should return the same client instance")
	}

	// Should create new http.Client
	if client.httpClient == originalClient {
		t.Error("WithTimeout should create a new http.Client")
	}

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("expected timeout=60s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_WithLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := NewClient("https://example.com")
	result := client.WithLogger(logger)

	// Should return same client instance
	if client != result {
		t.Error("WithLogger should return the same client instance")
	}

	if client.logger != logger {
		t.Error("logger should be set")
	}
}

func TestClient_FluentChaining(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := NewClient("https://example.com").
		WithAPIToken("my-token").
		WithTimeout(60 * time.Second).
		WithLogger(logger)

	if client.auth == nil {
		t.Error("auth should be set via chaining")
	}
	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("timeout should be 60s, got %v", client.httpClient.Timeout)
	}
	if client.logger != logger {
		t.Error("logger should be set via chaining")
	}
}

func TestClient_DefaultsWithoutConfiguration(t *testing.T) {
	client := NewClient("https://example.com")

	// Should have defaults
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("default timeout should be %v", DefaultTimeout)
	}

	// Should have nil for optional fields
	if client.auth != nil {
		t.Error("auth should be nil by default")
	}
	if client.logger != nil {
		t.Error("logger should be nil by default")
	}
}
