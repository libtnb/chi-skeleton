package request

import (
	"net/http"
)

type HasAuthorize interface {
	Authorize(r *http.Request) error
}

type HasPrepare interface {
	Prepare(r *http.Request) error
}

type HasRules interface {
	Rules(r *http.Request) map[string]string
}

type HasMessages interface {
	Messages(r *http.Request) map[string]string
}
