package usecase

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository"
	"time"
)

type User interface {
	Create(ctx context.Context, article *entity.User) (*entity.User, error)
	Update(ctx context.Context, article *entity.User) (*entity.User, error)
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64) (*entity.Users, error)
	UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error)
	UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error)
	UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error)
}

type userService struct {
	BaseUseCase
	repo       repository.Users
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.Users) User {
	return userService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, &user.CreatedAt, &user.UpdatedAt)

	return u.repo.Create(ctx, user)
}

func (u userService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, nil, &user.UpdatedAt)

	return u.repo.Update(ctx, user)
}

func (u userService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Delete(ctx, guid)
}

func (u userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Get(ctx, params)
}

func (u userService) List(ctx context.Context, limit, offset uint64) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.List(ctx, limit, offset)
}

func (u userService) UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UniqueEmail(ctx, request)
}

func (u userService) UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UpdateRefresh(ctx, request)
}

func (u userService) UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UpdatePassword(ctx, request)
}
