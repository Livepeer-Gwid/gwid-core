package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EC2 struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Tag             string    `json:"tag" gorm:"not null"`
	Ram             float64   `json:"ram" gorm:"not null"`
	Cpu             int       `json:"cpu" gorm:"not null"`
	Architecture    string    `json:"architecture" gorm:"not null"`
	CpuManufacturer string    `json:"cpu_manufacturer" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ec2 *EC2) BeforeCreate(tx *gorm.DB) (err error) {
	ec2.ID = uuid.New()

	return nil
}
