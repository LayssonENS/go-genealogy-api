package http

import (
	"github.com/LayssonENS/go-genealogy-api/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	PUseCase domain.PersonUseCase
}

func NewPersonHandler(routerGroup *gin.Engine, us domain.PersonUseCase) {
	handler := &PersonHandler{
		PUseCase: us,
	}
	routerGroup.GET("/person/:id", handler.GetByID)
	routerGroup.POST("/person", handler.CreatePerson)
}

func (h *PersonHandler) GetByID(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNotFound.Error()})
		return
	}

	personId := int64(idParam)

	response, err := h.PUseCase.GetByID(c, personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var person domain.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err := h.PUseCase.CreatePerson(c, person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Created"})
}
