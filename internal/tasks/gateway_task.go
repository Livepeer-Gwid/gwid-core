// Package tasks
package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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

	task := asynq.NewTask(utils.TypeDeployGateway, deploymentPayload, asynq.MaxRetry(1), asynq.Timeout(5*time.Minute))

	return task, nil
}

func (gt *GatewayTask) HandleDeployGatewayTask(ctx context.Context, task *asynq.Task) error {
	var payload types.DeployGatewayPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarsal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Println("processing task", task.ResultWriter().TaskID())

	invokeScript := "./scripts/invoke.sh"

	cmd := exec.Command(invokeScript, payload.GatewayName, payload.GatewayType, payload.RPCURL, payload.Password, payload.TranscodingProfile)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command failed with error: %v, Stderr: %s", err, stderr.String())

		log.Println(stdout.String())

		return nil
	}

	// log.Println(stdout.String())

	// taskId := task.ResultWriter().TaskID()

	// var successResponse map[string]any
	// var errorResponse map[string]any

	// if err := json.Unmarshal(stderr.Bytes(), &errorResponse); err != nil {
	//
	// 	fmt.Println("Error unmarshaling JSON:", err)
	//
	// 	return errors.New("error unmarshaling JSON")
	// }
	//
	// if err := json.Unmarshal(stdout.Bytes(), &successResponse); err != nil {
	// 	fmt.Println("Error unmarshaling JSON:", err)
	//
	// 	return errors.New("error unmarshaling JSON")
	// }
	//
	// // fmt.Printf("Stdout: %s\n", stdout.String())
	//
	// fmt.Println("json", successResponse["message"])

	return nil
}
