package repository

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
)

type Users interface {
	Create(ctx context.Context, kyc *entity.User) (*entity.User, error)
	Update(ctx context.Context, kyc *entity.User) (*entity.User, error)
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64) (*entity.Users, error)
	UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error)
	UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error)
	UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error)
}
