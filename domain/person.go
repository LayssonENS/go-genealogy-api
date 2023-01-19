package domain

import "github.com/gin-gonic/gin"

type Person struct {
	Name string
}

type PersonUseCase interface {
	GetByID(c *gin.Context, id int64) (Person, error)
}

type PersonRepository interface {
	GetByID(c *gin.Context, id int64) (Person, error)
}
