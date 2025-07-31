package services

import (
	"net/http"

	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
)

type EC2Service struct {
	awsCredentialsService *AWSCredentialsService
	ec2Repository         *repositories.EC2Repository
}

func NewEC2Service(awsCredentialsService *AWSCredentialsService, ec2Repository *repositories.EC2Repository) *EC2Service {
	return &EC2Service{
		awsCredentialsService: awsCredentialsService,
		ec2Repository:         ec2Repository,
	}
}

func (s *EC2Service) GetEC2InstanceTypes(params *middleware.QueryParams) (*[]models.EC2, int, error) {
	ec2Instances, err := s.ec2Repository.GetEC2InstancesTypes(params)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return ec2Instances, http.StatusOK, nil
}

func (s *EC2Service) GetEC2InstancesTypeCount() (int64, error) {
	count, err := s.ec2Repository.GetEC2TotalCount()

	if err != nil {
		return 0, err
	}

	return count, nil
}
