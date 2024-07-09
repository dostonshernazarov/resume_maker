package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"github.com/dostonshernazarov/resume_maker/api-service/api"
	grpcService "github.com/dostonshernazarov/resume_maker/api-service/internal/infrastructure/grpc_service_client"

	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/infrastructure/rabbitmq"

	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/redis"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/app_version"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/event"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	RedisDB        *redis.RedisDB
	server         *http.Server
	Enforcer       *casbin.Enforcer
	Clients        grpcService.ServiceClient
	BrokerProducer event.BrokerProducer
	appVersion     app_version.AppVersion
	writer         *rabbitmq.RabbitMQProducerImpl
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// kafka producer init
	// kafkaProducer := kafka.NewProducer(&cfg, logger)

	//RabbitMQ producer init
	writer, err := rabbitmq.NewRabbitMQProducer(cfg.RabbitMQ.Host)
	if err != nil {
		return nil, err
	}

	// err = writer.ProducerMessage("test-topic", []byte("\nthis message has come from produce"))
	// if err != nil {
	// 	return nil, err
	// }

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}

	// initialization enforcer
	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		return nil, err
	}

	// enforcer.SetCache(policy.NewCache(&redisdb.Client))

	//var (
	//	contextTimeout time.Duration
	//)

	// context timeout initialization
	_, err = time.ParseDuration(cfg.Context.Timeout)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:   &cfg,
		Logger:   logger,
		RedisDB:  redisdb,
		Enforcer: enforcer,
		writer:   writer,
		// BrokerProducer: kafkaProducer,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Enforcer:       a.Enforcer,
		Service:        clients,
		BrokerProducer: a.BrokerProducer,
		AppVersion:     a.appVersion,
		Writer:         a.writer,
	})
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)

	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	// close grpc connections
	a.Clients.Close()

	// kafka producer close
	a.BrokerProducer.Close()

	//rabbitmq producer close
	a.writer.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
