package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
)

type AWSCredentialsService struct {
	awsCredentialsRepository *repositories.AWSCredentialsRepository
	encryptionService        *EncryptionService
}

const ssmTrustPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

func NewAWSCredentialsService(
	awsCredentialsRepository *repositories.AWSCredentialsRepository,
	encryptionService *EncryptionService,
) *AWSCredentialsService {
	return &AWSCredentialsService{
		awsCredentialsRepository: awsCredentialsRepository,
		encryptionService:        encryptionService,
	}
}

func (s *AWSCredentialsService) ValidateAWSCredentials(accessKeyID, secretAccessKey, region string) (*types.AWSCredentailsProfile, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)

	if _, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{}); err != nil {
		return nil, err
	}

	roleName := fmt.Sprintf("EC2-SSM-Role-%d", time.Now().Unix())
	roleArn, err := createSSMRole(iamClient, roleName)

	if err != nil {
		return nil, err
	}

	profileName := fmt.Sprintf("EC2-SSM-Profile-%d", time.Now().Unix())
	profileArn, err := createInstanceProfile(iamClient, profileName, roleName)

	if err != nil {
		return nil, err
	}

	return &types.AWSCredentailsProfile{
		ProfileName: profileName,
		ProfileARN:  *profileArn,
		RoleName:    roleName,
		RoleARN:     *roleArn,
	}, nil
}

func createSSMRole(iamClient *iam.Client, roleName string) (*string, error) {
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(ssmTrustPolicy),
		Description:              aws.String("Role for EC2 instance to use SSM"),
	}

	roleResult, err := iamClient.CreateRole(context.TODO(), createRoleInput)

	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	attachPolicyInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"),
	}

	_, err = iamClient.AttachRolePolicy(context.TODO(), attachPolicyInput)
	if err != nil {
		return nil, fmt.Errorf("failed to attach SSM policy: %w", err)
	}

	return roleResult.Role.Arn, nil
}

func createInstanceProfile(iamClient *iam.Client, profileName, roleName string) (*string, error) {
	createProfileInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}

	profileResult, err := iamClient.CreateInstanceProfile(context.TODO(), createProfileInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance profile: %w", err)
	}

	addRoleInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}

	_, err = iamClient.AddRoleToInstanceProfile(context.TODO(), addRoleInput)
	if err != nil {
		return nil, fmt.Errorf("failed to add role to instance profile: %w", err)
	}

	return profileResult.InstanceProfile.Arn, nil
}

func (s *AWSCredentialsService) CreateAWSCredentials(credentials *models.AWSCredentials) (int, error) {
	_, result := s.awsCredentialsRepository.GetCredentialsByAccessKeyID(credentials.AccessKeyID)

	if result.RowsAffected > 0 {
		return http.StatusBadRequest, errors.New("aws credentials already exisits")
	}

	awsCredentialsProfile, err := s.ValidateAWSCredentials(credentials.AccessKeyID, credentials.SecretAccessKey, "eu-central-1")

	if err != nil {
		return http.StatusBadRequest, errors.New("invalid AWS credentials")
	}

	credentials.ProfieName = awsCredentialsProfile.ProfileName
	credentials.ProfileARN = awsCredentialsProfile.ProfileARN
	credentials.RoleName = awsCredentialsProfile.RoleName
	credentials.RoleARN = awsCredentialsProfile.RoleARN

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
