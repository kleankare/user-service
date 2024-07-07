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

type httpRoleHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *services.MockRoleService
}

func (suite *httpRoleHandlerTestSuite) SetupSuite() {
	suite.router = gin.Default()

	suite.mockService = new(services.MockRoleService)
	handler := handlers.NewHttpRoleHandler(suite.mockService)

	suite.router.POST("/roles", handler.CreateRole)
	suite.router.GET("/roles", handler.ReadRoles)
	suite.router.GET("/roles/:id", handler.ReadRole)
	suite.router.PUT("/roles/:id", handler.UpdateRole)
	suite.router.DELETE("/roles/:id", handler.DeleteRole)
}

func TestHttpRoleHandler(t *testing.T) {
	suite.Run(t, new(httpRoleHandlerTestSuite))
}

func (suite *httpRoleHandlerTestSuite) TestCreateRole_Success() {
	suite.mockService.
		On("CreateRole", mock.AnythingOfType("string")).
		Return(&domain.Role{}, nil).
		Once()

	w := httptest.NewRecorder()
	exampleRole := domain.Role{
		Name: "name",
	}
	roleJson, _ := json.Marshal(exampleRole)
	req := httptest.NewRequest("POST", "/roles", strings.NewReader(string(roleJson)))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpRoleHandlerTestSuite) TestReadRole_Success() {
	suite.mockService.
		On("ReadRoles").
		Return([]*domain.Role{}, nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roles", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpRoleHandlerTestSuite) TestReadRoles_Success() {
	suite.mockService.
		On("ReadRole", mock.AnythingOfType("string")).
		Return(&domain.Role{}, nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roles/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpRoleHandlerTestSuite) TestUpdateRoles_Success() {
	suite.mockService.
		On("UpdateRole", mock.AnythingOfType("string"), mock.AnythingOfType("domain.Role")).
		Return(&domain.Role{}, nil).
		Once()

	w := httptest.NewRecorder()
	exampleRole := domain.Role{
		Name: "name",
	}
	roleJson, _ := json.Marshal(exampleRole)
	req := httptest.NewRequest("PUT", "/roles/1", strings.NewReader(string(roleJson)))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *httpRoleHandlerTestSuite) TestDeleteRoles_Success() {
	suite.mockService.
		On("DeleteRole", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/roles/1", nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}
