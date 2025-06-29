// Package repositories contains handlers that directly interacts with database
package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, *gorm.DB)
	FindByID(id uuid.UUID) (*models.User, *gorm.DB)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) CreateUser(user *models.User) error {
	result := repo.db.Create(&user)

	return result.Error
}

func (repo *userRepository) FindByEmail(email string) (*models.User, *gorm.DB) {
	var user models.User

	result := repo.db.Where(&models.User{Email: email}).First(&user)

	return &user, result
}

func (repo *userRepository) FindByID(id uuid.UUID) (*models.User, *gorm.DB) {
	var user models.User

	result := repo.db.Where(&models.User{ID: id}).First(&user)

	return &user, result
}
