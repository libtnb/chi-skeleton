package server

import (
	"net/http"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
)

// HealthRoutes serves the probes; they stay out of the OpenAPI docs.
func HealthRoutes(i do.Injector) (transport.Endpoints, error) {
	return transport.Endpoints{
		{Method: http.MethodGet, Path: "/", Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello, World 👋!"))
		}},
		{Method: http.MethodGet, Path: "/healthz", Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		}},
		{Method: http.MethodGet, Path: "/readyz", Handler: func(w http.ResponseWriter, r *http.Request) {
			for name, err := range i.HealthCheckWithContext(r.Context()) {
				if err != nil {
					transport.Error(w, http.StatusServiceUnavailable, "%s unavailable", name)
					return
				}
			}
			_, _ = w.Write([]byte("ok"))
		}},
	}, nil
}
