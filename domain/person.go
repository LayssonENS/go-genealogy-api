package domain

import (
	"time"
)

type Person struct {
	ID          int64      `json:"id" swagger:"ignore"`
	Name        string     `json:"name" binding:"required"`
	Email       string     `json:"email"`
	DateOfBirth *time.Time `json:"date_of_birth" time_format:"2006-01-02"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PersonUseCase interface {
	GetByID(id int64) (Person, error)
	CreatePerson(person Person) error
	GetAllPerson() ([]Person, error)
}

type PersonRepository interface {
	GetByID(id int64) (Person, error)
	CreatePerson(person Person) error
	GetAllPerson() ([]Person, error)
}
