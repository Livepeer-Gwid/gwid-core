package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/di"
	"gwid.io/gwid-core/internal/router"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration %v", err)
	}

	gin.SetMode(conf.GinMode)

	container := di.NewContainer(conf)

	router := router.SetupRoutes(container)

	log.Printf("Server running on %s\n", conf.GetServerAddress())

	if err := router.Run(conf.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
