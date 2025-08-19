package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AWSRegion struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	RegionName string    `json:"region_name" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null"`
	Endpoint   string    `json:"endpoint" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (awsRegion *AWSRegion) BeforeCreate(tx *gorm.DB) (err error) {
	awsRegion.ID = uuid.New()

	return nil
}
