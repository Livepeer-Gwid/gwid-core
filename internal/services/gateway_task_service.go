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
	cloudflareService     *CloudflareService
}

func NewGatewayTaskService(awsCredentialsService *AWSCredentialsService, ec2Service *EC2Service, cloudflareService *CloudflareService) *GatewayTaskService {
	return &GatewayTaskService{
		awsCredentialsService: awsCredentialsService,
		ec2Service:            ec2Service,
		cloudflareService:     cloudflareService,
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

	ipAddress, err := gt.ec2Service.GetEC2IPAddress(payload.InstanceID, ctx, ec2Client)
	if err != nil {
		return fmt.Errorf("unable to get ec2 IP address", err, asynq.SkipRetry)
	}

	result, err := gt.cloudflareService.AddGatewayToCloudflare(ipAddress, payload.GatewayName)
	if err != nil {
		return fmt.Errorf("unable to add gateway to cloudflare", err, asynq.SkipRetry)
	}

	fmt.Println(result)

	if err := gt.ec2Service.WaitForSSM(payload.InstanceID, ctx, ssmClient); err != nil {
		return fmt.Errorf("%v: %w", err, asynq.SkipRetry)
	}

	generatedEncryptionKey, err := utils.GenerateKey32Bytes()
	if err != nil {
		return fmt.Errorf("unable to generate encryptionKey", err, asynq.SkipRetry)
	}

	subdomain := fmt.Sprintf("%s.gwid-gateway.link", payload.GatewayName)
	databaseName := "gwid-sub"
	walletAddress := ""
	password := payload.UnhashedPassword
	repoURL := "https://github.com/Livepeer-Gwid/gwid-sub.git"
	encryptionKey := generatedEncryptionKey
	gatewayName := payload.GatewayName

	command := fmt.Sprintf(`
		sudo apt update -y

		sudo apt install -y git

  	git clone "https://github.com/Livepeer-Gwid/gwid-setup-scripts.git"

		cd "gwid-setup-scripts"

		chmod +x *.sh

		./setup-server.sh %s %s %s %s %s %s %s
	`, subdomain, repoURL, encryptionKey, databaseName, gatewayName, walletAddress, password)

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
