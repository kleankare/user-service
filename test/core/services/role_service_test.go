package test

import (
	"testing"
	"time"

	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/kleankare/user-service/internal/core/services"
	"github.com/kleankare/user-service/test/mock/repositories"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type roleServiceTestSuite struct {
	suite.Suite
	mockRepo      *repositories.MockGormRoleRepository
	mockRedisRepo *repositories.MockRedisRepository
	service       ports.RoleService
}

func (suite *roleServiceTestSuite) SetupSuite() {
	suite.mockRepo = new(repositories.MockGormRoleRepository)
	suite.mockRedisRepo = new(repositories.MockRedisRepository)
	suite.service = services.NewRoleService(suite.mockRepo, suite.mockRedisRepo)
}

func TestRoleService(t *testing.T) {
	suite.Run(t, new(roleServiceTestSuite))
}

func (suite *roleServiceTestSuite) TestCreateRole_Success() {
	suite.mockRepo.
		On("CreateRole", mock.AnythingOfType("string")).
		Return(&domain.Role{}, nil).
		Once()

	_, err := suite.service.CreateRole("name")
	suite.NoError(err)
}

func (suite *roleServiceTestSuite) TestReadRoles_Success() {
	suite.mockRepo.
		On("ReadRoles").
		Return([]*domain.Role{}, nil).
		Once()

	_, err := suite.service.ReadRoles()
	suite.NoError(err)
}

func (suite *roleServiceTestSuite) TestReadRole_Success() {
	suite.mockRedisRepo.On("Get", mock.AnythingOfType("string"), mock.Anything).
		Return(nil).
		Once()
	suite.mockRepo.
		On("ReadRole", mock.AnythingOfType("string")).
		Return(&domain.Role{}, nil).
		Once()
	suite.mockRedisRepo.On("Set", mock.AnythingOfType("string"), mock.Anything, time.Second).
		Return(nil).
		Once()

	_, err := suite.service.ReadRole("1")
	suite.NoError(err)
}

func (suite *roleServiceTestSuite) TestUpdateRole_Success() {
	suite.mockRepo.
		On("UpdateRole", mock.AnythingOfType("string"), mock.AnythingOfType("domain.Role")).
		Return(&domain.Role{}, nil).
		Once()
	suite.mockRedisRepo.On("Delete", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	_, err := suite.service.UpdateRole("1", domain.Role{
		Name: "test",
	})
	suite.NoError(err)
}

func (suite *roleServiceTestSuite) TestDeleteRole_Success() {
	suite.mockRepo.
		On("DeleteRole", mock.AnythingOfType("string")).
		Return(nil).
		Once()
	suite.mockRedisRepo.On("Delete", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	err := suite.service.DeleteRole("1")
	suite.NoError(err)
}
