package repositories

import (
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (repo *gormUserRepository) CreateUser(username string, password string) (*domain.User, error) {
	user := domain.User{
		Username: username,
		Password: password,
	}

	if result := repo.db.Create(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *gormUserRepository) ReadUsers() ([]*domain.User, error) {
	var users []*domain.User

	if result := repo.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *gormUserRepository) ReadUser(id string) (*domain.User, error) {
	user := &domain.User{}

	if result := repo.db.First(&user, id); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *gormUserRepository) UpdateUser(id string, newUser domain.User) (*domain.User, error) {
	var user domain.User

	if result := repo.db.First(&user, id); result.Error != nil {
		return nil, result.Error
	}
	if result := repo.db.Model(&user).Updates(newUser); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (g *gormUserRepository) DeleteUser(id string) error {
	var user domain.User

	if result := g.db.Delete(&user, id); result.Error != nil {
		return result.Error
	}
	return nil
}
