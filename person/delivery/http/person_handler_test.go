package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LayssonENS/go-genealogy-api/domain"
	mockDomain "github.com/LayssonENS/go-genealogy-api/person/mocks"
	"github.com/LayssonENS/go-genealogy-api/person/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
)

func TestGetByID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockDomain.NewMockPersonRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any()).Return(domain.Person{
		ID:        1,
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		BirthDate: &time.Time{},
		CreatedAt: time.Now(),
	}, nil)

	// Creating an instance of the handler with mock repository
	handler := PersonHandler{PUseCase: usecase.NewPersonUseCase(mockRepo)}

	// Creating the test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Writer.Header().Set("Content-Type", "application/json")

	// Setting the ID parameter in the request
	c.Params = gin.Params{gin.Param{Key: "personId", Value: "1"}}

	// Execution of the route
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Error occurred while marshalling JSON", r)
		}
	}()
	handler.GetByID(c)

	// Checking the status of the response
	assert.Equal(t, http.StatusOK, c.Writer.Status())

	// Checking the body of the response
	var response domain.Person
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		return
	}

	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "johndoe@example.com", response.Email)
}

func TestGetAllPerson(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockDomain.NewMockPersonRepository(ctrl)
	mockRepo.EXPECT().GetAllPerson().Return([]domain.Person{
		{
			ID:        1,
			Name:      "John Doe",
			Email:     "johndoe@example.com",
			BirthDate: &time.Time{},
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Laysson",
			Email:     "Laysson@example.com",
			BirthDate: &time.Time{},
			CreatedAt: time.Now(),
		},
	}, nil)

	// Creating an instance of the handler with mock repository
	handler := PersonHandler{PUseCase: usecase.NewPersonUseCase(mockRepo)}

	// Creating the test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Writer.Header().Set("Content-Type", "application/json")

	// Setting the ID parameter in the request
	c.Params = gin.Params{gin.Param{Key: "personId", Value: "1"}}

	// Execution of the route
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Error occurred while marshalling JSON", r)
		}
	}()

	handler.GetAllPerson(c)

	// Checking the status of the response
	assert.Equal(t, http.StatusOK, c.Writer.Status())

	// Checking the body of the response
	var response []domain.Person
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		return
	}

	assert.Equal(t, int64(1), response[0].ID)
	assert.Equal(t, int64(2), response[1].ID)
	assert.Equal(t, "John Doe", response[0].Name)
	assert.Equal(t, "Laysson", response[1].Name)
	assert.Equal(t, "Laysson@example.com", response[1].Email)
}

func TestCreatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockDomain.NewMockPersonRepository(ctrl)
	mockRepo.EXPECT().CreatePerson(domain.PersonRequest{
		Name:      "Marcia",
		Email:     "teste@teste.com",
		BirthDate: "1998-01-01",
	}).Return(nil)

	// Creating an instance of the handler with mock repository
	handler := PersonHandler{PUseCase: usecase.NewPersonUseCase(mockRepo)}

	// Creating the test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonRequest := `{
						"name":"Marcia",
						"email":"teste@teste.com",
						"birth_date":"1998-01-01"
					}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(jsonRequest)))
	c.Writer.Header().Set("Content-Type", "application/json")

	// Execution of the route
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Error occurred while marshalling JSON", r)
		}
	}()

	handler.CreatePerson(c)

	// Checking the status of the response
	assert.Equal(t, http.StatusCreated, c.Writer.Status())

}
