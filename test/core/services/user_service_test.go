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

type userServiceTestSuite struct {
	suite.Suite
	mockRepo      *repositories.MockGormUserRepository
	mockRedisRepo *repositories.MockRedisRepository
	service       ports.UserService
}

func (suite *userServiceTestSuite) SetupSuite() {
	suite.mockRepo = new(repositories.MockGormUserRepository)
	suite.mockRedisRepo = new(repositories.MockRedisRepository)
	suite.service = services.NewUserService(suite.mockRepo, suite.mockRedisRepo)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(userServiceTestSuite))
}

func (suite *userServiceTestSuite) TestCreateUser_Success() {
	suite.mockRepo.
		On("CreateUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&domain.User{}, nil).
		Once()

	_, err := suite.service.CreateUser("username", "password")
	suite.NoError(err)
}

func (suite *userServiceTestSuite) TestReadUsers_Success() {
	suite.mockRepo.
		On("ReadUsers").
		Return([]*domain.User{}, nil).
		Once()

	_, err := suite.service.ReadUsers()
	suite.NoError(err)
}

func (suite *userServiceTestSuite) TestReadUser_Success() {
	suite.mockRedisRepo.On("Get", mock.AnythingOfType("string"), mock.Anything).
		Return(nil).
		Once()
	suite.mockRepo.
		On("ReadUser", mock.AnythingOfType("string")).
		Return(&domain.User{}, nil).
		Once()
	suite.mockRedisRepo.On("Set", mock.AnythingOfType("string"), mock.Anything, time.Second).
		Return(nil).
		Once()

	_, err := suite.service.ReadUser("1")
	suite.NoError(err)
}

func (suite *userServiceTestSuite) TestUpdateUser_Success() {
	suite.mockRepo.
		On("UpdateUser", mock.AnythingOfType("string"), mock.AnythingOfType("domain.User")).
		Return(&domain.User{}, nil).
		Once()
	suite.mockRedisRepo.On("Delete", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	_, err := suite.service.UpdateUser("1", domain.User{
		Username: "username",
		Password: "password",
	})
	suite.NoError(err)
}

func (suite *userServiceTestSuite) TestDeleteUser_Success() {
	suite.mockRepo.
		On("DeleteUser", mock.AnythingOfType("string")).
		Return(nil).
		Once()
	suite.mockRedisRepo.On("Delete", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	err := suite.service.DeleteUser("1")
	suite.NoError(err)
}
