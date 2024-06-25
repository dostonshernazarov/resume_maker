package grpc_service_clients

import (
	"fmt"

	pbr "github.com/dostonshernazarov/resume_maker/api-service/genproto/resume_service"
	pbu "github.com/dostonshernazarov/resume_maker/api-service/genproto/user_service"

	"google.golang.org/grpc"

	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
)

type ServiceClient interface {
	ResumeService() pbr.ResumeServiceClient
	UserService() pbu.UserServiceClient
	Close()
}

type serviceClient struct {
	connections   []*grpc.ClientConn
	resumeService pbr.ResumeServiceClient
	userService   pbu.UserServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	// dial to resume service
	connResumeService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ResumeService.Host, cfg.ResumeService.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// user service

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClient{
		resumeService: pbr.NewResumeServiceClient(connResumeService),
		userService:   pbu.NewUserServiceClient(connUserService),
		connections: []*grpc.ClientConn{
			connResumeService,
			connUserService,
		},
	}, nil
}

func (s *serviceClient) ResumeService() pbr.ResumeServiceClient {
	return s.resumeService
}

func (s *serviceClient) UserService() pbu.UserServiceClient {
	return s.userService
}

func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
