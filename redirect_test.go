package kanboard

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIsSameHost(t *testing.T) {
	tests := []struct {
		name     string
		urlA     string
		urlB     string
		expected bool
	}{
		{
			name:     "same host and scheme",
			urlA:     "https://example.com/path",
			urlB:     "https://example.com/other",
			expected: true,
		},
		{
			name:     "same host different scheme",
			urlA:     "http://example.com/path",
			urlB:     "https://example.com/path",
			expected: true,
		},
		{
			name:     "http with explicit port 80",
			urlA:     "http://example.com:80/path",
			urlB:     "http://example.com/path",
			expected: true,
		},
		{
			name:     "https with explicit port 443",
			urlA:     "https://example.com:443/path",
			urlB:     "https://example.com/path",
			expected: true,
		},
		{
			name:     "different hosts",
			urlA:     "https://example.com/path",
			urlB:     "https://other.com/path",
			expected: false,
		},
		{
			name:     "different non-standard ports",
			urlA:     "https://example.com:8080/path",
			urlB:     "https://example.com:9090/path",
			expected: false,
		},
		{
			name:     "non-standard port vs no port",
			urlA:     "https://example.com:8080/path",
			urlB:     "https://example.com/path",
			expected: false,
		},
		{
			name:     "case insensitive host",
			urlA:     "https://Example.COM/path",
			urlB:     "https://example.com/path",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := url.Parse(tt.urlA)
			b, _ := url.Parse(tt.urlB)
			result := isSameHost(a, b)
			if result != tt.expected {
				t.Errorf("isSameHost(%q, %q) = %v, want %v", tt.urlA, tt.urlB, result, tt.expected)
			}
		})
	}
}

func TestRedirectPreservesAuthForSameHost(t *testing.T) {
	redirectCount := 0
	var receivedAuth string

	// Server that redirects once, then returns success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")

		if redirectCount == 0 {
			redirectCount++
			// Redirect to same host with different path
			http.Redirect(w, r, "/redirected", http.StatusFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"test"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	// Override endpoint to point to test server
	client.endpoint = server.URL + "/jsonrpc.php"

	var result string
	err := client.call(t.Context(), "test", nil, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuth == "" {
		t.Error("Authorization header was not preserved after redirect")
	}
}

func TestRedirectPreservesCustomAuthHeader(t *testing.T) {
	redirectCount := 0
	var receivedAuth string
	customHeader := "X-Custom-Auth"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get(customHeader)

		if redirectCount == 0 {
			redirectCount++
			http.Redirect(w, r, "/redirected", http.StatusFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"test"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL).
		WithAuthHeader(customHeader).
		WithAPIToken("test-token")

	client.endpoint = server.URL + "/jsonrpc.php"

	var result string
	err := client.call(t.Context(), "test", nil, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuth == "" {
		t.Errorf("Custom auth header %q was not preserved after redirect", customHeader)
	}
}

// Note: Cross-domain redirect behavior is handled by Go's http.Client.
// Go preserves Authorization headers for same-domain redirects (including localhost:port1 to localhost:port2).
// Our custom redirect handler adds value for custom auth headers (e.g., "X-Custom-Auth")
// which Go's default behavior doesn't handle.

func TestRedirectLimit(t *testing.T) {
	redirectCount := 0

	// Server that always redirects
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectCount++
		http.Redirect(w, r, "/redirect"+string(rune('0'+redirectCount)), http.StatusFound)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")
	client.endpoint = server.URL + "/jsonrpc.php"

	var result string
	err := client.call(t.Context(), "test", nil, &result)

	if err == nil {
		t.Fatal("expected error for too many redirects")
	}

	if !errors.Is(err, ErrTooManyRedirects) {
		t.Errorf("expected ErrTooManyRedirects, got: %v", err)
	}

	if redirectCount > maxRedirects+1 {
		t.Errorf("followed %d redirects, expected max %d", redirectCount, maxRedirects)
	}
}

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
	}{
		{"http://example.com", "example.com"},
		{"http://example.com:80", "example.com"},
		{"http://example.com:8080", "example.com:8080"},
		{"https://example.com", "example.com"},
		{"https://example.com:443", "example.com"},
		{"https://example.com:8443", "example.com:8443"},
		{"https://Example.COM:443", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			u, _ := url.Parse(tt.rawURL)
			result := normalizeHost(u)
			if result != tt.expected {
				t.Errorf("normalizeHost(%q) = %q, want %q", tt.rawURL, result, tt.expected)
			}
		})
	}
}
