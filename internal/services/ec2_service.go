package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
	"gwid.io/gwid-core/internal/utils"
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

func (s *EC2Service) GetUbuntuImageID(ec2Client *ec2.Client) (string, error) {
	input := &ec2.DescribeImagesInput{
		Owners: []string{"099720109477"},
		Filters: []awsTypes.Filter{
			{
				Name:   aws.String("architecture"),
				Values: []string{"x86_64"},
			},
		},
	}

	result, err := ec2Client.DescribeImages(context.TODO(), input)
	if err != nil {
		return "", err
	}

	if len(result.Images) == 0 {
		return "", errors.New("no image found")
	}

	sort.Slice(result.Images, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, *result.Images[i].CreationDate)
		timeJ, _ := time.Parse(time.RFC3339, *result.Images[j].CreationDate)
		return timeI.After(timeJ)
	})

	latestImage := result.Images[0]

	return *latestImage.ImageId, nil
}

func (s *EC2Service) CreateEC2Instance(ec2InstanceReq types.CreateEC2InstanceReq, userID uuid.UUID) (string, int, error) {
	userCreds, int, err := s.awsCredentialsService.GetAWSCredentialsByID(ec2InstanceReq.CredentialsID, userID)

	if err != nil {
		return "", int, err
	}

	ec2InstanceType, result := s.ec2Repository.GetEC2InstanceTypeByID(ec2InstanceReq.EC2InstanceTypeID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", http.StatusNotFound, errors.New("ec2 intance type not found")
	}

	creds := credentials.NewStaticCredentialsProvider(userCreds.AccessKeyID, userCreds.SecretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion("eu-central-1"),
	)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("unable to load AWS config")
	}

	ec2Client := ec2.NewFromConfig(cfg)

	imageID, err := s.GetUbuntuImageID(ec2Client)

	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	instanceResult, err := ec2Client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		ImageId:      aws.String(imageID),
		InstanceType: awsTypes.InstanceType(ec2InstanceType.Tag),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		TagSpecifications: []awsTypes.TagSpecification{
			{
				ResourceType: awsTypes.ResourceTypeInstance,
				Tags: []awsTypes.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(utils.ToKebabCase(ec2InstanceReq.InstanceName)),
					},
				},
			},
		},
	})

	if err != nil {
		return "", http.StatusBadRequest, err
	}

	if len(instanceResult.Instances) == 0 {
		return "", http.StatusInternalServerError, errors.New("no instance was created")
	}

	for _, instance := range instanceResult.Instances {
		fmt.Printf("Created instance with ID: %s\n", *instance.InstanceId)
	}

	return *instanceResult.Instances[0].InstanceId, 200, nil
}

func (s *EC2Service) WaitForInstanceRunning(instanceID string, ctx context.Context, ec2Client *ec2.Client) error {
	waiter := ec2.NewInstanceRunningWaiter(ec2Client)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	err := waiter.Wait(ctx, input, 4*time.Minute)

	if err != nil {
		return err
	}

	result, err := ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to describe EC2 instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	instance := result.Reservations[0].Instances[0]

	if instance.IamInstanceProfile != nil {
		return nil
	} else {
		return fmt.Errorf("no IAM instance profile attached")
	}
}

func (s *EC2Service) RunCommand(instanceID string, ctx context.Context, command string, ssmClient *ssm.Client) (string, error) {
	sendInput := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]string{
			"commands": {*aws.String(command)},
		},
		Comment: aws.String("Command executed via Go SDK"),
	}

	sendOutput, err := ssmClient.SendCommand(ctx, sendInput)
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	commandID := *sendOutput.Command.CommandId

	return commandID, nil
}

func (s *EC2Service) WaitForCommandCompletion(instanceID string, ctx context.Context, commandID string, ssmClient *ssm.Client) (*types.CommandResult, error) {
	startTime := time.Now()
	maxWaitTime := 5 * time.Minute
	pollInterval := 2 * time.Second

	for {
		if time.Since(startTime) > maxWaitTime {
			return nil, fmt.Errorf("command execution timed out after %v", maxWaitTime)
		}

		getInput := &ssm.GetCommandInvocationInput{
			CommandId:  aws.String(commandID),
			InstanceId: aws.String(instanceID),
		}

		output, err := ssmClient.GetCommandInvocation(ctx, getInput)

		if err != nil {
			log.Printf("waiting for command to start... %v", err)
			time.Sleep(pollInterval)
			continue
		}

		status := output.Status

		log.Printf("Command status: %s", status)

		if status == ssmTypes.CommandInvocationStatusInProgress ||
			status == ssmTypes.CommandInvocationStatusPending {
			time.Sleep(pollInterval)
			continue
		}

		result := &types.CommandResult{
			CommandID:     commandID,
			Status:        string(status),
			ExecutionTime: time.Since(startTime),
			ExitCode:      utils.SafeInt32Value(&output.ResponseCode),
			StandardOut:   utils.SafeStringValue(output.StandardOutputContent),
			StandardErr:   utils.SafeStringValue(output.StandardErrorContent),
		}

		return result, nil
	}
}
