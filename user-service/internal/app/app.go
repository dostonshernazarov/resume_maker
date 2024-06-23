package app

import (
	"fmt"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/user_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/delivery/grpc/server"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/delivery/grpc/services"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/grpc_service_clients"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository/postgresql"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/config"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/postgres"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/usecase"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	GrpcServer     *grpc.Server
	ServiceClients grpc_service_clients.ServiceClients
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	log, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// init db
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	// grpc server init
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcctxtags.StreamServerInterceptor(),
			grpczap.StreamServerInterceptor(log),
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(server.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpcctxtags.UnaryServerInterceptor(),
				grpczap.UnaryServerInterceptor(log),
				grpcrecovery.UnaryServerInterceptor(),
			),
			server.UnaryInterceptorData(log),
		)),
	)

	return &App{
		Config:     cfg,
		Logger:     log,
		DB:         db,
		GrpcServer: grpcServer,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	articleRepo := postgresql.NewUsersRepo(a.DB)

	// useCase initialization
	articleUseCase := usecase.NewUserService(contextTimeout, articleRepo)

	user_service.RegisterUserServiceServer(a.GrpcServer, services.NewRPC(a.Logger, articleUseCase, &a.ServiceClients))
	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}

	return nil
}

func (a *App) Stop() {
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// zap logger sync
	err := a.Logger.Sync()
	if err != nil {
		return
	}
}
