// Package di mounts all services
package di

import (
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/database"
)

type Container struct {
	DB *gorm.DB
}

func NewContainer(conf *config.Config) *Container {
	database.InitDB(conf)

	db := database.DB

	return &Container{
		DB: db,
	}
}
