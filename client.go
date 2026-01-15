package kanboard

import (
	"net/http"
	"strings"
)

// Client is the Kanboard API client.
type Client struct {
	baseURL    string
	endpoint   string
	httpClient *http.Client
	auth       Authenticator
}

// NewClient creates a new Kanboard API client.
// The baseURL should be the base URL of the Kanboard instance (e.g., "https://kanboard.example.com").
// The path /jsonrpc.php is appended automatically.
// Supports subdirectory installations (e.g., "https://example.com/kanboard" → POST https://example.com/kanboard/jsonrpc.php).
func NewClient(baseURL string) *Client {
	// Ensure no trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL:    baseURL,
		endpoint:   baseURL + "/jsonrpc.php",
		httpClient: http.DefaultClient,
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
func (c *Client) WithHTTPClient(client *http.Client) *Client {
	c.httpClient = client
	return c
}
