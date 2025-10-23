package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hibiken/asynq"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
)

type GatewayTaskService struct {
	awsCredentialsService *AWSCredentialsService
	ec2Service            *EC2Service
}

func NewGatewayTaskService(awsCredentialsService *AWSCredentialsService, ec2Service *EC2Service) *GatewayTaskService {
	return &GatewayTaskService{
		awsCredentialsService: awsCredentialsService,
		ec2Service:            ec2Service,
	}
}

func (gt *GatewayTaskService) NewAWSDeployGatewayTask(payload types.DeployAWSGatewayPayload) (*asynq.Task, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(utils.TypeDeployAWSGateway, payloadJson, asynq.MaxRetry(2), asynq.Timeout(5*time.Minute))

	return task, nil
}

func (gt *GatewayTaskService) HandleAWSDeployGatewayTask(ctx context.Context, task *asynq.Task) error {
	var payload types.DeployAWSGatewayPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarsal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Println("processing task", task.ResultWriter().TaskID())

	userCreds, _, err := gt.awsCredentialsService.GetAWSCredentialsByID(payload.CredentialsID, payload.UserID)
	if err != nil {
		return fmt.Errorf("unable to get user AWS credentials: %v: %w", err, asynq.SkipRetry)
	}

	creds := credentials.NewStaticCredentialsProvider(userCreds.AccessKeyID, userCreds.SecretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(payload.Region),
	)
	if err != nil {
		return fmt.Errorf("unable to load AWS config: %v: %w", err, asynq.SkipRetry)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	ssmClient := ssm.NewFromConfig(cfg)

	if err := gt.ec2Service.WaitForInstanceRunning(payload.InstanceID, ctx, ec2Client); err != nil {
		return fmt.Errorf("unable to get instance running state: %v: %w", err, asynq.SkipRetry)
	}

	if err := gt.ec2Service.WaitForSSM(payload.InstanceID, ctx, ssmClient); err != nil {
		return fmt.Errorf("%v: %w", err, asynq.SkipRetry)
	}

	command := `echo "{\"success\": true, \"grafana_url\": \"https://gwid.io\"}"`

	// command := "ls -la /home && echo 'Hello from EC2!' && date"

	commandID, err := gt.ec2Service.RunCommand(payload.InstanceID, ctx, command, ssmClient)
	if err != nil {
		return fmt.Errorf("something went wrong: %v: %w", err, asynq.SkipRetry)
	}

	commandResult, err := gt.ec2Service.WaitForCommandCompletion(payload.InstanceID, ctx, commandID, ssmClient)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("something went wrong with command execution: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("\n=== Command Execution Result ===\n")
	fmt.Printf("Command ID: %s\n", commandResult.CommandID)
	fmt.Printf("Status: %s\n", commandResult.Status)
	fmt.Printf("Exit Code: %d\n", commandResult.ExitCode)
	fmt.Printf("Execution Time: %v\n", commandResult.ExecutionTime)

	if commandResult.StandardOut != "" {
		fmt.Printf("\n--- Standard Output ---\n%s\n", commandResult.StandardOut)
	}

	if commandResult.StandardErr != "" {
		fmt.Printf("\n--- Standard Error ---\n%s\n", commandResult.StandardErr)
	}
	fmt.Printf("===============================\n")

	return nil
}
