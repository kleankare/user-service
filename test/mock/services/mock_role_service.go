package services

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) CreateRole(name string) (*domain.Role, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) ReadRoles() ([]*domain.Role, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (m *MockRoleService) ReadRole(id string) (*domain.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) UpdateRole(id string, role domain.Role) (*domain.Role, error) {
	args := m.Called(id, role)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) DeleteRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
