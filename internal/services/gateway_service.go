package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
)

type GatewayService struct {
	cfg                *config.Config
	gatewayTaskService *GatewayTaskService
	gatewayRepository  *repositories.GatewayRepository
	ec2Service         *EC2Service
}

func NewGatewayService(
	cfg *config.Config,
	gatewayTaskService *GatewayTaskService,
	gatewayRepository *repositories.GatewayRepository,
	ec2Service *EC2Service,
) *GatewayService {
	return &GatewayService{
		cfg:                cfg,
		gatewayTaskService: gatewayTaskService,
		gatewayRepository:  gatewayRepository,
		ec2Service:         ec2Service,
	}
}

func (s *GatewayService) getAsynqClient() *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: s.cfg.RedisAddress, Password: s.cfg.RedisPassword})

	return client
}

func (s *GatewayService) CreateGatewayWithAWS(createGatewayWithAWSReq types.CreateGatewayWithAWSReq, userID uuid.UUID) (*models.Gateway, int, error) {
	formattedGatewayName := utils.ToKebabCase(createGatewayWithAWSReq.GatewayName)

	_, result := s.gatewayRepository.GetGatewayByName(formattedGatewayName)

	if result.RowsAffected > 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("%s already existis", createGatewayWithAWSReq.GatewayName)
	}

	gateway := models.Gateway{
		Provider:           "aws",
		Region:             createGatewayWithAWSReq.Region,
		GatewayName:        formattedGatewayName,
		GatewayType:        createGatewayWithAWSReq.GatewayType,
		RPCURL:             createGatewayWithAWSReq.RPCURL,
		Password:           createGatewayWithAWSReq.Password,
		TranscodingProfile: createGatewayWithAWSReq.TranscodingProfile,
		UserID:             userID,
		AWSCredentialsID:   createGatewayWithAWSReq.CredentialsID,
	}

	err := s.gatewayRepository.CreateGateway(&gateway)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	ec2InstancePayload := types.CreateEC2InstanceReq{
		CredentialsID:     createGatewayWithAWSReq.CredentialsID,
		EC2InstanceTypeID: createGatewayWithAWSReq.EC2InstanceTypeID,
		InstanceName:      createGatewayWithAWSReq.GatewayName,
	}

	instanceID, statusCode, err := s.ec2Service.CreateEC2Instance(ec2InstancePayload, userID)

	if err != nil {
		gateway.Status = models.GatewayFailed

		err := s.gatewayRepository.UpdateRepository(&gateway)

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return nil, statusCode, err
	}

	gateway.InstanceID = &instanceID

	client := s.getAsynqClient()

	payload := types.DeployAWSGatewayPayload{
		GatewayID:        gateway.ID,
		UnhashedPassword: createGatewayWithAWSReq.Password,
		CredentialsID:    gateway.AWSCredentialsID,
		InstanceID:       instanceID,
		UserID:           userID,
		Region:           createGatewayWithAWSReq.Region,
	}

	task, err := s.gatewayTaskService.NewAWSDeployGatewayTask(payload)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	info, err := client.Enqueue(task)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("unable to queue task")
	}

	gateway.QueueID = &info.ID

	if err := s.gatewayRepository.UpdateRepository(&gateway); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer client.Close()

	return &gateway, http.StatusCreated, nil
}

func (s *GatewayService) GetUserGateways(userID uuid.UUID, params *middleware.QueryParams) (*[]models.Gateway, int, error) {
	gateways, err := s.gatewayRepository.GetUserGateways(userID, params)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gateways, http.StatusOK, nil
}

func (s *GatewayService) GetUserGatewaysCount(userID uuid.UUID) (int64, error) {
	count, err := s.gatewayRepository.GetUserGatewaysCount(userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
