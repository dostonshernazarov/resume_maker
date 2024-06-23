package repo

import (
	"context"

	"github.com/dostonshernazarov/resume_maker/internal/entity"
)

type AppVersionRepo interface {
	Get(ctx context.Context) (*entity.AppVersion, error)
	Create(ctx context.Context, m *entity.AppVersion) error
	Update(ctx context.Context, m *entity.AppVersion) error
}
