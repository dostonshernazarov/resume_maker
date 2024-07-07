package repository

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
)

type Resumes interface {
	CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	DeleteResume(ctx context.Context, resumeID string) error
	DeleteUserResumes(ctx context.Context, userID string) error
	GetResumeByID(ctx context.Context, resumeID string) (*entity.Resume, error)
	GetUserResume(ctx context.Context, userID string, limit, offset uint64) (*entity.ListResume, error)
	ListResume(ctx context.Context, request *entity.ListRequest) (*entity.ListResume, error)
}
