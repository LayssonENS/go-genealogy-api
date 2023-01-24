package database

import (
	"database/sql"
	"fmt"
	"github.com/LayssonENS/go-genealogy-api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

func NewPostgresConnection(dbConfig config.DbConfig) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open postgres connection")
	}

	return db, nil
}

func DBMigrate(dbInstance *sql.DB, dbConfig config.DbConfig) error {

	driver, err := postgres.WithInstance(dbInstance, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to instantiate postgres driver")
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
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
