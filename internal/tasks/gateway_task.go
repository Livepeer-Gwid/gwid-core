package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
)

type GatewayTask struct{}

func NewGatewayTask() *GatewayTask {
	return &GatewayTask{}
}

func (gt *GatewayTask) NewDeployGatewayTask(payload types.DeployGatewayPayload) (*asynq.Task, error) {
	deploymentPayload, err := json.Marshal(types.DeployGatewayPayload{RPCURL: payload.RPCURL, Password: payload.Password, GatewayType: payload.GatewayType, GatewayName: payload.GatewayName, TranscodingProfile: payload.TranscodingProfile})
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(utils.TypeDeployGateway, deploymentPayload, asynq.MaxRetry(5), asynq.Timeout(5*time.Minute))

	return task, nil
}

func (gt *GatewayTask) HandleDeployGatewayTask(ctx context.Context, task *asynq.Task) error {
	var payload types.DeployGatewayPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarsal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Println("processing task", task.ResultWriter().TaskID())

	// invokeScript := "./script/init.sh"
	//
	// cmd := exec.Command(invokeScript, payload.GatewayName, payload.GatewayType, payload.RpcUrl, payload.Password, payload.TranscodingProfile)
	//
	// var stdout, stderr bytes.Buffer
	//
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	//
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatalf("Command failed with error: %v, Stderr: %s", err, stderr.String())
	// }
	//
	// // taskId := task.ResultWriter().TaskID()
	//
	// fmt.Printf("Stdout: %s\n", stdout.String())

	return nil
}
