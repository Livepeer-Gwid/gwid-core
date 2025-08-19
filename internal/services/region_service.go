package services

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/types"
)

type RegionService struct {
	awsCredentialsService *AWSCredentialsService
	encryptionService     *EncryptionService
}

func NewRegionService(awsCredentialsService *AWSCredentialsService, encryptionService *EncryptionService) *RegionService {
	return &RegionService{
		awsCredentialsService: awsCredentialsService,
		encryptionService:     encryptionService,
	}
}

func (s *RegionService) GetAWSRegions(userID uuid.UUID, credentialsID uuid.UUID) ([]*types.RegionRes, int, error) {
	ctx := context.Background()

	credential, statusCode, err := s.awsCredentialsService.GetAWSCredentialsByID(credentialsID, userID)
	if err != nil {
		return nil, statusCode, err
	}

	awsRegion := "eu-central-1"

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(credential.AccessKeyID, credential.SecretAccessKey, "")))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	ec2Client := ec2.NewFromConfig(cfg)

	shouldGetAllRegions := true

	regions, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: &shouldGetAllRegions,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var data []*types.RegionRes

	for _, region := range regions.Regions {
		regionRes := types.RegionRes{
			ID:         uuid.New(),
			RegionName: *region.RegionName,
			Status:     *region.OptInStatus,
			Endpoint:   *region.Endpoint,
		}

		data = append(data, &regionRes)
	}

	return data, http.StatusOK, nil
}
