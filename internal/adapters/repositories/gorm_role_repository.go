package repositories

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"gorm.io/gorm"
)

type gormRoleRepository struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) ports.RoleRepository {
	return &gormRoleRepository{
		db: db,
	}
}

func (repo *gormRoleRepository) CreateRole(name string) (*domain.Role, error) {
	role := domain.Role{
		Name: name,
	}

	if result := repo.db.Create(&role); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (repo *gormRoleRepository) ReadRoles() ([]*domain.Role, error) {
	var role []*domain.Role

	if result := repo.db.Find(&role); result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}

func (repo *gormRoleRepository) ReadRole(id string) (*domain.Role, error) {
	role := &domain.Role{}

	if result := repo.db.First(&role, id); result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}

func (repo *gormRoleRepository) UpdateRole(id string, newRole domain.Role) (*domain.Role, error) {
	var role domain.Role

	if result := repo.db.First(&role, id); result.Error != nil {
		return nil, result.Error
	}
	if result := repo.db.Model(&role).Updates(newRole); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (repo *gormRoleRepository) DeleteRole(id string) error {
	var role domain.Role

	if result := repo.db.Delete(&role, id); result.Error != nil {
		return result.Error
	}
	return nil
}
