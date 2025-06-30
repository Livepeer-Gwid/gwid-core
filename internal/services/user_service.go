package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
)

type UserService interface {
	GetCurrentUserProfile(userId uuid.UUID) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetCurrentUserProfile(userId uuid.UUID) (*models.User, error) {
	user, result := s.userRepository.FindByID(userId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return user, nil
}
