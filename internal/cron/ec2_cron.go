package cron

import (
	"context"
	"fmt"
	"log"

	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Cron struct {
	ec2Repository *repositories.EC2Repository
	config        *config.Config
}

func NewEC2Cron(
	ec2Repository *repositories.EC2Repository,
	config *config.Config,
) *EC2Cron {
	return &EC2Cron{
		ec2Repository: ec2Repository,
		config:        config,
	}
}

func (s *EC2Cron) SyncEC2Instances() {
	ec2Count, err := s.ec2Repository.GetEC2TotalCount()

	if err != nil {
		return
	}

	if ec2Count > 0 {
		return
	}

	accessID := s.config.AwsAccessID
	secretAccessKey := s.config.AwsSecretAccessKey

	ctx := context.TODO()

	creds := credentials.NewStaticCredentialsProvider(accessID, secretAccessKey, "")

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(creds),
		awsConfig.WithRegion("eu-central-1"),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	paginator := ec2.NewDescribeInstanceTypesPaginator(ec2Client, &ec2.DescribeInstanceTypesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("processor-info.supported-architecture"),
				Values: []string{"x86_64"},
			},
		},
	})

	fmt.Println("Syncing EC2 Instance Types:")

	var ec2Instances []models.EC2

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to get instance types: %v", err)
		}

		for _, it := range page.InstanceTypes {
			vcpu := it.VCpuInfo.DefaultVCpus
			memoryGiB := float64(*it.MemoryInfo.SizeInMiB) / 1024.0

			if vcpu != nil && *vcpu >= 4 && memoryGiB >= 8.0 {
				ec2Instance := models.EC2{
					Tag:             string(it.InstanceType),
					Ram:             memoryGiB,
					Cpu:             int(*vcpu),
					Architecture:    string(it.ProcessorInfo.SupportedArchitectures[0]),
					CpuManufacturer: *it.ProcessorInfo.Manufacturer,
				}

				ec2Instances = append(ec2Instances, ec2Instance)

			}
		}
	}

	if err := s.ec2Repository.CreateEC2InstanceTypes(&ec2Instances); err != nil {
		log.Println(err)

		return
	}

	log.Println("EC2 Syncing done")
}
