package server

import (
	"go.uber.org/fx"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/controllers"
	"gwid.io/gwid-core/internal/database"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/router"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/tasks"
)

func RunGwidCore() {
	fx.New(
		fx.Provide(
			config.NewConfig,

			database.NewDatabase,

			repositories.NewUserRepository,
			repositories.NewGatewayRepository,

			services.NewAuthService,
			services.NewJwtService,
			services.NewUserService,
			services.NewGatewayService,

			controllers.NewAuthController,
			controllers.NewUserController,
			controllers.NewGatewayController,

			tasks.NewGatewayTask,

			router.NewRouter,
			NewGinServer,
		),
		fx.Invoke(RunServer, RunQueueServer),
	).Run()
}
