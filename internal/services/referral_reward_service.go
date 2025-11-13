package services

import (
	"errors"

	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
)

type ReferralRewardService struct {
	referralRewardRepositotry *repositories.ReferralRewardRepository
	userRepository            *repositories.UserRepository
}

func NewReferralRewardService(
	referralRewardRepository *repositories.ReferralRewardRepository,
	userRepository *repositories.UserRepository,
) *ReferralRewardService {
	return &ReferralRewardService{
		referralRewardRepositotry: referralRewardRepository,
		userRepository:            userRepository,
	}
}

func (s *ReferralRewardService) CreateReferralReward(referralCode string) error {
	user, result := s.userRepository.FindByReferralCode(referralCode)

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	referralReward := models.ReferralReward{
		UserID: user.ID,
	}

	if err := s.referralRewardRepositotry.CreateReferralRewardRepository(&referralReward); err != nil {
		return err
	}

	return nil
}
