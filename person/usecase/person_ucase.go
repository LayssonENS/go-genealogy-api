package usecase

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
)

type personUseCase struct {
	personRepository domain.PersonRepository
}

func NewPersonUseCase(personRepository domain.PersonRepository) domain.PersonUseCase {
	return &personUseCase{
		personRepository: personRepository,
	}
}

func (a *personUseCase) GetByID(id int64) (domain.Person, error) {
	person, err := a.personRepository.GetByID(id)
	if err != nil {
		return person, err
	}

	return person, nil
}

func (a *personUseCase) CreatePerson(person domain.PersonRequest) error {
	err := a.personRepository.CreatePerson(person)
	if err != nil {
		return err
	}

	return nil
}

func (a *personUseCase) GetAllPerson() ([]domain.Person, error) {
	person, err := a.personRepository.GetAllPerson()
	if err != nil {
		return person, err
	}

	return person, nil
}
