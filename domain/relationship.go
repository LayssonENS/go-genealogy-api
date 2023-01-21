package domain

import (
	"github.com/gin-gonic/gin"
)

type Relationship struct {
	ID              int64  `json:"id"`
	PersonID        int64  `json:"person_id"`
	RelatedPersonID int64  `json:"related_person_id"`
	Relationship    string `json:"relationship"`
}

type FamilyMember struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
}

type FamilyMembers struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Members []FamilyMember `json:"members"`
}

type RelationshipUseCase interface {
	GetRelationshipByID(c *gin.Context, id int64) (*FamilyMembers, error)
	CreateRelationship(c *gin.Context, person Relationship) error
}

type RelationshipRepository interface {
	GetRelationshipByID(c *gin.Context, id int64) (*FamilyMembers, error)
	CreateRelationship(c *gin.Context, relationship Relationship) error
}
