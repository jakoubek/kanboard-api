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
	baseURL        string
	endpoint       string
	httpClient     *http.Client
	auth           Authenticator
	logger         *slog.Logger
	authHeaderName string // custom auth header, empty = use "Authorization"
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

// WithAuthHeader configures a custom header name for authentication.
// If not set, the standard "Authorization" header is used.
// This must be called before configuring authentication (WithAPIToken, WithBasicAuth, etc.).
// Example: WithAuthHeader("X-API-Auth")
func (c *Client) WithAuthHeader(headerName string) *Client {
	c.authHeaderName = headerName
	return c
}

// WithAPIToken configures the client to use API token authentication.
// Uses "jsonrpc" as the username for HTTP Basic Auth.
func (c *Client) WithAPIToken(token string) *Client {
	c.auth = &apiTokenAuth{token: token, headerName: c.authHeaderName}
	return c
}

// WithAPITokenUser configures the client to use API token authentication with a custom username.
// If user is empty, "jsonrpc" will be used as the default.
func (c *Client) WithAPITokenUser(token, user string) *Client {
	c.auth = &apiTokenAuth{token: token, user: user, headerName: c.authHeaderName}
	return c
}

// WithBasicAuth configures the client to use username/password authentication.
func (c *Client) WithBasicAuth(username, password string) *Client {
	c.auth = &basicAuth{username: username, password: password, headerName: c.authHeaderName}
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
