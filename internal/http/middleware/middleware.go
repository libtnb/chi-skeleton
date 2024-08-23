package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.URLFormat,
		middleware.Logger,
		middleware.Recoverer,
	}
}
