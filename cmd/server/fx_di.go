package server

import (
	"go.uber.org/fx"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/controllers"
	"gwid.io/gwid-core/internal/cron"
	"gwid.io/gwid-core/internal/database"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/router"
	"gwid.io/gwid-core/internal/services"
)

func RunGwidCore() {
	fx.New(
		fx.Provide(
			config.NewConfig,

			database.NewDatabase,

			repositories.NewUserRepository,
			repositories.NewGatewayRepository,
			repositories.NewAWSCredentialsRepository,
			repositories.NewEC2Repository,
			repositories.NewReferralRewardRepository,

			services.NewAuthService,
			services.NewJwtService,
			services.NewUserService,
			services.NewGatewayService,
			services.NewRegionService,
			services.NewAWSCredentialsService,
			services.NewEncryptionService,
			services.NewEC2Service,
			services.NewGatewayTaskService,
			services.NewReferralRewardService,

			controllers.NewAuthController,
			controllers.NewUserController,
			controllers.NewGatewayController,
			controllers.NewRegionController,
			controllers.NewAWSCredentialsController,
			controllers.NewEC2Controller,

			cron.NewCronService,
			cron.NewEC2Cron,

			router.NewRouter,
			NewGinServer,
		),
		fx.Invoke(RunServer, RunQueueServer),
	).Run()
}
