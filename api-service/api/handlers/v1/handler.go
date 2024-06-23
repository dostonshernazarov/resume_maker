package v1

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	grpcClients "github.com/dostonshernazarov/resume_maker/api-service/internal/infrastructure/grpc_service_client"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"

	appV "github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/app_version"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/event"
	// "github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/refresh_token"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	BrokerProducer event.BrokerProducer
	Enforcer       *casbin.Enforcer
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	BrokerProducer event.BrokerProducer
	Enforcer       *casbin.Enforcer
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		Service:        c.Service,
		JwtHandler:     c.JwtHandler,
		AppVersion:     c.AppVersion,
		BrokerProducer: c.BrokerProducer,
		Enforcer:       c.Enforcer,
	}
}
