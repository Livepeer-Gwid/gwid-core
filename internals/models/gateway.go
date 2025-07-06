package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Gateway struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	GatwayType         string    `json:"gatway_type" gorm:"not null"`
	Provider           string    `json:"provider" gorm:"not null"`
	Region             string    `json:"region" gorm:"not null"`
	GatewatName        string    `json:"gatewat_name" gorm:"not null"`
	RPCURL             string    `json:"rpc_url" gorm:"not null"`
	Password           string    `json:"-" gorm:"not null"`
	TranscodingProfile string    `json:"transcoding_profile" gorm:"not null"`
	QueueID            string    `json:"queue_id" gorm:"not null"`
	UserID             uuid.UUID `json:"user_id" gorm:"index"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"user" gorm:"references.ID"`
}

func (gateway *Gateway) BeforeCreate(tx *gorm.DB) (err error) {
	gateway.ID = uuid.New()

	gateway.HashPassword(gateway.Password)

	return nil
}

func (gateway *Gateway) HashPassword(password string) error {
	hashedPasswordbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("unable to hash password")
	}

	gateway.Password = string(hashedPasswordbytes)

	return nil
}
