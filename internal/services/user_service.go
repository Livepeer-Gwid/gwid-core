package services

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
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

func (s *UserService) GenerateUniqueReferralCode() (string, error) {
	var currentReferralCode string

	for {
		referralCode, err := utils.GenerateReferralCode()
		if err != nil {
			return "", err
		}

		if _, result := s.userRepository.FindByReferralCode(referralCode); result.RowsAffected == 0 {
			currentReferralCode = referralCode

			break
		}
	}

	return currentReferralCode, nil
}

func (s *UserService) UpdateUserProfile(updateProfileReq types.UpdateProfileReq, userID uuid.UUID) (*models.User, int, error) {
	user, result := s.userRepository.FindByID(userID)

	if result.RowsAffected == 0 {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	user.Name = updateProfileReq.Name

	if result := s.userRepository.UpdateUser(user); result.RowsAffected == 0 {
		return nil, http.StatusInternalServerError, errors.New("something went wrong")
	}

	return user, http.StatusOK, nil
}
