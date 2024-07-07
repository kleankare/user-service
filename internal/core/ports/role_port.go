package ports

import "github.com/kleankare/user-service/internal/core/domain"

type RoleService interface {
	CreateRole(name string) (*domain.Role, error)
	ReadRoles() ([]*domain.Role, error)
	ReadRole(id string) (*domain.Role, error)
	UpdateRole(id string, role domain.Role) (*domain.Role, error)
	DeleteRole(id string) error
}

type RoleRepository interface {
	CreateRole(name string) (*domain.Role, error)
	ReadRoles() ([]*domain.Role, error)
	ReadRole(id string) (*domain.Role, error)
	UpdateRole(id string, role domain.Role) (*domain.Role, error)
	DeleteRole(id string) error
}
