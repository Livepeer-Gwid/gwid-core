package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	GinMode     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	env := &Config{
		Environment: GetEnv("ENVIRONMENT", "development"),
		Port:        GetEnv("PORT", "5000"),
		GinMode:     GetEnv("GIN_MODE", "debug"),
	}

	return env, nil
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.Port)
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
