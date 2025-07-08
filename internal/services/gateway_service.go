package services

import (
	"errors"
	"net/http"

	"github.com/hibiken/asynq"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/tasks"
	"gwid.io/gwid-core/internal/types"
)

type GatewayService struct {
	cfg               *config.Config
	gatewayTask       *tasks.GatewayTask
	gatewayRepository *repositories.GatewayRepository
}

func NewGatewayService(cfg *config.Config, gatewayTask *tasks.GatewayTask, gatewayRepository *repositories.GatewayRepository) *GatewayService {
	return &GatewayService{
		cfg:               cfg,
		gatewayTask:       gatewayTask,
		gatewayRepository: gatewayRepository,
	}
}

func (gs *GatewayService) getGatewayClient() *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: gs.cfg.RedisAddress})

	return client
}

func (gs *GatewayService) CreateGateway(gateway *models.Gateway) (int, error) {
	client := gs.getGatewayClient()

	payload := types.DeployGatewayPayload{
		RPCURL:             gateway.RPCURL,
		Password:           gateway.Password,
		GatewayType:        gateway.GatewayType,
		GatewayName:        gateway.GatewayName,
		TranscodingProfile: gateway.TranscodingProfile,
		Provider:           gateway.Provider,
	}

	task, err := gs.gatewayTask.NewDeployGatewayTask(payload)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	info, err := client.Enqueue(task)
	if err != nil {
		return http.StatusInternalServerError, errors.New("unable to queue task")
	}

	gateway.QueueID = info.ID
	gateway.User = nil

	if err := gs.gatewayRepository.CreateGateway(gateway); err != nil {
		return http.StatusInternalServerError, errors.New("unable to create gateway")
	}

	defer client.Close()

	return http.StatusCreated, nil
}
