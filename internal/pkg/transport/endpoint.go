package transport

import (
	"net/http"
)

// Endpoint declares one HTTP endpoint; without Request/Response samples it
// stays out of the OpenAPI document.
type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Summary string
	Tags    []string
	// Request documents parameters and body via uri/query/json tags;
	// constraints come from the validate tags.
	Request any
	// Response documents the response body; Status defaults to 200.
	Response any
	Status   int
}

// Endpoints is a module's route contribution, registered under
// registry.RoutePrefix.
type Endpoints []Endpoint
