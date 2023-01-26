package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/LayssonENS/go-genealogy-api/domain"
	mock_relationship "github.com/LayssonENS/go-genealogy-api/relationships/mocks"
	usecase2 "github.com/LayssonENS/go-genealogy-api/relationships/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRelationshipHandler_GetRelationshipByID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_relationship.NewMockRelationshipRepository(ctrl)
	mockRepo.EXPECT().GetRelationshipByID(gomock.Any()).Return(&domain.Member{
		XMLName: xml.Name{},
		ID:      1,
		Members: []domain.Family{
			{
				ID: 2,
			}},
	}, nil)

	// Criação da instância do handler com o mocks do repositório
	handler := RelationshipHandler{RUseCase: usecase2.NewRelationshipUseCase(mockRepo)}

	// Criação do contexto de teste
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Writer.Header().Set("Content-Type", "application/json")

	// Configuração do parâmetro de ID na requisição
	c.Params = gin.Params{gin.Param{Key: "personId", Value: "1"}}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Add("Accept", "application/json")

	// Execução da rota
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Error occurred while marshalling JSON", r)
		}
	}()

	handler.GetRelationshipByID(c)

	// Verificação do status da resposta
	assert.Equal(t, http.StatusOK, c.Writer.Status())

	// Verificação do conteúdo da resposta
	var response domain.Member
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		return
	}

	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, int64(2), response.Members[0].ID)

}

func TestRelationshipHandler_CreateRelationship(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_relationship.NewMockRelationshipRepository(ctrl)
	mockRepo.EXPECT().CreateRelationship(gomock.Any(), gomock.Any()).Return(nil)

	// Criação da instância do handler com o mocks do repositório
	handler := RelationshipHandler{RUseCase: usecase2.NewRelationshipUseCase(mockRepo)}

	// Criação do contexto de teste
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Writer.Header().Set("Content-Type", "application/json")

	// Configuração do parâmetro de ID na requisição
	c.Params = gin.Params{gin.Param{Key: "personId", Value: "1"}}

	jsonRequest := `{
						"children":1,
						"parent":2
					}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(jsonRequest)))
	c.Writer.Header().Set("Content-Type", "application/json")

	// Execução da rota
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Error occurred while marshalling JSON", r)
		}
	}()

	handler.CreateRelationship(c)

	// Verificação do status da resposta
	assert.Equal(t, http.StatusCreated, c.Writer.Status())

}
