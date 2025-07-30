package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/uuid"
)

type EC2Service struct {
	awsCredentialsService *AWSCredentialsService
}

func NewEC2Service(awsCredentialsService *AWSCredentialsService) *EC2Service {
	return &EC2Service{
		awsCredentialsService: awsCredentialsService,
	}
}

func (s *EC2Service) GetEC2InstanceTypes(userID uuid.UUID, credentialsID uuid.UUID, region string) ([]types.InstanceTypeInfo, int, error) {
	credential, statusCode, err := s.awsCredentialsService.GetAWSCredentialsByID(credentialsID, userID)
	if err != nil {
		return nil, statusCode, err
	}

	creds := credentials.NewStaticCredentialsProvider(credential.AccessKeyID, credential.SecretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("unable to load aws config")
	}

	ec2Client := ec2.NewFromConfig(cfg)

	output, err := ec2Client.DescribeInstanceTypes(context.TODO(), &ec2.DescribeInstanceTypesInput{
		MaxResults: aws.Int32(10),
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("unable to get ec2 list")
	}

	ec2 := output.InstanceTypes

	return ec2, http.StatusOK, nil
}
