package repositories

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockGormUserRepository struct {
	mock.Mock
}

func (m *MockGormUserRepository) CreateUser(username string, password string) (*domain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockGormUserRepository) ReadUsers() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockGormUserRepository) ReadUser(id string) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockGormUserRepository) UpdateUser(id string, user domain.User) (*domain.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockGormUserRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
