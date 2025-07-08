package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment        string
	Port               string
	GinMode            string
	JwtSecret          string
	RedisAddress       string
	DevPostgresConfig  PostgresConfig
	ProdPostgresConfig PostgresConfig
}

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	env := &Config{
		Environment:  GetEnv("ENVIRONMENT", "development"),
		Port:         GetEnv("PORT", "5000"),
		GinMode:      GetEnv("GIN_MODE", "debug"),
		JwtSecret:    GetEnv("JWT_SECRET", "the-fallback-key"),
		RedisAddress: GetEnv("REDIS_ADDRESS", ""),
		DevPostgresConfig: PostgresConfig{
			Host:         GetEnv("DEV_DB_HOST", "localhost"),
			Port:         GetEnv("DEV_DB_PORT", "5432"),
			User:         GetEnv("DEV_DB_USER", "postgres"),
			Password:     GetEnv("DEV_DB_PASSWORD", "postgres"),
			DatabaseName: GetEnv("DEV_DB_NAME", "gwid"),
			SSLMode:      GetEnv("DEV_DB_SSLMODE", "disable"),
		},
	}

	return env
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.Port)
}

func (c *Config) GetDSN() string {
	environment := c.Environment

	var dsn string

	if environment == "development" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.DevPostgresConfig.Host, c.DevPostgresConfig.Port, c.DevPostgresConfig.User, c.DevPostgresConfig.Password, c.DevPostgresConfig.DatabaseName, c.DevPostgresConfig.SSLMode)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.ProdPostgresConfig.Host, c.ProdPostgresConfig.Port, c.ProdPostgresConfig.User, c.ProdPostgresConfig.Password, c.ProdPostgresConfig.DatabaseName, c.ProdPostgresConfig.SSLMode)
	}

	return dsn
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return fallback
}
