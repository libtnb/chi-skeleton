package route

import (
	"net/http"
	"runtime/debug"

	"github.com/dromara/carbon/v2"
	"github.com/go-chi/chi/v5"
	"github.com/libtnb/validator/contrib/openapi"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/registry"
)

// Prefix marks route contributions, collected by the assemblers below and
// registered sorted by service name.
const RoutePrefix = "routes:"

// Package lists every route contribution; add yours here.
var Package = do.Package(
	do.LazyNamed(RoutePrefix+"health", HealthRoutes),
	do.LazyNamed(RoutePrefix+"user", UserRoutes),
	do.LazyNamed(RoutePrefix+"ws", WsRoutes),
)

// Endpoint declares one HTTP endpoint: how to serve it and, through the
// Request/Response samples, how to document it. Endpoints without either
// stay out of the OpenAPI document (probes, websockets).
type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Summary string
	Tags    []string
	// Request documents parameters and body: uri tags become path
	// parameters, query tags query parameters, json tags the body —
	// constraints are read from the validate tags.
	Request any
	// Response documents the response body; Status defaults to 200.
	Response any
	Status   int
}

// Endpoints is a module's contribution to the HTTP router.
type Endpoints []Endpoint

// HTTP registers every "routes:*" contribution on r.
func HTTP(i do.Injector, r chi.Router) error {
	groups, err := registry.Collect[Endpoints](i, RoutePrefix)
	if err != nil {
		return err
	}
	for _, endpoints := range groups {
		for _, e := range endpoints {
			r.Method(e.Method, e.Path, e.Handler)
		}
	}

	return nil
}

// SpecJSON assembles the OpenAPI 3.1 document from every documented endpoint.
func SpecJSON(i do.Injector, title string) ([]byte, error) {
	g := openapi.New(title, buildVersion(),
		openapi.WithType(carbon.DateTime{}, &openapi.Schema{Type: "string", Format: "date-time"}),
	)

	groups, err := registry.Collect[Endpoints](i, RoutePrefix)
	if err != nil {
		return nil, err
	}
	for _, endpoints := range groups {
		for _, e := range endpoints {
			if e.Request == nil && e.Response == nil {
				continue
			}
			err := g.Add(e.Method, e.Path, openapi.Op{
				Summary:  e.Summary,
				Tags:     e.Tags,
				Request:  e.Request,
				Response: e.Response,
				Status:   e.Status,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return g.JSON()
}

func buildVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}
