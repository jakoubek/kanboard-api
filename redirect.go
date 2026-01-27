package kanboard

import (
	"net/http"
	"net/url"
	"strings"
)

// maxRedirects is the maximum number of redirects to follow (Go's default).
const maxRedirects = 10

// redirectBehavior is a CheckRedirect handler that preserves authentication
// headers for same-host redirects. Go's http.Client strips the Authorization
// header on redirects by default (security feature since Go 1.8).
func (c *Client) redirectBehavior(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return ErrTooManyRedirects
	}

	if len(via) == 0 {
		return nil
	}

	// Check if we're redirecting to the same host
	originalReq := via[0]
	if isSameHost(originalReq.URL, req.URL) {
		// Preserve auth header for same-host redirects
		headerName := "Authorization"
		if c.authHeaderName != "" {
			headerName = c.authHeaderName
		}

		if authValue := originalReq.Header.Get(headerName); authValue != "" {
			req.Header.Set(headerName, authValue)
		}
	}

	return nil
}

// isSameHost compares two URLs to determine if they have the same host.
// It normalizes default ports (80 for HTTP, 443 for HTTPS).
func isSameHost(a, b *url.URL) bool {
	return normalizeHost(a) == normalizeHost(b)
}

// normalizeHost returns the host with default ports removed.
// http://example.com:80 -> example.com
// https://example.com:443 -> example.com
// http://example.com:8080 -> example.com:8080
func normalizeHost(u *url.URL) string {
	host := strings.ToLower(u.Host)

	// Remove default ports
	switch u.Scheme {
	case "http":
		host = strings.TrimSuffix(host, ":80")
	case "https":
		host = strings.TrimSuffix(host, ":443")
	}

	return host
}
