package services

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/user_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/grpc_service_clients"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/usecase"
	"time"

	"go.uber.org/zap"
)

type userRPC struct {
	logger      *zap.Logger
	userUseCase usecase.User
	client      grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger, userUseCase usecase.User, client *grpc_service_clients.ServiceClients) user_service.UserServiceServer {
	return &userRPC{
		logger:      logger,
		userUseCase: userUseCase,
		client:      *client,
	}
}

func (s userRPC) CreateUser(ctx context.Context, in *user_service.User) (*user_service.UserWithGUID, error) {
	user, err := s.userUseCase.Create(ctx, &entity.User{
		GUID:        in.Id,
		Name:        in.Name,
		Image:       in.Image,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Refresh:     in.Refresh,
		Password:    in.Password,
		Role:        in.Role,
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &user_service.UserWithGUID{
		Guid: user.GUID,
	}, nil
}

func (s userRPC) UpdateUser(ctx context.Context, in *user_service.User) (*user_service.User, error) {
	user, err := s.userUseCase.Update(ctx, &entity.User{
		GUID:        in.Id,
		Name:        in.Name,
		Image:       in.Image,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Role:        in.Role,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &user_service.User{
		Id:          user.GUID,
		Name:        user.Name,
		Image:       user.Image,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Refresh:     user.Refresh,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) DeleteUser(ctx context.Context, in *user_service.UserWithGUID) (*user_service.ResponseStatus, error) {
	if err := s.userUseCase.Delete(ctx, in.Guid); err != nil {
		s.logger.Error(err.Error())
		return &user_service.ResponseStatus{Status: false}, err
	}

	return &user_service.ResponseStatus{Status: true}, nil
}

func (s userRPC) GetUser(ctx context.Context, in *user_service.Filter) (*user_service.User, error) {
	user, err := s.userUseCase.Get(ctx, in.Filter)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &user_service.User{
		Id:          user.GUID,
		Name:        user.Name,
		Image:       user.Image,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Refresh:     user.Refresh,
		Password:    user.Password,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetAllUsers(ctx context.Context, in *user_service.ListUserRequest) (*user_service.ListUserResponse, error) {
	offset := in.Limit * (in.Page - 1)
	users, err := s.userUseCase.List(ctx, uint64(in.Limit), uint64(offset))
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := user_service.ListUserResponse{}
	for _, u := range users.Users {
		temp := &user_service.User{
			Id:          u.GUID,
			Name:        u.Name,
			Image:       u.Image,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Refresh:     u.Refresh,
			Role:        u.Role,
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		}

		response.Users = append(response.Users, temp)
	}
	response.TotalCount = users.Total

	return &response, nil
}

func (s userRPC) UniqueEmail(ctx context.Context, in *user_service.IsUnique) (*user_service.ResponseStatus, error) {
	response, err := s.userUseCase.UniqueEmail(ctx, &entity.IsUnique{Email: in.Email})

	if err != nil {
		s.logger.Error(err.Error())
		return &user_service.ResponseStatus{Status: true}, err
	}
	if response.Status {
		return &user_service.ResponseStatus{Status: true}, nil
	}

	return &user_service.ResponseStatus{Status: false}, nil
}

func (s userRPC) UpdateRefresh(ctx context.Context, in *user_service.RefreshRequest) (*user_service.ResponseStatus, error) {
	_, err := s.userUseCase.UpdateRefresh(ctx, &entity.UpdateRefresh{
		UserID:       in.UserId,
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return &user_service.ResponseStatus{Status: false}, err
	}

	return &user_service.ResponseStatus{Status: true}, nil
}

func (s userRPC) UpdatePassword(ctx context.Context, in *user_service.UpdatePasswordRequest) (*user_service.ResponseStatus, error) {
	_, err := s.userUseCase.UpdatePassword(ctx, &entity.UpdatePassword{
		UserID:      in.UserId,
		NewPassword: in.NewPassword,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return &user_service.ResponseStatus{Status: false}, err
	}

	return &user_service.ResponseStatus{Status: true}, nil
}
