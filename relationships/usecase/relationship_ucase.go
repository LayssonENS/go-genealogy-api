package usecase

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
)

type relationshipUseCase struct {
	relationshipRepository domain.RelationshipRepository
}

func NewRelationshipUseCase(relationshipRepository domain.RelationshipRepository) domain.RelationshipUseCase {
	return &relationshipUseCase{
		relationshipRepository: relationshipRepository,
	}
}

func (a *relationshipUseCase) GetRelationshipByID(personId int64) (*domain.Member, error) {
	person, err := a.relationshipRepository.GetRelationshipByID(personId)
	if err != nil {
		return person, err
	}

	return person, nil
}

func (a *relationshipUseCase) CreateRelationship(personId int64, relationship domain.Relationship) error {
	err := a.relationshipRepository.CreateRelationship(personId, relationship)
	if err != nil {
		return err
	}

	return nil
}

func (a *relationshipUseCase) DeleteRelationship(personId int64) error {
	err := a.relationshipRepository.DeleteRelationship(personId)
	if err != nil {
		return err
	}

	return nil
}
