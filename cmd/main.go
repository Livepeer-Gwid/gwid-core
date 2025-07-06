package main

import (
	"go.uber.org/fx"
	"gwid.io/gwid-core/cmd/server"
	"gwid.io/gwid-core/internals/config"
	"gwid.io/gwid-core/internals/database"
	"gwid.io/gwid-core/internals/router"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfig,
			database.NewDatabase,
			router.NewRouter,
			server.NewGinServer,
		),
		fx.Invoke(server.RunServer),
	).Run()
}
