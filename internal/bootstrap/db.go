package bootstrap

import (
	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-rat/chi-skeleton/internal/migration"
)

func NewDB(conf *koanf.Koanf) (*gorm.DB, error) {
	logLevel := logger.Error
	if conf.Bool("database.debug") {
		logLevel = logger.Info
	}
	// You can use any other database, like MySQL or PostgreSQL.
	db, err := gorm.Open(sqlite.Open(conf.MustString("database.path")), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
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
	if err := migrator.Migrate(); err != nil {
		return err
	}

	return nil
}
