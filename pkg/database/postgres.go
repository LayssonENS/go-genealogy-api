package database

import (
	"database/sql"
	"fmt"
	"github.com/LayssonENS/go-genealogy-api/pkg/config"
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
