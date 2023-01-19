package repository

import (
	"context"
	"database/sql"
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
)

type postgresPersonRepo struct {
	DB *sql.DB
}

// NewPostgresPersonRepository will create an implementation of author.Repository
func NewPostgresPersonRepository(db *sql.DB) domain.PersonRepository {
	return &postgresPersonRepo{
		DB: db,
	}
}

func (p *postgresPersonRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Person, err error) {
	stmt, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Person{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Person{}

	err = row.Scan(
		&res.Name,
	)
	return
}

func (p *postgresPersonRepo) GetByID(c *gin.Context, id int64) (domain.Person, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return p.getOne(c, query, id)
}
