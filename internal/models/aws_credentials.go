package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AWSCredentials struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	AccessKeyID     string    `json:"access_key_id" gorm:"not null;uniqueIndex"`
	SecretAccessKey string    `json:"secret_access_key" gorm:"not null"`
	RoleName        string    `json:"role_name" gorm:"not null"`
	RoleARN         string    `json:"role_arn" gorm:"not null"`
	ProfieName      string    `json:"profie_name" gorm:"not null"`
	ProfileARN      string    `json:"profile_arn" gorm:"not null"`

	UserID uuid.UUID `json:"user_id" gorm:"index"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *User `json:"user" gorm:"foreignKey.UserID"`
}

func (awsCredentials *AWSCredentials) BeforeCreate(db *gorm.DB) (err error) {
	awsCredentials.ID = uuid.New()

	return
}
