package kanboard

import "net/http"

// Authenticator applies authentication to HTTP requests.
type Authenticator interface {
	Apply(req *http.Request)
}

// apiTokenAuth implements API token authentication.
type apiTokenAuth struct {
	token string
}

// Apply adds HTTP Basic Auth with username "jsonrpc" and the API token.
func (a *apiTokenAuth) Apply(req *http.Request) {
	req.SetBasicAuth("jsonrpc", a.token)
}

// basicAuth implements username/password authentication.
type basicAuth struct {
	username string
	password string
}

// Apply adds HTTP Basic Auth with username and password.
func (a *basicAuth) Apply(req *http.Request) {
	req.SetBasicAuth(a.username, a.password)
}
