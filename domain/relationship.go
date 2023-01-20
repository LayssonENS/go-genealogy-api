package domain

import (
	"github.com/gin-gonic/gin"
)

type Relationship struct {
	ID              int    `json:"id"`
	PersonID        int    `json:"person_id"`
	RelatedPersonID int    `json:"related_person_id"`
	Relationship    string `json:"relationship"`
}

type RelationshipUseCase interface {
	GetRelationshipByID(c *gin.Context, id int64) (Relationship, error)
	CreateRelationship(c *gin.Context, person Relationship) error
}

type RelationshipRepository interface {
	GetRelationshipByID(c *gin.Context, id int64) (Relationship, error)
	CreateRelationship(c *gin.Context, relationship Relationship) error
}
