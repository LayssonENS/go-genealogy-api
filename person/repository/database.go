package repository

import (
	"database/sql"
	"github.com/LayssonENS/go-genealogy-api/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

func DBMigrate(dbInstance *sql.DB, dbConfig config.DbConfig) error {

	driver, err := postgres.WithInstance(dbInstance, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to instantiate postgres driver")
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://person/repository/migrations",
		dbConfig.Name, driver)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrate up")
	}
	return nil
}
