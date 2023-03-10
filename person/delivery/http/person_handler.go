package http

import (
	"net/http"
	"strconv"

	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	PUseCase domain.PersonUseCase
}

func NewPersonHandler(routerGroup *gin.Engine, us domain.PersonUseCase) {
	handler := &PersonHandler{
		PUseCase: us,
	}

	routerGroup.GET("/v1/person/:personId", handler.GetByID)
	routerGroup.GET("/v1/person/all", handler.GetAllPerson)
	routerGroup.POST("/v1/person", handler.CreatePerson)
}

// GetByID godoc
// @Summary Route to fetch person by ID
// @Description Fetch person
// @Tags Person
// @Accept  json
// @Produce  json
// @Param personId path int true "Person ID"
// @Success 200 {object} domain.Person
// @Failure 400 {object} domain.ErrorResponse
// @Router /v1/person/{personId} [GET]
func (h *PersonHandler) GetByID(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("personId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	personId := int64(idParam)

	response, err := h.PUseCase.GetByID(personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllPerson godoc
// @Summary Route to return all people
// @Description All people
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Person
// @Failure 400 {object} domain.ErrorResponse
// @Router /v1/person/all [GET]
func (h *PersonHandler) GetAllPerson(c *gin.Context) {
	response, err := h.PUseCase.GetAllPerson()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreatePerson godoc
// @Summary Route to create person
// @Description Create person
// @Tags Person
// @Accept  json
// @Produce  json
// @Param Payload body domain.PersonRequest true "Payload"
// @Success 201 {object} string
// @Failure 400 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Router /v1/person [POST]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var person domain.PersonRequest
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	err := h.PUseCase.CreatePerson(person)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Created"})
}
