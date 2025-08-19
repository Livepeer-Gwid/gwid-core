// Package repositories contains handlers that directly interacts with database
package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	result := repo.db.Create(&user)

	return result.Error
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, *gorm.DB) {
	var user models.User

	result := repo.db.Where(&models.User{Email: email}).First(&user)

	return &user, result
}

func (repo *UserRepository) FindByID(id uuid.UUID) (*models.User, *gorm.DB) {
	var user models.User

	result := repo.db.Where(&models.User{ID: id}).First(&user)

	return &user, result
}
