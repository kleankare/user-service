package ports

import (
	"github.com/kleankare/user-service/internal/core/domain"
)

type UserService interface {
	CreateUser(username string, password string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	UpdateUser(id string, user domain.User) (*domain.User, error)
	DeleteUser(id string) error
}

type UserRepository interface {
	CreateUser(username string, password string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	UpdateUser(id string, user domain.User) (*domain.User, error)
	DeleteUser(id string) error
}
