// Package database sets up postgres database and handles migrations automatically
package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
)

var DB *gorm.DB

func InitDB(conf *config.Config) {
	var err error

	loggerLevel := logger.Silent
	if conf.Environment == "development" {
		loggerLevel = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(conf.GetDSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(loggerLevel),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established")
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
