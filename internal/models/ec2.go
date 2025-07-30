package models

import (
	"github.com/google/uuid"
	"time"
)

type EC2 struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Tag       string    `json:"tag" gorm:"not null"`
	Ram       float64   `json:"ram" gorm:"not null"`
	Cpu       int       `json:"cpu" gorm:"not null"`
	Storage   float64   `json:"storage" gorm:"not null"`
	Bandwidth float64   `json:"bandwidth" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
