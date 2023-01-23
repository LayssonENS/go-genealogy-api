package http

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RelationshipHandler struct {
	RUseCase domain.RelationshipUseCase
}

func NewRelationshipHandler(routerGroup *gin.Engine, us domain.RelationshipUseCase) {
	handler := &RelationshipHandler{
		RUseCase: us,
	}
	routerGroup.GET("/relationships/:personId", handler.GetRelationshipByID)
	routerGroup.POST("/relationships", handler.CreateRelationship)
}

// GetRelationshipByID godoc
// @Summary Route to Get relationships
// @Description Get relationships
// @Tags Relationship
// @Accept  json
// @Produce  json
// @Param personId path int true "Person ID"
// @Success 200 {object} domain.FamilyMembers
// @Failure 400 {object} string
// @Router /relationships/{personId} [GET]
func (h *RelationshipHandler) GetRelationshipByID(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("personId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	relationshipId := int64(idParam)

	response, err := h.RUseCase.GetRelationshipByID(c, relationshipId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateRelationship godoc
// @Summary Route to create relationships
// @Description Create relationships
// @Tags Relationship
// @Accept  json
// @Produce  json
// @Param Payload body domain.Relationship true "Payload"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 422 {object} string
// @Router /relationships [POST]
func (h *RelationshipHandler) CreateRelationship(c *gin.Context) {
	var relationship domain.Relationship
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err := h.RUseCase.CreateRelationship(c, relationship)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Created"})
}
