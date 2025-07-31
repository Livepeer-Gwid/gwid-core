package cron

import (
	"github.com/robfig/cron/v3"
)

type CronService struct {
	ec2Cron *EC2Cron
}

func NewCronService(ec2Cron *EC2Cron) *CronService {
	return &CronService{
		ec2Cron: ec2Cron,
	}
}

func (s *CronService) StartScheduler() *cron.Cron {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("@daily", s.ec2Cron.SyncEC2Instances)

	c.Start()

	return c
}
