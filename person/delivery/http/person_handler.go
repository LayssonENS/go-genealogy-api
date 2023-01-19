package http

import (
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
	routerGroup.GET("/articles", handler.GetByID)
}

func (h PersonHandler) GetByID(context *gin.Context) {

}
