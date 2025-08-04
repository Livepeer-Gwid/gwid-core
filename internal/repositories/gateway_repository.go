package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/middleware"
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

func (repo *GatewayRepository) GetUserGateways(userID uuid.UUID, params *middleware.QueryParams) (*[]models.Gateway, error) {
	var gateways []models.Gateway

	result := repo.db.Offset(params.Offset).Limit(params.Limit).Order(params.Sort + " " + params.Order).Where(&models.Gateway{UserID: userID}).Find(&gateways)

	return &gateways, result.Error
}

func (repo *GatewayRepository) GetUserGatewaysCount(userID uuid.UUID) (int64, error) {
	var count int64

	result := repo.db.Model(&models.Gateway{}).Where(&models.Gateway{UserID: userID}).Count(&count)

	return count, result.Error
}

func (repo *GatewayRepository) GetGateway(id uuid.UUID) {}

func (repo *GatewayRepository) GetGatewayByName(name string) (*models.Gateway, *gorm.DB) {
	var gateway models.Gateway

	result := repo.db.Where(&models.Gateway{GatewayName: name}).Find(&gateway)

	return &gateway, result
}

func (repo *GatewayRepository) UpdateRepository(gateway *models.Gateway) error {
	result := repo.db.Updates(&gateway)

	return result.Error
}
