package repositories

import (
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
)

type ReferralRewardRepository struct {
	db *gorm.DB
}

func NewReferralRewardRepository(db *gorm.DB) *ReferralRewardRepository {
	return &ReferralRewardRepository{
		db: db,
	}
}

func (repo *ReferralRewardRepository) CreateReferralRewardRepository(referralReward *models.ReferralReward) error {
	result := repo.db.Create(referralReward)

	return result.Error
}
