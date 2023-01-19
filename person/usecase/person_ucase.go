package usecase

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
)

type personUseCase struct {
	personRepository domain.PersonRepository
}

func NewPersonUseCase(personRepository domain.PersonRepository) domain.PersonUseCase {
	return &personUseCase{
		personRepository: personRepository,
	}
}

func (a *personUseCase) GetByID(c *gin.Context, id int64) (res domain.Person, err error) {

	return
}
