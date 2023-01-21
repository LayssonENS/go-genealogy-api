package usecase

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
)

type relationshipUseCase struct {
	relationshipRepository domain.RelationshipRepository
}

func NewRelationshipUseCase(relationshipRepository domain.RelationshipRepository) domain.RelationshipUseCase {
	return &relationshipUseCase{
		relationshipRepository: relationshipRepository,
	}
}

func (a *relationshipUseCase) GetRelationshipByID(c *gin.Context, id int64) (*domain.FamilyMembers, error) {
	person, err := a.relationshipRepository.GetRelationshipByID(c, id)
	if err != nil {
		return person, err
	}

	return person, nil
}

func (a *relationshipUseCase) CreateRelationship(c *gin.Context, relationship domain.Relationship) error {
	err := a.relationshipRepository.CreateRelationship(c, relationship)
	if err != nil {
		return err
	}

	return nil
}
