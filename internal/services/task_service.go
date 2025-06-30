package services

import "gwid.io/gwid-core/internal/config"

type TaskService interface {
}

type taskService struct {
	configService *config.Config
}
