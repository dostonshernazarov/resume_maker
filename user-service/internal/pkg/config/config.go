package config

import (
	"os"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}

	TelegramService struct {
		Host string
		Port string
	}

	ResumeService struct {
		Host string
		Port string
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":9090")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "postgres")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "root")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "resume")

	// telegram service
	config.TelegramService.Host = getEnv("TELEGRAM_SERVICE_RPC_HOST", "telegram-service")
	config.TelegramService.Port = getEnv("TELEGRAM_SERVICE_RPC_PORT", ":8090")

	// resume service
	config.ResumeService.Host = getEnv("RESUME_SERVICE_RPC_HOST", "resume-service")
	config.ResumeService.Port = getEnv("RESUME_SERVICE_RPC_PORT", ":9080")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
