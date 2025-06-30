package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"gwid.io/gwid-core/internal/types"
)

type deployInstancePayload struct {
	rpcUrl             string
	password           string
	gatewayType        string
	gatewayName        string
	transcodingProfile string
}

const (
	TypeDeployInstance = "deploy:instance"
)

type DeployInstanceTask interface {
	DeployInstanceTask(types.DeployInstancePayloadReq) (*asynq.Task, error)
}

type deployInstanceTask struct{}

func NewDeployInstanceTask() DeployInstanceTask {
	return &deployInstanceTask{}
}

func (d *deployInstanceTask) DeployInstanceTask(payload types.DeployInstancePayloadReq) (*asynq.Task, error) {
	log.Println("adding deployment to queue")

	deploymentPayload, err := json.Marshal(deployInstancePayload{rpcUrl: payload.RpcUrl, password: payload.Password, gatewayType: payload.GatewayType, gatewayName: payload.GatewayName, transcodingProfile: payload.TranscodingProfile})

	if err != nil {
		return nil, err
	}

	queue := asynq.NewTask(TypeDeployInstance, deploymentPayload, asynq.MaxRetry(5), asynq.Timeout(5*time.Minute))

	return queue, nil
}

func (d *deployInstancePayload) HandleDeployInstanceTask(ctx context.Context, task *asynq.Task) error {
	var payload deployInstancePayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarsal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Println("deploying gateway instance")

	return nil
}
