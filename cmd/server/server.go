package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gwid.io/gwid-core/internal/config"
)

func NewGinServer(router *gin.Engine, cfg *config.Config) *http.Server {
	server := &http.Server{
		Addr:           cfg.GetServerAddress(),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server
}

func RunServer(lc fx.Lifecycle, server *http.Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting server on port %s", cfg.Port)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server")
			return server.Shutdown(ctx)
		},
	})
}
