package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
)

type AWSCredentialsService struct {
	awsCredentialsRepository *repositories.AWSCredentialsRepository
	encryptionService        *EncryptionService
}

func NewAWSCredentialsService(
	awsCredentialsRepository *repositories.AWSCredentialsRepository,
	encryptionService *EncryptionService,
) *AWSCredentialsService {
	return &AWSCredentialsService{
		awsCredentialsRepository: awsCredentialsRepository,
		encryptionService:        encryptionService,
	}
}

func (s *AWSCredentialsService) ValidateAWSCredentials(accessKeyID, secretAccessKey, region string) bool {
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return false
	}

	stsClient := sts.NewFromConfig(cfg)

	if _, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{}); err != nil {
		return false
	}

	return true
}

func (s *AWSCredentialsService) CreateAWSCredentials(credentials *models.AWSCredentials) (int, error) {
	_, result := s.awsCredentialsRepository.GetCredentialsByAccessKeyID(credentials.AccessKeyID)

	if result.RowsAffected > 0 {
		return http.StatusBadRequest, errors.New("aws credentials already exisits")
	}

	isValid := s.ValidateAWSCredentials(credentials.AccessKeyID, credentials.SecretAccessKey, "eu-central-1")

	if !isValid {
		return http.StatusBadRequest, errors.New("invalid AWS credentials")
	}

	if encryptedSecretAccessKey, err := s.encryptionService.EncryptData([]byte(credentials.SecretAccessKey)); err != nil {
		return http.StatusBadRequest, err
	} else {
		credentials.SecretAccessKey = encryptedSecretAccessKey
	}

	if err := s.awsCredentialsRepository.CreateAWSCredentials(credentials); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusCreated, nil
}

func (s *AWSCredentialsService) GetUserAWSCredentials(userID uuid.UUID, params *middleware.QueryParams) (*[]models.AWSCredentials, int, error) {
	credentials, err := s.awsCredentialsRepository.GetUserCredentials(userID, params)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return credentials, http.StatusOK, nil
}

func (s *AWSCredentialsService) GetUserCredentialsCount(userID uuid.UUID) (int64, error) {
	count, err := s.awsCredentialsRepository.GetUserCredentialsCount(userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *AWSCredentialsService) GetAWSCredentialsByID(id uuid.UUID, userID uuid.UUID) (*models.AWSCredentials, int, error) {
	credential, result := s.awsCredentialsRepository.GetCredentialsByID(id, userID)

	if result.RowsAffected == 0 {
		return nil, http.StatusNotFound, errors.New("credentials not found")
	}

	decryptedSecretAccessKey, err := s.encryptionService.DecryptData(credential.SecretAccessKey)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	credential.SecretAccessKey = decryptedSecretAccessKey

	return credential, http.StatusOK, nil
}
