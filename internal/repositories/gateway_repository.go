package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/models"
)

type GatewayRepository struct {
	db *gorm.DB
}

func NewGatewayRepository(db *gorm.DB) *GatewayRepository {
	return &GatewayRepository{
		db: db,
	}
}

func (repo *GatewayRepository) CreateGateway(gateway *models.Gateway) error {
	result := repo.db.Create(&gateway)

	return result.Error
}

func (repo *GatewayRepository) GetUserGateways(userID uuid.UUID) error {
	return nil
}

func (repo *GatewayRepository) GetGateway(id uuid.UUID) {}
