package server

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/tasks"
	"gwid.io/gwid-core/internal/utils"
)

func RunQueueServer(lc fx.Lifecycle, cfg *config.Config, gatewayTask *tasks.GatewayTask) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddress},
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
	mux.HandleFunc(utils.TypeDeployGateway, gatewayTask.HandleDeployGatewayTask)

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
