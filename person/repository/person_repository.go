package repository

import (
	"database/sql"
	"fmt"
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

func (p *postgresPersonRepo) GetByID(c *gin.Context, id int64) (domain.Person, error) {
	var person domain.Person
	err := p.DB.QueryRow(
		"SELECT id, name, created_at FROM person WHERE id = $1", id).Scan(
		&person.Id, &person.Name, &person.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return person, fmt.Errorf("No data found with ID %d", id)
		}
		return person, err
	}
	return person, nil

}

func (p *postgresPersonRepo) CreatePerson(c *gin.Context, person domain.Person) error {
	query := `INSERT INTO person (name) VALUES ($1) `
	_, err := p.DB.Exec(query, person.Name)
	if err != nil {
		return err
	}

	return nil
}
