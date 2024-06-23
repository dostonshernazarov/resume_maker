package grpc_service_clients

import (
	"fmt"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/resume_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/config"

	"google.golang.org/grpc"
)

type ServiceClients interface {
	ResumeService() resume_service.ResumeServiceClient
	Close()
}

type serviceClients struct {
	resumeService resume_service.ResumeServiceClient
	services      []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {
	resumeServiceConnection, err := grpc.Dial(
		fmt.Sprintf("%s%s", config.ResumeService.Host, config.ResumeService.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClients{
		resumeService: resume_service.NewResumeServiceClient(resumeServiceConnection),
		services: []*grpc.ClientConn{
			resumeServiceConnection,
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

func (s *serviceClients) ResumeService() resume_service.ResumeServiceClient {
	return s.resumeService
}
