package services

import (
	"fmt"
	"time"

	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/kleankare/user-service/internal/logger"
)

type roleService struct {
	repo  ports.RoleRepository
	cache ports.CacheRepository
}

func NewRoleService(repo ports.RoleRepository, cache ports.CacheRepository) ports.RoleService {
	return &roleService{
		repo:  repo,
		cache: cache,
	}
}

func (service *roleService) CreateRole(name string) (*domain.Role, error) {
	return service.repo.CreateRole(name)
}

func (service *roleService) ReadRoles() ([]*domain.Role, error) {
	return service.repo.ReadRoles()
}

func (service *roleService) ReadRole(id string) (*domain.Role, error) {
	cacheKey := fmt.Sprintf("role:%s", id)
	var (
		role *domain.Role
		err  error
	)

	if err = service.cache.Get(cacheKey, &role); err == nil {
		return role, nil
	}
	if role, err = service.repo.ReadRole(id); err != nil {
		return nil, err
	}
	if err = service.cache.Set(cacheKey, role, time.Minute*10); err != nil {
		logger.Log.Error("Error storing role in cache: %v", err)
	}
	return role, nil
}

func (service *roleService) UpdateRole(id string, role domain.Role) (*domain.Role, error) {
	updatedRole, err := service.repo.UpdateRole(id, role)
	if err != nil {
		return nil, err
	}
	if err := service.cache.Delete(id); err != nil {
		logger.Log.Error("Error deleting role from cache: %v", err)
	}
	return updatedRole, nil
}

func (service *roleService) DeleteRole(id string) error {
	if err := service.repo.DeleteRole(id); err != nil {
		return err
	}
	if err := service.cache.Delete(id); err != nil {
		logger.Log.Error("Error deleting role from cache: %v", err)
	}
	return nil
}
