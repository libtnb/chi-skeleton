package request

import (
	"net/http"
)

type Request[T any] interface {
	*T
	PrepareForValidation(r *http.Request) error
}
