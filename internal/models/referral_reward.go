package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReferralReward struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`

	UserID uuid.UUID `json:"user_id" gorm:"index"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User *User `json:"user" gorm:"foreignKey.UserID"`
}

func (referralReward *ReferralReward) BeforeCreate(tx *gorm.DB) (err error) {
	referralReward.ID = uuid.New()

	return nil
}
