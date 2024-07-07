package usecase

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository"
	"time"
)

type Resume interface {
	CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	DeleteResume(ctx context.Context, resumeID string) error
	DeleteUserResumes(ctx context.Context, userID string) error
	GetResumeByID(ctx context.Context, resumeID string) (*entity.Resume, error)
	GetUserResume(ctx context.Context, userID string, page, limit uint64) (*entity.ListResume, error)
	ListResume(ctx context.Context, request *entity.ListRequest) (*entity.ListResume, error)
}

type resumeService struct {
	BaseUseCase
	repo       repository.Resumes
	ctxTimeout time.Duration
}

func NewResumeService(ctxTimeout time.Duration, repo repository.Resumes) Resume {
	return resumeService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u resumeService) CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.CreateResume(ctx, resume)
}

func (u resumeService) UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.UpdateResume(ctx, resume)
}

func (u resumeService) DeleteResume(ctx context.Context, resumeID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteResume(ctx, resumeID)
}

func (u resumeService) DeleteUserResumes(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.DeleteUserResumes(ctx, userID)
}

func (u resumeService) GetResumeByID(ctx context.Context, resumeID string) (*entity.Resume, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetResumeByID(ctx, resumeID)
}

func (u resumeService) GetUserResume(ctx context.Context, userID string, limit, offset uint64) (*entity.ListResume, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.GetUserResume(ctx, userID, limit, offset)
}

func (u resumeService) ListResume(ctx context.Context, request *entity.ListRequest) (*entity.ListResume, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListResume(ctx, request)
}
