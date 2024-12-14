package bootstrap

import (
	"log/slog"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/go-rat/chi-skeleton/internal/migration"
)

func NewDB(conf *koanf.Koanf, log *slog.Logger) (*gorm.DB, error) {
	// You can use any other database, like MySQL or PostgreSQL.
	db, err := gorm.Open(sqlite.Open(conf.MustString("database.path")), &gorm.Config{
		Logger:                                   slogGorm.New(slogGorm.WithHandler(log.Handler())),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMigrate(db *gorm.DB) error {
	migrator := gormigrate.New(db, &gormigrate.Options{
		UseTransaction: true, // Note: MySQL not support DDL transaction
	}, migration.Migrations)

	return migrator.Migrate()
}
