package services

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(username string, password string) (*domain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) ReadUsers() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserService) ReadUser(id string) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, user domain.User) (*domain.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
