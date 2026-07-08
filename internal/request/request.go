package request

import (
	"net/http"
)

// WithPrepare runs after binding and before validation: fill defaults or
// normalize values. Authorization belongs in middleware or usecases.
type WithPrepare interface {
	Prepare(r *http.Request) error
}

// WithRules ANDs extra rules onto the struct tags at runtime.
type WithRules interface {
	Rules(r *http.Request) map[string]string
}

// WithFilters applies value filters (trim, lower, ...) before validation.
type WithFilters interface {
	Filters(r *http.Request) map[string]string
}

// WithMessages overrides message templates for this request only.
type WithMessages interface {
	Messages(r *http.Request) map[string]string
}
