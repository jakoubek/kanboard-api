package kanboard

import (
	"encoding/base64"
	"net/http"
)

// Authenticator applies authentication to HTTP requests.
type Authenticator interface {
	Apply(req *http.Request)
}

// apiTokenAuth implements API token authentication.
type apiTokenAuth struct {
	user       string
	token      string
	headerName string
}

// Apply adds HTTP Basic Auth with the configured user (or "jsonrpc" if empty) and the API token.
func (a *apiTokenAuth) Apply(req *http.Request) {
	user := a.user
	if user == "" {
		user = "jsonrpc"
	}
	if a.headerName != "" {
		req.Header.Set(a.headerName, "Basic "+basicAuthValue(user, a.token))
	} else {
		req.SetBasicAuth(user, a.token)
	}
}

// basicAuth implements username/password authentication.
type basicAuth struct {
	username   string
	password   string
	headerName string
}

// Apply adds HTTP Basic Auth with username and password.
func (a *basicAuth) Apply(req *http.Request) {
	if a.headerName != "" {
		req.Header.Set(a.headerName, "Basic "+basicAuthValue(a.username, a.password))
	} else {
		req.SetBasicAuth(a.username, a.password)
	}
}

// basicAuthValue returns the base64-encoded value for HTTP Basic Auth.
func basicAuthValue(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
