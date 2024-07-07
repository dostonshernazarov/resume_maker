package services

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/resume_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/grpc_service_clients"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/usecase"
	"go.uber.org/zap"
)

type resumeRPC struct {
	logger        *zap.Logger
	resumeUseCase usecase.Resume
	client        grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger, resumeUseCase usecase.Resume, client *grpc_service_clients.ServiceClients) resume_service.ResumeServiceServer {
	return &resumeRPC{
		logger:        logger,
		resumeUseCase: resumeUseCase,
		client:        *client,
	}
}

func (s resumeRPC) CreateResume(ctx context.Context, in *resume_service.Resume) (*resume_service.ResumeWithID, error) {
	resume, err := s.resumeUseCase.CreateResume(ctx, &entity.Resume{
		ID:          in.Id,
		UserID:      in.UserId,
		URL:         in.Url,
		Salary:      int64(in.Salary),
		JobTitle:    in.JobTitle,
		Region:      in.Region,
		JobLocation: in.JobLocation,
		JobType:     in.JobType,
		Experience:  in.Experience,
		Template:    in.Template,
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &resume_service.ResumeWithID{
		ResumeId: resume.ID,
	}, nil
}

func (s resumeRPC) UpdateResume(ctx context.Context, in *resume_service.Resume) (*resume_service.Resume, error) {
	resume, err := s.resumeUseCase.UpdateResume(ctx, &entity.Resume{
		ID:          in.Id,
		UserID:      in.UserId,
		URL:         in.Url,
		Salary:      int64(in.Salary),
		JobTitle:    in.JobTitle,
		Region:      in.Region,
		JobLocation: in.JobLocation,
		JobType:     in.JobType,
		Experience:  in.Experience,
		Template:    in.Template,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &resume_service.Resume{
		Id:          resume.ID,
		UserId:      resume.UserID,
		Url:         resume.URL,
		Salary:      uint64(resume.Salary),
		JobTitle:    resume.JobTitle,
		Region:      resume.Region,
		JobLocation: resume.JobLocation,
		JobType:     resume.JobType,
		Experience:  resume.Experience,
		Template:    resume.Template,
	}, nil
}

func (s resumeRPC) DeleteResume(ctx context.Context, in *resume_service.ResumeWithID) (*resume_service.Status, error) {
	if err := s.resumeUseCase.DeleteResume(ctx, in.ResumeId); err != nil {
		s.logger.Error(err.Error())
		return &resume_service.Status{Action: false}, err
	}

	return &resume_service.Status{Action: true}, nil
}

func (s resumeRPC) DeleteUserResume(ctx context.Context, in *resume_service.UserWithID) (*resume_service.Status, error) {
	if err := s.resumeUseCase.DeleteUserResumes(ctx, in.UserId); err != nil {
		s.logger.Error(err.Error())
		return &resume_service.Status{Action: false}, err
	}

	return &resume_service.Status{Action: true}, nil
}

func (s resumeRPC) GetResumeByID(ctx context.Context, in *resume_service.ResumeWithID) (*resume_service.Resume, error) {
	resume, err := s.resumeUseCase.GetResumeByID(ctx, in.ResumeId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &resume_service.Resume{
		Id:          resume.ID,
		UserId:      resume.UserID,
		Url:         resume.URL,
		Salary:      uint64(resume.Salary),
		JobTitle:    resume.JobTitle,
		Region:      resume.Region,
		JobLocation: resume.JobLocation,
		JobType:     resume.JobType,
		Experience:  resume.Experience,
		Template:    resume.Template,
	}, nil
}

func (s resumeRPC) GetUserResume(ctx context.Context, in *resume_service.UserWithID) (*resume_service.ListResumeResponse, error) {
	offset := in.Limit * (in.Page - 1)
	resumes, err := s.resumeUseCase.GetUserResume(ctx, in.UserId, in.Limit, offset)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := &resume_service.ListResumeResponse{}

	for _, resume := range resumes.Resumes {
		response.Resumes = append(response.Resumes, &resume_service.Resume{
			Id:          resume.ID,
			UserId:      resume.UserID,
			Url:         resume.URL,
			Salary:      uint64(resume.Salary),
			JobTitle:    resume.JobTitle,
			Region:      resume.Region,
			JobLocation: resume.JobLocation,
			JobType:     resume.JobType,
			Experience:  resume.Experience,
			Template:    resume.Template,
		})
	}
	response.TotalCount = resumes.TotalCount

	return response, nil
}

func (s resumeRPC) ListResume(ctx context.Context, in *resume_service.ListRequest) (*resume_service.ListResumeResponse, error) {
	resumes, err := s.resumeUseCase.ListResume(ctx, &entity.ListRequest{
		Page:        int64(in.Page),
		Limit:       int64(in.Limit),
		JobTitle:    in.JobTitle,
		JobLocation: in.JobLocation,
		JobType:     in.JobType,
		Salary:      in.Salary,
		Region:      in.Region,
		Experience:  in.Experience,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := &resume_service.ListResumeResponse{}

	for _, resume := range resumes.Resumes {

		response.Resumes = append(response.Resumes, &resume_service.Resume{
			Id:          resume.ID,
			UserId:      resume.UserID,
			Url:         resume.URL,
			Salary:      uint64(resume.Salary),
			JobTitle:    resume.JobTitle,
			Region:      resume.Region,
			JobLocation: resume.JobLocation,
			JobType:     resume.JobType,
			Experience:  resume.Experience,
			Template:    resume.Template,
		})
	}
	response.TotalCount = resumes.TotalCount

	return response, nil
}
