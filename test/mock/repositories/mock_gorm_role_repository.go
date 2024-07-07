package repositories

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockGormRoleRepository struct {
	mock.Mock
}

func (m *MockGormRoleRepository) CreateRole(name string) (*domain.Role, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockGormRoleRepository) ReadRoles() ([]*domain.Role, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (m *MockGormRoleRepository) ReadRole(id string) (*domain.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockGormRoleRepository) UpdateRole(id string, role domain.Role) (*domain.Role, error) {
	args := m.Called(id, role)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockGormRoleRepository) DeleteRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
