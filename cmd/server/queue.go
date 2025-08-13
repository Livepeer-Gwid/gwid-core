package server

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/utils"
)

func RunQueueServer(lc fx.Lifecycle, cfg *config.Config, gatewayTaskService *services.GatewayTaskService) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddress, Password: cfg.RedisPassword},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(utils.TypeDeployAWSGateway, gatewayTaskService.HandleAWSDeployGatewayTask)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Starting asynq queue server")
				if err := srv.Run(mux); err != nil {
					log.Fatalf("could not run mux server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Shutdown()
			return nil
		},
	})
}
