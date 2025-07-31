package repositories

import (
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
)

type EC2Repository struct {
	db *gorm.DB
}

func NewEC2Repository(db *gorm.DB) *EC2Repository {
	return &EC2Repository{
		db: db,
	}
}

func (repo *EC2Repository) GetEC2InstancesTypes(params *middleware.QueryParams) (*[]models.EC2, error) {
	var ec2InstanceTypes []models.EC2

	result := repo.db.Offset(params.Offset).Limit(params.Limit).Order(params.Sort + " " + params.Order).Where(&models.EC2{}).Find(&ec2InstanceTypes)

	return &ec2InstanceTypes, result.Error
}

func (repo *EC2Repository) GetEC2TotalCount() (int64, error) {
	var count int64

	result := repo.db.Model(&models.EC2{}).Count(&count)

	return count, result.Error
}

func (repo *EC2Repository) CreateEC2InstanceTypes(instances *[]models.EC2) error {
	result := repo.db.Create(instances)

	return result.Error
}
