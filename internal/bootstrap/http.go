package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-rat/chi-skeleton/internal/app"
	"github.com/go-rat/chi-skeleton/internal/http/middleware"
	"github.com/go-rat/chi-skeleton/internal/route"
)

func initHttp() {
	app.Http = chi.NewRouter()

	// add middleware
	app.Http.Use(middleware.GlobalMiddleware()...)

	// add route
	route.Http(app.Http)

	if err := http.ListenAndServe(app.Conf.MustString("http.address"), app.Http); err != nil {
		panic(fmt.Sprintf("failed to start http server: %v", err))
	}
}
