package config

import (
	"os"
	"strings"
	"time"
)

const (
	OtpSecret = "some_secret"
)

type webAddress struct {
	Host string
	Port string
}

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	Server      struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}
	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}
	Context struct {
		Timeout string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		Name     string
	}
	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
		SignInKey  string
	}
	Minio struct {
		Endpoint              string
		AccessKey             string
		SecretKey             string
		Location              string
		MovieUploadBucketName string
	}
	Kafka struct {
		Address []string
		Topic   struct {
			UserCreateTopic string
		}
	}
	EstablishmentService webAddress
	UserService          webAddress
	BookingService       webAddress
	OTLPCollector        webAddress
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "7s")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "api-service")
	config.Server.Port = getEnv("SERVER_PORT", ":8080")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "postgres")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "touristandb")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "doston")

	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "redis-db")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	config.EstablishmentService.Host = getEnv("ESTABLISHMENT_SERVICE_GRPC_HOST", "establishment-service")
	config.EstablishmentService.Port = getEnv("ESTABLISHMENT_SERVICE_GRPC_PORT", ":50024")

	// user configuration
	config.UserService.Host = getEnv("USER_SERVICE_GRPC_HOST", "user-service")
	config.UserService.Port = getEnv("USER_SERVICE_GRPC_PORT", ":50025")

	// booking configuration
	config.BookingService.Host = getEnv("BOOKING_SERVICE_GRPC_HOST", "booking-service")
	config.BookingService.Port = getEnv("BOOKING_SERVICE_GRPC_PORT", ":50023")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "2h"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "48h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNIN_KEY", "debug_booking")

	// otlp collector configuration
	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:9092"), ",")
	config.Kafka.Topic.UserCreateTopic = getEnv("KAFKA_USER_CREATE", "user.create.api")

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
