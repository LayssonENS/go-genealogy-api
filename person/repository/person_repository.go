package personRepository

import (
	"database/sql"
	"github.com/LayssonENS/go-genealogy-api/domain"
	"time"
)

const dateLayout = "2006-01-02"

type postgresPersonRepo struct {
	DB *sql.DB
}

// NewPostgresPersonRepository will create an implementation of person.Repository
func NewPostgresPersonRepository(db *sql.DB) domain.PersonRepository {
	return &postgresPersonRepo{
		DB: db,
	}
}

func (p *postgresPersonRepo) GetByID(id int64) (domain.Person, error) {
	var person domain.Person
	err := p.DB.QueryRow(
		"SELECT id, name, email, birth_date, created_at FROM person WHERE id = $1", id).Scan(
		&person.ID, &person.Name, &person.Email, &person.BirthDate, &person.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return person, domain.ErrRegistrationNotFound
		}
		return person, err
	}
	return person, nil

}

func (p *postgresPersonRepo) CreatePerson(person domain.PersonRequest) error {
	date, _ := time.Parse(dateLayout, person.BirthDate)
	birthDate := date

	query := `INSERT INTO person (name, email, birth_date) VALUES ($1, $2, $3) `
	_, err := p.DB.Exec(query, person.Name, person.Email, birthDate)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresPersonRepo) GetAllPerson() ([]domain.Person, error) {
	var people []domain.Person

	rows, err := p.DB.Query("SELECT id, name, email, birth_date, created_at FROM person")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var person domain.Person
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Email,
			&person.BirthDate,
			&person.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	if len(people) == 0 {
		return nil, domain.ErrRegistrationNotFound
	}

	return people, nil

}
