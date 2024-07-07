package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kleankare/user-service/internal/adapters/handlers"
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/test/mock/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type httpUserHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *services.MockUserService
}

func (suite *httpUserHandlerTestSuite) SetupSuite() {
	suite.router = gin.Default()

	suite.mockService = new(services.MockUserService)
	handler := handlers.NewHttpUserHandler(suite.mockService)

	suite.router.POST("/users", handler.CreateUser)
	suite.router.GET("/users", handler.ReadUsers)
	suite.router.GET("/users/:id", handler.ReadUser)
	suite.router.PUT("/users/:id", handler.UpdateUser)
	suite.router.DELETE("/users/:id", handler.DeleteUser)
}

func TestHttpUserHandler(t *testing.T) {
	suite.Run(t, new(httpUserHandlerTestSuite))
}

func (suite *httpUserHandlerTestSuite) TestCreateUser_Success() {
	suite.mockService.
		On("CreateUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&domain.User{}, nil).
		Once()

	w := httptest.NewRecorder()
	exampleUser := domain.User{
		Username: "username",
		Password: "password",
	}
	userJson, _ := json.Marshal(exampleUser)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(userJson)))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpUserHandlerTestSuite) TestReadUsers_Success() {
	suite.mockService.
		On("ReadUsers").
		Return([]*domain.User{}, nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpUserHandlerTestSuite) TestReadUser_Success() {
	suite.mockService.
		On("ReadUser", mock.AnythingOfType("string")).
		Return(&domain.User{}, nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpUserHandlerTestSuite) TestUpdateUser_Success() {
	suite.mockService.
		On("UpdateUser", mock.AnythingOfType("string"), mock.AnythingOfType("domain.User")).
		Return(&domain.User{}, nil).
		Once()

	w := httptest.NewRecorder()
	exampleUser := domain.User{
		Username: "username",
		Password: "password",
	}
	userJson, _ := json.Marshal(exampleUser)
	req := httptest.NewRequest("PUT", "/users/1", strings.NewReader(string(userJson)))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpUserHandlerTestSuite) TestDeleteUser_Success() {
	suite.mockService.
		On("DeleteUser", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/users/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}
