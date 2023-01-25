package domain

import (
	"time"
)

type PersonRequest struct {
	Name      string `json:"name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	BirthDate string `json:"birth_date"`
}

type Person struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	BirthDate *time.Time `json:"birth_date"`
	CreatedAt time.Time  `json:"created_at"`
}

type PersonUseCase interface {
	GetByID(id int64) (Person, error)
	CreatePerson(person PersonRequest) error
	GetAllPerson() ([]Person, error)
}

type PersonRepository interface {
	GetByID(id int64) (Person, error)
	CreatePerson(person PersonRequest) error
	GetAllPerson() ([]Person, error)
}
