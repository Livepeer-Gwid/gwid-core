package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
)

type AWSCredentialsRepository struct {
	db *gorm.DB
}

func NewAWSCredentialsRepository(db *gorm.DB) *AWSCredentialsRepository {
	return &AWSCredentialsRepository{
		db: db,
	}
}

func (repo *AWSCredentialsRepository) CreateAWSCredentials(awsCredentials *models.AWSCredentials) error {
	result := repo.db.Create(awsCredentials)

	return result.Error
}

func (repo *AWSCredentialsRepository) GetCredentialsByAccessKeyID(accessKeyID string) (*models.AWSCredentials, *gorm.DB) {
	var credentials models.AWSCredentials

	result := repo.db.Where(models.AWSCredentials{AccessKeyID: accessKeyID}).First(&credentials)

	return &credentials, result
}

func (repo *AWSCredentialsRepository) GetCredentialsByID(id uuid.UUID, userID uuid.UUID) (*models.AWSCredentials, *gorm.DB) {
	var credentials models.AWSCredentials

	result := repo.db.Where(models.AWSCredentials{ID: id, UserID: userID}).First(&credentials)

	return &credentials, result
}

func (repo *AWSCredentialsRepository) GetUserCredentials(userID uuid.UUID, params *middleware.QueryParams) (*[]models.AWSCredentials, error) {
	var credentials []models.AWSCredentials

	result := repo.db.Offset(params.Offset).Limit(params.Limit).Order(params.Sort + " " + params.Order).Where(&models.AWSCredentials{UserID: userID}).Find(&credentials)

	return &credentials, result.Error
}

func (repo *AWSCredentialsRepository) GetUserCredentialsCount(userID uuid.UUID) (int64, error) {
	var count int64

	result := repo.db.Model(&models.AWSCredentials{}).Where(&models.AWSCredentials{UserID: userID}).Count(&count)

	return count, result.Error
}
