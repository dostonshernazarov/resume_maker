package config

import (
	"os"
	"strings"
	"time"
)

type webAddress struct {
	Host string
	Port string
}

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	APIToken    string
	ChatID      string
	Server      struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
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
		AccessTTL  time.Duration
		RefreshTTL time.Duration
		SignInKey  string
	}
	Minio struct {
		Host       string
		Port       string
		AccessKey  string
		SecretKey  string
		Location   string
		BucketName string
	}
	Kafka struct {
		Address []string
		Topic   struct {
			UserCreateTopic string
		}
	}
	RabbitMQ struct {
		Host  string
		Topic string
	}
	ResumeService   webAddress
	UserService     webAddress
	TelegramService webAddress
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "7s")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "cvmaker_api-service")
	config.Server.Port = getEnv("SERVER_PORT", ":9070")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "cvmaker_redis-db")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	// minio configuration
	config.Minio.Host = getEnv("MINIO_HOST", "35.198.173.128")
	config.Minio.Port = getEnv("MINIO_PORT", ":9000")
	config.Minio.AccessKey = getEnv("MINIO_ACCESS_KEY", "minioadmin")
	config.Minio.SecretKey = getEnv("MINIO_SECRET_KEY", "minioadmin")
	config.Minio.BucketName = getEnv("MINIO_BUCKET_NAME", "resumes")

	config.ResumeService.Host = getEnv("RESUME_SERVICE_GRPC_HOST", "cvmaker_resume-service")
	config.ResumeService.Port = getEnv("RESUME_SERVICE_GRPC_PORT", ":9070")

	// user configuration
	config.UserService.Host = getEnv("USER_SERVICE_GRPC_HOST", "cvmaker_user-service")
	config.UserService.Port = getEnv("USER_SERVICE_GRPC_PORT", ":9090")

	// telegram configuration
	config.TelegramService.Host = getEnv("TELEGRAM_SERVICE_GRPC_HOST", "localhost")
	config.TelegramService.Port = getEnv("TELEGRAM_SERVICE_GRPC_PORT", ":8090")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "6h"))
	if err != nil {
		return nil, err
	}

	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "168h"))
	if err != nil {
		return nil, err
	}

	//JWT key
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNING_KEY", "token_secret")

	//Telegram config
	config.APIToken = getEnv("TELEGRAM_BOT_TOKEN", "7303220559:AAHgpp6y1f_dk-iLsZ_gGrjwoI5-9mTVrPY")
	config.ChatID = getEnv("CHAT_ID", "-1002142909351")

	//RabbitMQ config
	config.RabbitMQ.Host = getEnv("AMQP_SERVER", "amqp://guest:guest@rabbitmq:5672/")
	config.RabbitMQ.Topic = getEnv("QUEUE_NAME", "cvmaker_queue")

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
