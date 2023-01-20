package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Person struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type PersonUseCase interface {
	GetByID(c *gin.Context, id int64) (Person, error)
	CreatePerson(c *gin.Context, person Person) error
}

type PersonRepository interface {
	GetByID(c *gin.Context, id int64) (Person, error)
	CreatePerson(c *gin.Context, person Person) error
}
