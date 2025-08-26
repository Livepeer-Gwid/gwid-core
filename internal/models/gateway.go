package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type GatewayStatus string

const (
	GatewayInitializing GatewayStatus = "initializing"
	GatewayRunning      GatewayStatus = "running"
	GatewayStopped      GatewayStatus = "stopped"
	GatewayFailed       GatewayStatus = "failed"
)

type Gateway struct {
	ID                 uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;"`
	Provider           string        `json:"provider" gorm:"not null"`
	Region             string        `json:"region" gorm:"not null"`
	GatewayName        string        `json:"gateway_name" gorm:"uniqueIndex,not null"`
	GatewayType        string        `json:"gateway_type" gorm:"not null"`
	RPCURL             string        `json:"rpc_url" gorm:"not null"`
	Password           string        `json:"-" gorm:"not null"`
	TranscodingProfile string        `json:"transcoding_profile" gorm:"not null"`
	Status             GatewayStatus `json:"status" gorm:"default:'initializing';not null"`
	ErrorStatus        string        `json:"error_status"`
	QueueID            *string       `json:"queue_id"`
	InstanceID         *string       `json:"instance_id"`
	UserID             uuid.UUID     `json:"user_id" gorm:"index"`
	AWSCredentialsID   uuid.UUID     `json:"aws_credentials_id" gorm:"index"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User           *User           `json:"user" gorm:"foreignKey.UserID"`
	AWSCredentials *AWSCredentials `json:"aws_credentials" gorm:"foreignKey.AWSCredentialsID"`
}

func (gateway *Gateway) BeforeCreate(tx *gorm.DB) (err error) {
	gateway.ID = uuid.New()

	gateway.Status = GatewayInitializing

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
