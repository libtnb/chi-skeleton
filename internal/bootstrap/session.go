package bootstrap

import (
	"log/slog"

	"github.com/libtnb/gormstore"
	"github.com/libtnb/sessions"
	"github.com/samber/do/v2"

	"github.com/libtnb/chi-skeleton/internal/config"
	"github.com/libtnb/chi-skeleton/internal/data"
)

func NewSession(i do.Injector) (*sessions.Manager, error) {
	conf := do.MustInvoke[*config.Config](i)

	manager, err := sessions.NewManager(&sessions.ManagerOptions{
		Key:                  conf.App.Key,
		Lifetime:             conf.Session.Lifetime,
		GcInterval:           conf.Session.GcInterval,
		DisableDefaultDriver: true,
		// background errors (GC, middleware saves) land in the app log
		Logger: do.MustInvoke[*slog.Logger](i),
	})
	if err != nil {
		return nil, err
	}

	// extend gorm store driver
	store := gormstore.New(do.MustInvoke[*data.Data](i).DB)
	if err = manager.Extend("default", store); err != nil {
		return nil, err
	}

	return manager, nil
}
