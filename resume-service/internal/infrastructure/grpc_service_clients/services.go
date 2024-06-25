package grpc_service_clients

import (
	"fmt"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/user_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/config"

	"google.golang.org/grpc"
)

type ServiceClients interface {
	UserService() user_service.UserServiceClient
	Close()
}

type serviceClients struct {
	userService user_service.UserServiceClient
	services    []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {
	userServiceConnection, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.UserService.Host, config.UserService.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClients{
		userService: user_service.NewUserServiceClient(userServiceConnection),
		services: []*grpc.ClientConn{
			userServiceConnection,
		},
	}, nil
}

func (s *serviceClients) Close() {
	// closing store-management service
	for _, conn := range s.services {
		err := conn.Close()
		if err != nil {
			return
		}
	}
}

func (s *serviceClients) UserService() user_service.UserServiceClient {
	return s.userService
}
