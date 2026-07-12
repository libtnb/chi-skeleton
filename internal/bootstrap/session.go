package bootstrap

import (
	"log/slog"

	"github.com/go-rio/rio"
	"github.com/libtnb/sessions"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/conf"
)

func NewSession(i do.Injector) (*sessions.Manager, error) {
	config := do.MustInvoke[*conf.Config](i)

	manager, err := sessions.NewManager(&sessions.ManagerOptions{
		Key:                  config.App.Key,
		Lifetime:             config.Session.Lifetime,
		GcInterval:           config.Session.GcInterval,
		DisableDefaultDriver: true,
		// background errors (GC, middleware saves) land in the app log
		Logger: do.MustInvoke[*slog.Logger](i),
	})
	if err != nil {
		return nil, err
	}

	store, err := newSessionStore(do.MustInvoke[*rio.DB](i))
	if err != nil {
		return nil, err
	}
	if err = manager.Extend("default", store); err != nil {
		return nil, err
	}

	return manager, nil
}
