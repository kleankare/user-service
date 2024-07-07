package services

import (
	"fmt"
	"time"

	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/kleankare/user-service/internal/logger"
)

type userService struct {
	repo  ports.UserRepository
	cache ports.CacheRepository
}

func NewUserService(repo ports.UserRepository, cache ports.CacheRepository) ports.UserService {
	return &userService{
		repo:  repo,
		cache: cache,
	}
}

func (service *userService) CreateUser(username string, password string) (*domain.User, error) {
	return service.repo.CreateUser(username, password)
}

func (service *userService) ReadUsers() ([]*domain.User, error) {
	return service.repo.ReadUsers()
}

func (service *userService) ReadUser(id string) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:%s", id)
	var (
		user *domain.User
		err  error
	)
	if err = service.cache.Get(cacheKey, &user); err == nil {
		return user, nil
	}
	if user, err = service.repo.ReadUser(id); err != nil {
		return nil, err
	}
	if err = service.cache.Set(cacheKey, user, time.Minute*10); err != nil {
		logger.Log.Error("Error storing role in cache: %v", err)
	}
	return user, nil
}

func (service *userService) UpdateUser(id string, user domain.User) (*domain.User, error) {
	updatedUser, err := service.repo.UpdateUser(id, user)
	if err != nil {
		return nil, err
	}
	if err := service.cache.Delete(id); err != nil {
		logger.Log.Error("Error deleting user from cache: %v", err)
	}
	return updatedUser, nil
}

func (service *userService) DeleteUser(id string) error {
	if err := service.repo.DeleteUser(id); err != nil {
		return err
	}
	if err := service.cache.Delete(id); err != nil {
		logger.Log.Error("Error deleting user from cache: %v", err)
	}
	return nil
}
