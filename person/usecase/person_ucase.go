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

func (a *personUseCase) GetByID(c *gin.Context, id int64) (domain.Person, error) {
	person, err := a.personRepository.GetByID(c, id)
	if err != nil {
		return person, err
	}

	return person, nil
}

func (a *personUseCase) CreatePerson(c *gin.Context, person domain.Person) error {
	err := a.personRepository.CreatePerson(c, person)
	if err != nil {
		return err
	}

	return nil
}
