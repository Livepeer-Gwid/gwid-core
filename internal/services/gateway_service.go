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
	"gwid.io/gwid-core/internal/tasks"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
)

type GatewayService struct {
	cfg               *config.Config
	gatewayTask       *tasks.GatewayTask
	gatewayRepository *repositories.GatewayRepository
	ec2Service        *EC2Service
}

func NewGatewayService(
	cfg *config.Config,
	gatewayTask *tasks.GatewayTask,
	gatewayRepository *repositories.GatewayRepository,
	ec2Service *EC2Service,
) *GatewayService {
	return &GatewayService{
		cfg:               cfg,
		gatewayTask:       gatewayTask,
		gatewayRepository: gatewayRepository,
		ec2Service:        ec2Service,
	}
}

func (s *GatewayService) getGatewayClient() *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: s.cfg.RedisAddress})

	return client
}

func (s *GatewayService) CreateGatewayWithAWS(createGatewayWithAWSReq types.CreateGatewayWithAWSReq, userID uuid.UUID) (any, int, error) {
	formattedGatewayName := utils.ToKebabCase(createGatewayWithAWSReq.GatewayName)

	_, result := s.gatewayRepository.GetGatewayByName(formattedGatewayName)

	if result.RowsAffected > 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("%s already existis", createGatewayWithAWSReq.GatewayName)
	}

	gateway := models.Gateway{
		Provider:           createGatewayWithAWSReq.Provider,
		Region:             createGatewayWithAWSReq.Region,
		GatewayName:        formattedGatewayName,
		GatewayType:        createGatewayWithAWSReq.GatewayType,
		RPCURL:             createGatewayWithAWSReq.RPCURL,
		Password:           createGatewayWithAWSReq.Password,
		TranscodingProfile: createGatewayWithAWSReq.TranscodingProfile,
		UserID:             userID,
	}

	err := s.gatewayRepository.CreateGateway(&gateway)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	ec2InstancePayload := types.CreateEC2InstanceReq{
		CredentialsID:     createGatewayWithAWSReq.CredentialsID,
		EC2InstanceTypeID: createGatewayWithAWSReq.EC2InstanceTypeID,
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

	return nil, 0, nil
}

func (s *GatewayService) CreateGateway(gateway *models.Gateway) (int, error) {
	client := s.getGatewayClient()

	gateway.GatewayName = utils.ToKebabCase(gateway.GatewayName)

	payload := types.DeployGatewayPayload{
		RPCURL:             gateway.RPCURL,
		Password:           gateway.Password,
		GatewayType:        gateway.GatewayType,
		GatewayName:        gateway.GatewayName,
		TranscodingProfile: gateway.TranscodingProfile,
		Provider:           gateway.Provider,
	}

	task, err := s.gatewayTask.NewDeployGatewayTask(payload)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	info, err := client.Enqueue(task)
	if err != nil {
		return http.StatusInternalServerError, errors.New("unable to queue task")
	}

	gateway.QueueID = info.ID
	gateway.User = nil

	if err := s.gatewayRepository.CreateGateway(gateway); err != nil {
		return http.StatusInternalServerError, errors.New("unable to create gateway")
	}

	defer client.Close()

	return http.StatusCreated, nil
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
