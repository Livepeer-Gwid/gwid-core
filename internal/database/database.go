// Package database sets up database and models
package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	loggerLevel := logger.Silent
	if cfg.Environment == "development" {
		loggerLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(loggerLevel),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Gateway{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")

	return db
}
