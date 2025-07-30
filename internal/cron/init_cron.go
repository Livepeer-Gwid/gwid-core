package cron

import (
	"github.com/robfig/cron/v3"
)

type CronService struct{}

func NewCronService() *CronService {
	return &CronService{}
}

func (s *CronService) StartScheduler() *cron.Cron {
	c := cron.New(cron.WithSeconds())

	c.Start()

	return c
}
