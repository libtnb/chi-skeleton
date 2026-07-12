package server

import (
	"context"
	"net/http"

	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/pkg/transport"
)

// healthChecker is the one container capability Readyz needs; keeping the
// dependency this narrow makes the service easy to fake in tests.
type healthChecker interface {
	HealthCheckWithContext(ctx context.Context) map[string]error
}

// HealthService serves the container/orchestrator probes.
type HealthService struct {
	checker healthChecker
}

func NewHealthService(i do.Injector) (*HealthService, error) {
	return &HealthService{
		checker: i,
	}, nil
}

// Healthz is the liveness probe: the process is up and serving.
func (r *HealthService) Healthz(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

// Readyz is the readiness probe: every health-checkable service in the
// container (the database, and whatever you add later) must pass.
func (r *HealthService) Readyz(w http.ResponseWriter, req *http.Request) {
	for name, err := range r.checker.HealthCheckWithContext(req.Context()) {
		if err != nil {
			transport.Error(w, http.StatusServiceUnavailable, "%s unavailable", name)
			return
		}
	}

	_, _ = w.Write([]byte("ok"))
}

// HealthRoutes endpoints are undocumented, so they carry no Request/Response samples.
func HealthRoutes(i do.Injector) (Endpoints, error) {
	health := do.MustInvoke[*HealthService](i)

	return Endpoints{
		{Method: http.MethodGet, Path: "/", Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello, World 👋!"))
		}},
		{Method: http.MethodGet, Path: "/healthz", Handler: health.Healthz},
		{Method: http.MethodGet, Path: "/readyz", Handler: health.Readyz},
	}, nil
}
