package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetCurrentUserProfile(userID uuid.UUID) (*models.User, error) {
	user, result := s.userRepository.FindByID(userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return user, nil
}
