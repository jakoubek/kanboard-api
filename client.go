package kanboard

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// DefaultTimeout is the default HTTP timeout for API requests.
const DefaultTimeout = 30 * time.Second

// Client is the Kanboard API client.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
	baseURL    string
	endpoint   string
	httpClient *http.Client
	auth       Authenticator
	logger     *slog.Logger
}

// NewClient creates a new Kanboard API client.
// The baseURL should be the base URL of the Kanboard instance (e.g., "https://kanboard.example.com").
// The path /jsonrpc.php is appended automatically.
// Supports subdirectory installations (e.g., "https://example.com/kanboard" → POST https://example.com/kanboard/jsonrpc.php).
func NewClient(baseURL string) *Client {
	// Ensure no trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL:  baseURL,
		endpoint: baseURL + "/jsonrpc.php",
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// WithAPIToken configures the client to use API token authentication.
func (c *Client) WithAPIToken(token string) *Client {
	c.auth = &apiTokenAuth{token: token}
	return c
}

// WithBasicAuth configures the client to use username/password authentication.
func (c *Client) WithBasicAuth(username, password string) *Client {
	c.auth = &basicAuth{username: username, password: password}
	return c
}

// WithHTTPClient sets a custom HTTP client.
// This replaces the default client entirely, including any timeout settings.
func (c *Client) WithHTTPClient(client *http.Client) *Client {
	c.httpClient = client
	return c
}

// WithTimeout sets the HTTP client timeout.
// This creates a new HTTP client with the specified timeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient = &http.Client{
		Timeout:   timeout,
		Transport: c.httpClient.Transport,
	}
	return c
}

// WithLogger sets the logger for debug output.
// If set, the client will log request/response details at debug level.
func (c *Client) WithLogger(logger *slog.Logger) *Client {
	c.logger = logger
	return c
}
